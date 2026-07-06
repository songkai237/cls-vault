package strategy_test

import (
	"math/big"
	"testing"

	"github.com/clsVault/keeper/internal/math"
	"github.com/clsVault/keeper/internal/strategy"
)

var q96big = new(big.Int).Lsh(big.NewInt(1), 96)

// spotQuote returns a QuoteFunc that uses the spot-price formula.
// Since there is no price impact, sqrtPriceX96After equals the input price.
func spotQuote(sqrtPriceX96 *big.Int, feePPM uint32) strategy.QuoteFunc {
	feeNumer := new(big.Int).SetUint64(uint64(1_000_000 - feePPM))
	feeDenom := big.NewInt(1_000_000)

	return func(amountIn *big.Int, zeroForOne bool) (*big.Int, *big.Int, error) {
		if zeroForOne {
			t := new(big.Int).Mul(amountIn, sqrtPriceX96)
			t.Div(t, q96big)
			t.Mul(t, sqrtPriceX96)
			t.Div(t, q96big)
			t.Mul(t, feeNumer)
			t.Div(t, feeDenom)
			return t, sqrtPriceX96, nil
		}
		t := new(big.Int).Mul(amountIn, q96big)
		t.Div(t, sqrtPriceX96)
		t.Mul(t, q96big)
		t.Div(t, sqrtPriceX96)
		t.Mul(t, feeNumer)
		t.Div(t, feeDenom)
		return t, sqrtPriceX96, nil
	}
}

// rangeAt builds sqrtLower/sqrtUpper from a centre tick and half-range,
// aligned to tickSpacing.
func rangeAt(centerTick, halfRange, tickSpacing int32) (sqrtLower, sqrtUpper *big.Int) {
	lo, hi := math.ComputeTickRange(centerTick, halfRange, tickSpacing)
	return math.GetSqrtRatioAtTick(lo), math.GetSqrtRatioAtTick(hi)
}

func logInputs(t *testing.T, amount0, amount1 *big.Int, tick int32, halfRange, tickSpacing int32) {
	lo, hi := math.ComputeTickRange(tick, halfRange, tickSpacing)
	t.Logf("input: amount0=%s amount1=%s tick=%d range=[%d,%d]", amount0, amount1, tick, lo, hi)
}

func logSwapResult(t *testing.T, sug *strategy.SwapSuggestion, preLiq *big.Int) {
	if sug == nil {
		t.Log("result: no swap suggestion")
		return
	}
	dir := "token0→token1"
	if !sug.ZeroForOne {
		dir = "token1→token0"
	}
	t.Logf("result: dir=%s amountIn=%s estimatedOut=%s", dir, sug.AmountIn, sug.EstimatedOut)
	t.Logf("        postAmount0=%s postAmount1=%s", sug.PostAmount0, sug.PostAmount1)
	if preLiq != nil {
		t.Logf("        preLiq=%s postLiq=%s (delta=%s)", preLiq, sug.PostLiquidity,
			new(big.Int).Sub(sug.PostLiquidity, preLiq))
	} else {
		t.Logf("        postLiq=%s", sug.PostLiquidity)
	}
}

func TestFindBestSwap_ExcessToken0(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	amount0 := new(big.Int).SetUint64(1_000_000)
	amount1 := new(big.Int).SetUint64(100)

	minSwap0 := new(big.Int).SetUint64(10)
	minSwap1 := new(big.Int).SetUint64(10)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, amount1, sqrtPrice, sqrtLo, sqrtHi, qf, minSwap0, minSwap1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, amount0, amount1)
	logInputs(t, amount0, amount1, 0, 1000, 10)
	logSwapResult(t, sug, preLiq)

	if sug == nil {
		t.Fatal("expected a swap suggestion, got nil")
	}
	if !sug.ZeroForOne {
		t.Errorf("expected zeroForOne=true (excess token0), got false")
	}
	if sug.AmountIn.Sign() <= 0 {
		t.Errorf("expected positive AmountIn, got %s", sug.AmountIn)
	}
	if sug.EstimatedOut.Sign() <= 0 {
		t.Errorf("expected positive EstimatedOut, got %s", sug.EstimatedOut)
	}
	if sug.PostLiquidity.Cmp(preLiq) <= 0 {
		t.Errorf("post-swap liq (%s) should exceed pre-swap liq (%s)", sug.PostLiquidity, preLiq)
	}
	if sug.PostAmount0.Sign() < 0 {
		t.Errorf("PostAmount0 negative: %s", sug.PostAmount0)
	}
	if sug.PostAmount1.Sign() < 0 {
		t.Errorf("PostAmount1 negative: %s", sug.PostAmount1)
	}
}

func TestFindBestSwap_ExcessToken1(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	amount0 := new(big.Int).SetUint64(100)
	amount1 := new(big.Int).SetUint64(1_000_000)

	minSwap0 := new(big.Int).SetUint64(10)
	minSwap1 := new(big.Int).SetUint64(10)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, amount1, sqrtPrice, sqrtLo, sqrtHi, qf, minSwap0, minSwap1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, amount0, amount1)
	logInputs(t, amount0, amount1, 0, 1000, 10)
	logSwapResult(t, sug, preLiq)

	if sug == nil {
		t.Fatal("expected a swap suggestion, got nil")
	}
	if sug.ZeroForOne {
		t.Errorf("expected zeroForOne=false (excess token1), got true")
	}
	if sug.PostLiquidity.Cmp(preLiq) <= 0 {
		t.Errorf("post-swap liq (%s) should exceed pre-swap liq (%s)", sug.PostLiquidity, preLiq)
	}
}

func TestFindBestSwap_Balanced_NoSwap(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	refLiq := new(big.Int).SetUint64(1_000_000_000)
	n0, n1 := math.GetAmountsForLiquidity(sqrtPrice, sqrtLo, sqrtHi, refLiq)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(n0, n1, sqrtPrice, sqrtLo, sqrtHi, qf,
		big.NewInt(0), big.NewInt(0))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	logInputs(t, n0, n1, 0, 1000, 10)
	t.Logf("        balanced refLiq=%s", refLiq)
	logSwapResult(t, sug, nil)

	if sug != nil {
		t.Errorf("expected nil suggestion for balanced amounts, got: zeroForOne=%v amountIn=%s",
			sug.ZeroForOne, sug.AmountIn)
	}
}

func TestFindBestSwap_ExcessBelowMinSwap(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	refLiq := new(big.Int).SetUint64(1_000_000_000)
	n0, n1 := math.GetAmountsForLiquidity(sqrtPrice, sqrtLo, sqrtHi, refLiq)
	amount0 := new(big.Int).Add(n0, big.NewInt(1))

	minSwap0 := new(big.Int).SetUint64(1000)
	minSwap1 := new(big.Int).SetUint64(1000)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, n1, sqrtPrice, sqrtLo, sqrtHi, qf, minSwap0, minSwap1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	excess := new(big.Int).Sub(amount0, n0)
	logInputs(t, amount0, n1, 0, 1000, 10)
	t.Logf("        excess0=%s minSwap0=%s", excess, minSwap0)
	logSwapResult(t, sug, nil)

	if sug != nil {
		t.Errorf("expected nil for excess below minSwap, got amountIn=%s", sug.AmountIn)
	}
}

func TestFindBestSwap_AmountIn_GEQ_MinSwap(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	amount0 := new(big.Int).SetUint64(500_000)
	amount1 := new(big.Int).SetUint64(100)
	minSwap0 := new(big.Int).SetUint64(50)
	minSwap1 := new(big.Int).SetUint64(50)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, amount1, sqrtPrice, sqrtLo, sqrtHi, qf, minSwap0, minSwap1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, amount0, amount1)
	logInputs(t, amount0, amount1, 0, 1000, 10)
	t.Logf("        minSwap0=%s", minSwap0)
	logSwapResult(t, sug, preLiq)

	if sug == nil {
		t.Fatal("expected a suggestion")
	}
	if sug.AmountIn.Cmp(minSwap0) < 0 {
		t.Errorf("AmountIn (%s) < minSwap0 (%s)", sug.AmountIn, minSwap0)
	}
}

func TestFindBestSwap_PriceAboveRange(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(200)
	sqrtLo := math.GetSqrtRatioAtTick(0)
	sqrtHi := math.GetSqrtRatioAtTick(100)

	amount0 := new(big.Int).SetUint64(500_000)
	amount1 := new(big.Int).SetUint64(0)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, amount1, sqrtPrice, sqrtLo, sqrtHi, qf,
		big.NewInt(100), big.NewInt(100))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, amount0, amount1)
	t.Logf("input: amount0=%s amount1=%s tick=200 range=[0,100] (price above range)", amount0, amount1)
	t.Logf("        preLiq=%s", preLiq)
	logSwapResult(t, sug, preLiq)

	_ = sug
}

func TestFindBestSwap_NilAmounts(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	qf := spotQuote(sqrtPrice, 3000)
	_, err := strategy.FindBestSwap(nil, big.NewInt(100), sqrtPrice, sqrtLo, sqrtHi, qf,
		big.NewInt(0), big.NewInt(0))
	if err == nil {
		t.Fatal("expected error for nil amount0")
	}
}

func TestFindBestSwap_NilQuoteFunc(t *testing.T) {
	sqrtPrice := math.GetSqrtRatioAtTick(0)
	sqrtLo, sqrtHi := rangeAt(0, 1000, 10)

	_, err := strategy.FindBestSwap(big.NewInt(100), big.NewInt(100), sqrtPrice, sqrtLo, sqrtHi, nil,
		big.NewInt(0), big.NewInt(0))
	if err == nil {
		t.Fatal("expected error for nil quoteFunc")
	}
}

func TestFindBestSwap_SepoliaLike(t *testing.T) {
	const currentTick = int32(177249)
	const halfRange = int32(600)
	const tickSpacing = int32(60)

	sqrtPrice := math.GetSqrtRatioAtTick(currentTick)
	sqrtLo, sqrtHi := rangeAt(currentTick, halfRange, tickSpacing)

	amount0 := new(big.Int).SetUint64(10_000_000)
	amount1, _ := new(big.Int).SetString("1000000000000000", 10)

	minSwap0 := new(big.Int).SetUint64(100_000)
	minSwap1, _ := new(big.Int).SetString("10000000000000", 10)

	qf := spotQuote(sqrtPrice, 3000)
	sug, err := strategy.FindBestSwap(amount0, amount1, sqrtPrice, sqrtLo, sqrtHi, qf, minSwap0, minSwap1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	preLiq := math.GetLiquidityForAmounts(sqrtPrice, sqrtLo, sqrtHi, amount0, amount1)
	logInputs(t, amount0, amount1, currentTick, halfRange, tickSpacing)
	t.Logf("        minSwap0=%s minSwap1=%s fee=3000", minSwap0, minSwap1)
	logSwapResult(t, sug, preLiq)

	if sug == nil {
		t.Fatal("expected a swap suggestion for unbalanced USDC/WETH position")
	}

	if sug.PostLiquidity.Cmp(preLiq) <= 0 {
		t.Errorf("post-swap liq (%s) should be > pre-swap liq (%s)", sug.PostLiquidity, preLiq)
	}
}
