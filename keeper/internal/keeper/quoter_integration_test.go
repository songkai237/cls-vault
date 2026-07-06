package keeper_test

import (
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/clsVault/keeper/internal/contracts"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func TestQuoter_QuoteExactInputSingle(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	rpcURL := os.Getenv("RPC_URL")
	quoterAddr := os.Getenv("QUOTER_ADDRESS")
	token0Str := os.Getenv("TOKEN0")
	token1Str := os.Getenv("TOKEN1")
	feeStr := "3000"
	fmt.Println(rpcURL, quoterAddr, token0Str, token1Str)
	if rpcURL == "" || quoterAddr == "" || token0Str == "" || token1Str == "" {

		t.Skip("set RPC_URL, QUOTER_ADDRESS, TOKEN0_ADDRESS, TOKEN1_ADDRESS, FEE to run this test")
	}

	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		t.Fatalf("dial rpc: %v", err)
	}
	defer client.Close()

	token0 := common.HexToAddress(token0Str)
	token1 := common.HexToAddress(token1Str)
	feeBI, ok := new(big.Int).SetString(feeStr, 10)
	if !ok {
		t.Fatalf("invalid FEE: %s", feeStr)
	}

	quoter, err := contracts.NewIQuoterV2(common.HexToAddress(quoterAddr), client)
	if err != nil {
		t.Fatalf("bind quoter: %v", err)
	}

	raw := &contracts.IQuoterV2CallerRaw{Contract: &quoter.IQuoterV2Caller}
	callOpts := &bind.CallOpts{}

	// Default: 1 WETH-equivalent for token0, 1 USDC-equivalent for token1.
	amount0, _ := new(big.Int).SetString("1000000000000000000", 10) // 1e18
	amount1, _ := new(big.Int).SetString("1000000", 10)             // 1e6
	if v := os.Getenv("AMOUNT_IN"); v != "" {
		amount0, ok = new(big.Int).SetString(v, 10)
		if !ok {
			t.Fatalf("invalid AMOUNT_IN: %s", v)
		}
		amount1 = new(big.Int).Set(amount0)
	}

	t.Logf("quoter=%s token0=%s token1=%s fee=%s", quoterAddr, token0Str, token1Str, feeStr)
	t.Logf("rpc=%s", rpcURL)
	fmt.Println()

	// ── token0 → token1 ─────────────────────────────────────────────────────
	params0to1 := contracts.IQuoterV2QuoteExactInputSingleParams{
		TokenIn:           token0,
		TokenOut:          token1,
		AmountIn:          amount0,
		Fee:               feeBI,
		SqrtPriceLimitX96: minSqrtRatioPlusOne(),
	}

	var out0to1 []interface{}
	if err := raw.Call(callOpts, &out0to1, "quoteExactInputSingle", params0to1); err != nil {
		t.Fatalf("quoteExactInputSingle(token0→token1): %v", err)
	}

	fmt.Printf("═══ token0 → token1 ═══\n")
	fmt.Printf("  amountIn:           %s\n", amount0.String())
	fmt.Printf("  amountOut:          %s\n", out0to1[0].(*big.Int).String())
	fmt.Printf("  sqrtPriceX96After:  %s\n", out0to1[1].(*big.Int).String())
	fmt.Printf("  ticksCrossed:       %d\n", out0to1[2].(uint32))
	fmt.Printf("  gasEstimate:        %s\n", out0to1[3].(*big.Int).String())
	fmt.Println()

	// ── token1 → token0 ─────────────────────────────────────────────────────
	params1to0 := contracts.IQuoterV2QuoteExactInputSingleParams{
		TokenIn:           token1,
		TokenOut:          token0,
		AmountIn:          amount1,
		Fee:               feeBI,
		SqrtPriceLimitX96: maxSqrtRatioMinusOne(),
	}

	var out1to0 []interface{}
	if err := raw.Call(callOpts, &out1to0, "quoteExactInputSingle", params1to0); err != nil {
		t.Fatalf("quoteExactInputSingle(token1→token0): %v", err)
	}

	fmt.Printf("═══ token1 → token0 ═══\n")
	fmt.Printf("  amountIn:           %s\n", amount1.String())
	fmt.Printf("  amountOut:          %s\n", out1to0[0].(*big.Int).String())
	fmt.Printf("  sqrtPriceX96After:  %s\n", out1to0[1].(*big.Int).String())
	fmt.Printf("  ticksCrossed:       %d\n", out1to0[2].(uint32))
	fmt.Printf("  gasEstimate:        %s\n", out1to0[3].(*big.Int).String())
	fmt.Println()

	t.Log("quoter integration test passed")
}

func minSqrtRatioPlusOne() *big.Int {
	return big.NewInt(4295128740)
}

func maxSqrtRatioMinusOne() *big.Int {
	r, _ := new(big.Int).SetString("1461446703485210103287273052203988822378723970341", 10)
	return r
}
