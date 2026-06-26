package keeper

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"

	"github.com/clsVault/keeper/internal/config"
	"github.com/clsVault/keeper/internal/contracts"
	"github.com/clsVault/keeper/internal/math"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type Action string

const (
	ActionNone     Action = "none"
	ActionMaintain Action = "maintain"
	ActionCollect  Action = "collect"
)

type Decision struct {
	Action  Action
	Reason  string
	Summary string
}

type Keeper struct {
	cfg      config.Config
	client   *ethclient.Client
	strategy *contracts.UniswapV3Strategy
	pool     *contracts.UniswapV3Pool
	token0   *contracts.ERC20
	token1   *contracts.ERC20
	chainID  *big.Int

	strategyAddr common.Address
	token0Addr   common.Address
	token1Addr   common.Address

	halfRangeTicks int32
	minSwap0       *big.Int
	minSwap1       *big.Int
	minIdle0       *big.Int
	minIdle1       *big.Int

	transactOpts *bind.TransactOpts
}

func New(cfg config.Config) (*Keeper, error) {
	client, err := ethclient.Dial(cfg.RPCURL)
	if err != nil {
		return nil, fmt.Errorf("dial rpc: %w", err)
	}

	strategyAddr := common.HexToAddress(cfg.StrategyAddress)
	strategy, err := contracts.NewUniswapV3Strategy(strategyAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind strategy: %w", err)
	}

	poolAddr, err := strategy.GetPool(nil)
	if err != nil {
		return nil, fmt.Errorf("get pool: %w", err)
	}
	pool, err := contracts.NewUniswapV3Pool(poolAddr, client)
	if err != nil {
		return nil, fmt.Errorf("bind pool: %w", err)
	}

	token0Addr, err := strategy.GetToken0(nil)
	if err != nil {
		return nil, fmt.Errorf("get token0: %w", err)
	}
	token1Addr, err := strategy.GetToken1(nil)
	if err != nil {
		return nil, fmt.Errorf("get token1: %w", err)
	}

	token0, err := contracts.NewERC20(token0Addr, client)
	if err != nil {
		return nil, fmt.Errorf("bind token0: %w", err)
	}
	token1, err := contracts.NewERC20(token1Addr, client)
	if err != nil {
		return nil, fmt.Errorf("bind token1: %w", err)
	}

	halfRangeBI, err := strategy.GetHalfRangeTicks(nil)
	if err != nil {
		return nil, fmt.Errorf("get half range: %w", err)
	}
	minSwap0, err := strategy.GetMinSwapAmount0(nil)
	if err != nil {
		return nil, fmt.Errorf("get min swap0: %w", err)
	}
	minSwap1, err := strategy.GetMinSwapAmount1(nil)
	if err != nil {
		return nil, fmt.Errorf("get min swap1: %w", err)
	}

	minIdle0 := new(big.Int).Set(minSwap0)
	minIdle1 := new(big.Int).Set(minSwap1)
	if cfg.MinIdleAmount0 != "" {
		v, ok := new(big.Int).SetString(cfg.MinIdleAmount0, 10)
		if !ok {
			return nil, fmt.Errorf("invalid MIN_IDLE_AMOUNT0")
		}
		minIdle0 = v
	}
	if cfg.MinIdleAmount1 != "" {
		v, ok := new(big.Int).SetString(cfg.MinIdleAmount1, 10)
		if !ok {
			return nil, fmt.Errorf("invalid MIN_IDLE_AMOUNT1")
		}
		minIdle1 = v
	}

	var chainID *big.Int
	if cfg.ChainID != 0 {
		chainID = big.NewInt(cfg.ChainID)
	} else {
		chainID, err = client.ChainID(context.Background())
		if err != nil {
			return nil, fmt.Errorf("chain id: %w", err)
		}
	}

	privateKey, err := crypto.HexToECDSA(strings.TrimPrefix(cfg.PrivateKey, "0x"))
	if err != nil {
		return nil, fmt.Errorf("private key: %w", err)
	}
	transactOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return nil, fmt.Errorf("transactor: %w", err)
	}

	k := &Keeper{
		cfg:            cfg,
		client:         client,
		strategy:       strategy,
		pool:           pool,
		token0:         token0,
		token1:         token1,
		chainID:        chainID,
		strategyAddr:   strategyAddr,
		token0Addr:     token0Addr,
		token1Addr:     token1Addr,
		halfRangeTicks: math.Int24FromBigInt(halfRangeBI),
		minSwap0:       minSwap0,
		minSwap1:       minSwap1,
		minIdle0:       minIdle0,
		minIdle1:       minIdle1,
		transactOpts:   transactOpts,
	}

	log.Printf("keeper ready strategy=%s pool=%s token0=%s token1=%s halfRange=%d",
		strategyAddr.Hex(), poolAddr.Hex(), token0Addr.Hex(), token1Addr.Hex(), k.halfRangeTicks)

	return k, nil
}

func (k *Keeper) Run(ctx context.Context) error {
	ticker := time.NewTicker(k.cfg.PollInterval)
	defer ticker.Stop()

	for {
		if err := k.tick(ctx); err != nil {
			log.Printf("keeper tick error: %v", err)
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}

func (k *Keeper) tick(ctx context.Context) error {
	decision, state, err := k.Evaluate(ctx)
	if err != nil {
		return err
	}

	log.Printf("monitor tick=%d pos=[%d,%d] liq=%s idle0=%s idle1=%s | %s",
		state.CurrentTick, state.TickLower, state.TickUpper,
		state.Liquidity.String(), state.Idle0.String(), state.Idle1.String(),
		state.Summary)

	if decision.Action == ActionNone {
		log.Printf("no action: %s", decision.Reason)
		return nil
	}

	log.Printf("action=%s reason=%s", decision.Action, decision.Reason)
	if k.cfg.DryRun {
		log.Printf("dry-run enabled, skip tx")
		return nil
	}

	switch decision.Action {
	case ActionMaintain:
		return k.sendMaintain(ctx)
	case ActionCollect:
		return k.sendCollect(ctx)
	}
	return nil
}

type State struct {
	CurrentTick int32
	TickLower   int32
	TickUpper   int32
	Liquidity   *big.Int
	Idle0       *big.Int
	Idle1       *big.Int
	TokensOwed0 *big.Int
	TokensOwed1 *big.Int
	Summary     string
}

func (k *Keeper) Evaluate(ctx context.Context) (Decision, State, error) {
	_ = ctx
	callOpts := &bind.CallOpts{Context: ctx}

	slot0, err := k.pool.Slot0(callOpts)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("slot0: %w", err)
	}
	currentTick := math.Int24FromBigInt(slot0.Tick)

	pos, err := k.strategy.GetPosition(callOpts)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("getPosition: %w", err)
	}

	tickLower := math.Int24FromBigInt(pos.TickLower)
	tickUpper := math.Int24FromBigInt(pos.TickUpper)
	hasPosition := pos.Liquidity.Sign() > 0

	idle0, err := k.token0.BalanceOf(callOpts, k.strategyAddr)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("idle0: %w", err)
	}
	idle1, err := k.token1.BalanceOf(callOpts, k.strategyAddr)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("idle1: %w", err)
	}

	spacingBI, err := k.pool.TickSpacing(callOpts)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("tickSpacing: %w", err)
	}
	tickSpacing := math.Int24FromBigInt(spacingBI)

	needRebalanceOnChain, err := k.strategy.NeedRebalance(callOpts)
	if err != nil {
		return Decision{}, State{}, fmt.Errorf("needRebalance: %w", err)
	}
	needRebalanceLocal := math.NeedRebalance(currentTick, tickLower, tickUpper, hasPosition)

	expectedLower, expectedUpper := math.ComputeTickRange(currentTick, k.halfRangeTicks, tickSpacing)

	state := State{
		CurrentTick: currentTick,
		TickLower:   tickLower,
		TickUpper:   tickUpper,
		Liquidity:   pos.Liquidity,
		Idle0:       idle0,
		Idle1:       idle1,
		TokensOwed0: pos.TokensOwed0,
		TokensOwed1: pos.TokensOwed1,
		Summary: fmt.Sprintf("needRebalance(chain=%v local=%v) expectedRange=[%d,%d]",
			needRebalanceOnChain, needRebalanceLocal, expectedLower, expectedUpper),
	}

	if needRebalanceOnChain {
		return Decision{
			Action: ActionMaintain,
			Reason: "price out of position range or no position",
		}, state, nil
	}

	sqrtPrice := slot0.SqrtPriceX96
	var sqrtLower, sqrtUpper *big.Int
	if hasPosition {
		sqrtLower = math.GetSqrtRatioAtTick(tickLower)
		sqrtUpper = math.GetSqrtRatioAtTick(tickUpper)
	} else {
		sqrtLower = math.GetSqrtRatioAtTick(expectedLower)
		sqrtUpper = math.GetSqrtRatioAtTick(expectedUpper)
	}

	if math.ShouldSwap(idle0, idle1, sqrtPrice, sqrtLower, sqrtUpper, k.minSwap0, k.minSwap1) {
		return Decision{
			Action: ActionMaintain,
			Reason: "idle token excess above min swap threshold",
		}, state, nil
	}

	hasIdle := idle0.Cmp(k.minIdle0) >= 0 || idle1.Cmp(k.minIdle1) >= 0
	if hasIdle && hasPosition {
		return Decision{
			Action: ActionMaintain,
			Reason: "idle tokens to mint into existing position",
		}, state, nil
	}

	if !hasPosition && (idle0.Sign() > 0 || idle1.Sign() > 0) {
		return Decision{
			Action: ActionMaintain,
			Reason: "idle tokens without open position",
		}, state, nil
	}

	if pos.TokensOwed0.Sign() > 0 || pos.TokensOwed1.Sign() > 0 {
		return Decision{
			Action: ActionCollect,
			Reason: "uncollected fees on position",
		}, state, nil
	}

	return Decision{Action: ActionNone, Reason: "position healthy"}, state, nil
}

func (k *Keeper) sendMaintain(ctx context.Context) error {
	opts := *k.transactOpts
	opts.Context = ctx
	tx, err := k.strategy.Maintain(&opts)
	if err != nil {
		return fmt.Errorf("maintain tx: %w", err)
	}
	return k.waitTx(ctx, tx, "maintain")
}

func (k *Keeper) sendCollect(ctx context.Context) error {
	opts := *k.transactOpts
	opts.Context = ctx
	tx, err := k.strategy.Collect(&opts)
	if err != nil {
		return fmt.Errorf("collect tx: %w", err)
	}
	return k.waitTx(ctx, tx, "collect")
}

func (k *Keeper) waitTx(ctx context.Context, tx *types.Transaction, label string) error {
	log.Printf("sent %s tx=%s", label, tx.Hash().Hex())
	receipt, err := bind.WaitMined(ctx, k.client, tx)
	if err != nil {
		return fmt.Errorf("wait %s: %w", label, err)
	}
	if receipt.Status != types.ReceiptStatusSuccessful {
		return fmt.Errorf("%s reverted gasUsed=%d", label, receipt.GasUsed)
	}
	log.Printf("%s confirmed block=%d gasUsed=%d", label, receipt.BlockNumber.Uint64(), receipt.GasUsed)
	return nil
}
