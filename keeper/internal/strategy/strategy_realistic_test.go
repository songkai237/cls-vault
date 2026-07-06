package strategy_test

import (
	"math/big"
	"testing"

	"github.com/clsVault/keeper/internal/math"
	"github.com/clsVault/keeper/internal/strategy"
)

// Realistic pool parameters aligned with script/Deploy.s.sol and fork tests.

const (
	mainnetWETHUSDCTick        = int32(-195800)
	mainnetWETHUSDCTickSpacing  = int32(10)
	mainnetWETHUSDCFee          = uint32(500)

	sepoliaUSDCWETHTick        = int32(177249)
	sepoliaUSDCWETHTickSpacing  = int32(60)
	sepoliaUSDCWETHFee          = uint32(3000)

	defaultHalfRange = int32(600)
)

var (
	mainnetMinSwapWETH = mustBig("10000000000000000")
	mainnetMinSwapUSDC = mustBig("100000000")

	sepoliaMinSwapUSDC = mustBig("100000000")
	sepoliaMinSwapWETH = mustBig("10000000000000000")
)

type realisticCase struct {
	name            string
	pool            string
	tick            int32
	halfRange       int32
	tickSpacing     int32
	fee             uint32
	amount0         *big.Int
	amount1         *big.Int
	minSwap0        *big.Int
	minSwap1        *big.Int
	wantSwap        bool
	wantZeroForOne  *bool
}

func mustBig(s string) *big.Int {
	v, ok := new(big.Int).SetString(s, 10)
	if !ok {
		panic("invalid big.Int: " + s)
	}
	return v
}

func eth(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), mustBig("1000000000000000000"))
}

func usdc(n int64) *big.Int {
	return new(big.Int).Mul(big.NewInt(n), big.NewInt(1_000_000))
}

func runRealisticCase(t *testing.T, c realisticCase) {
	t.Helper()
	t.Logf("=== %s ===", c.name)
	t.Logf("pool: %s", c.pool)

	sqrtPrice := math.GetSqrtRatioAtTick(c.tick)
	sqrtLo, sqrtHi := rangeAt(c.tick, c.halfRange, c.tickSpacing)
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, c.amount0, c.amount1)

	logInputs(t, c.amount0, c.amount1, c.tick, c.halfRange, c.tickSpacing)
	t.Logf("        fee=%d minSwap0=%s minSwap1=%s preLiq=%s", c.fee, c.minSwap0, c.minSwap1, preLiq)

	qf := spotQuote(sqrtPrice, c.fee)
	sug, err := strategy.FindBestSwap(
		c.amount0, c.amount1, sqrtPrice, sqrtLo, sqrtHi,
		qf, c.minSwap0, c.minSwap1,
	)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	logSwapResult(t, sug, preLiq)

	if c.wantSwap && sug == nil {
		t.Fatal("expected swap suggestion, got nil")
	}
	if !c.wantSwap && sug != nil {
		t.Fatalf("expected no swap, got amountIn=%s", sug.AmountIn)
	}
	if sug == nil {
		return
	}

	if sug.AmountIn.Sign() <= 0 {
		t.Errorf("AmountIn must be positive, got %s", sug.AmountIn)
	}
	if sug.EstimatedOut.Sign() <= 0 {
		t.Errorf("EstimatedOut must be positive, got %s", sug.EstimatedOut)
	}
	if sug.PostAmount0.Sign() < 0 || sug.PostAmount1.Sign() < 0 {
		t.Errorf("negative post amounts: amount0=%s amount1=%s", sug.PostAmount0, sug.PostAmount1)
	}
	if sug.PostLiquidity.Cmp(preLiq) <= 0 {
		t.Errorf("postLiq (%s) should exceed preLiq (%s)", sug.PostLiquidity, preLiq)
	}
	if c.wantZeroForOne != nil && sug.ZeroForOne != *c.wantZeroForOne {
		t.Errorf("zeroForOne=%v, want %v", sug.ZeroForOne, *c.wantZeroForOne)
	}
}

func TestFindBestSwap_Realistic_MainnetForkDeposit(t *testing.T) {
	runRealisticCase(t, realisticCase{
		name:        "mainnet fork deposit ratio",
		pool:        "WETH/USDC 0.05% mainnet",
		tick:        mainnetWETHUSDCTick,
		halfRange:   defaultHalfRange,
		tickSpacing: mainnetWETHUSDCTickSpacing,
		fee:         mainnetWETHUSDCFee,
		amount0:     eth(10),
		amount1:     usdc(30_000),
		minSwap0:    mainnetMinSwapWETH,
		minSwap1:    mainnetMinSwapUSDC,
		wantSwap:    true,
	})
}

func TestFindBestSwap_Realistic_MainnetExcessWETH(t *testing.T) {
	want := true
	runRealisticCase(t, realisticCase{
		name:           "mainnet excess WETH",
		pool:           "WETH/USDC 0.05% mainnet",
		tick:           mainnetWETHUSDCTick,
		halfRange:      defaultHalfRange,
		tickSpacing:    mainnetWETHUSDCTickSpacing,
		fee:            mainnetWETHUSDCFee,
		amount0:        eth(15),
		amount1:        usdc(30_000),
		minSwap0:       mainnetMinSwapWETH,
		minSwap1:       mainnetMinSwapUSDC,
		wantSwap:       true,
		wantZeroForOne: &want,
	})
}

func TestFindBestSwap_Realistic_MainnetExcessUSDC(t *testing.T) {
	want := false
	runRealisticCase(t, realisticCase{
		name:           "mainnet excess USDC",
		pool:           "WETH/USDC 0.05% mainnet",
		tick:           mainnetWETHUSDCTick,
		halfRange:      defaultHalfRange,
		tickSpacing:    mainnetWETHUSDCTickSpacing,
		fee:            mainnetWETHUSDCFee,
		amount0:        eth(10),
		amount1:        usdc(50_000),
		minSwap0:       mainnetMinSwapWETH,
		minSwap1:       mainnetMinSwapUSDC,
		wantSwap:       true,
		wantZeroForOne: &want,
	})
}

func TestFindBestSwap_Realistic_MainnetSmallRetail(t *testing.T) {
	halfEth := mustBig("500000000000000000")
	runRealisticCase(t, realisticCase{
		name:        "mainnet 0.5 ETH + 1500 USDC",
		pool:        "WETH/USDC 0.05% mainnet",
		tick:        mainnetWETHUSDCTick,
		halfRange:   defaultHalfRange,
		tickSpacing: mainnetWETHUSDCTickSpacing,
		fee:         mainnetWETHUSDCFee,
		amount0:     halfEth,
		amount1:     usdc(1_500),
		minSwap0:    mainnetMinSwapWETH,
		minSwap1:    mainnetMinSwapUSDC,
		wantSwap:    true,
	})
}

func TestFindBestSwap_Realistic_MainnetNearRangeEdge(t *testing.T) {
	_, hi := math.ComputeTickRange(mainnetWETHUSDCTick, defaultHalfRange, mainnetWETHUSDCTickSpacing)
	tickNearUpper := hi - mainnetWETHUSDCTickSpacing

	want := true
	runRealisticCase(t, realisticCase{
		name:           "mainnet tick near range upper",
		pool:           "WETH/USDC 0.05% mainnet",
		tick:           tickNearUpper,
		halfRange:      defaultHalfRange,
		tickSpacing:    mainnetWETHUSDCTickSpacing,
		fee:            mainnetWETHUSDCFee,
		amount0:        eth(5),
		amount1:        usdc(5_000),
		minSwap0:       mainnetMinSwapWETH,
		minSwap1:       mainnetMinSwapUSDC,
		wantSwap:       true,
		wantZeroForOne: &want,
	})
}

func TestFindBestSwap_Realistic_SepoliaMediumVault(t *testing.T) {
	runRealisticCase(t, realisticCase{
		name:        "sepolia medium vault",
		pool:        "USDC/WETH 0.3% sepolia",
		tick:        sepoliaUSDCWETHTick,
		halfRange:   defaultHalfRange,
		tickSpacing: sepoliaUSDCWETHTickSpacing,
		fee:         sepoliaUSDCWETHFee,
		amount0:     usdc(100),
		amount1:     mustBig("50000000000000000"),
		minSwap0:    sepoliaMinSwapUSDC,
		minSwap1:    sepoliaMinSwapWETH,
		wantSwap:    true,
	})
}

func TestFindBestSwap_Realistic_SepoliaBalancedNoSwap(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(sepoliaUSDCWETHTick)
	sqrtLo, sqrtHi := rangeAt(sepoliaUSDCWETHTick, defaultHalfRange, sepoliaUSDCWETHTickSpacing)

	refLiq := mustBig("5000000000000000000")
	n0, n1 := math.GetAmountsForLiquidity(sqrtPrice, sqrtLo, sqrtHi, refLiq)

	runRealisticCase(t, realisticCase{
		name:        "sepolia balanced LP amounts",
		pool:        "USDC/WETH 0.3% sepolia",
		tick:        sepoliaUSDCWETHTick,
		halfRange:   defaultHalfRange,
		tickSpacing: sepoliaUSDCWETHTickSpacing,
		fee:         sepoliaUSDCWETHFee,
		amount0:     n0,
		amount1:     n1,
		minSwap0:    sepoliaMinSwapUSDC,
		minSwap1:    sepoliaMinSwapWETH,
		wantSwap:    false,
	})
}

func TestFindBestSwap_Realistic_MaintainAfterCollect(t *testing.T) {
	idle0 := eth(2)
	idle1 := usdc(8_000)
	owed0 := mustBig("500000000000000")
	owed1 := usdc(50)

	sim0 := new(big.Int).Add(idle0, owed0)
	sim1 := new(big.Int).Add(idle1, owed1)

	t.Log("=== maintain simulation (idle + owed fees) ===")
	runRealisticCase(t, realisticCase{
		name:        "mainnet maintain idle+fees",
		pool:        "WETH/USDC 0.05% mainnet",
		tick:        mainnetWETHUSDCTick,
		halfRange:   defaultHalfRange,
		tickSpacing: mainnetWETHUSDCTickSpacing,
		fee:         mainnetWETHUSDCFee,
		amount0:     sim0,
		amount1:     sim1,
		minSwap0:    mainnetMinSwapWETH,
		minSwap1:    mainnetMinSwapUSDC,
		wantSwap:    true,
	})
}
