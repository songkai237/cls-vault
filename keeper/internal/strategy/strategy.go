package strategy

import (
	"errors"
	"math/big"

	"github.com/clsVault/keeper/internal/math"
)

// QuoteFunc returns (amountOut, sqrtPriceX96After, error) for an exact-input swap.
type QuoteFunc func(amountIn *big.Int, zeroForOne bool) (*big.Int, *big.Int, error)

type SwapSuggestion struct {
	ZeroForOne    bool
	AmountIn      *big.Int
	EstimatedOut  *big.Int
	PostAmount0   *big.Int
	PostAmount1   *big.Int
	PostLiquidity *big.Int
}

// FindBestSwap uses binary search to find the swap amount that maximises
// deployable liquidity within the tick range [sqrtLower, sqrtUpper]
func FindBestSwap(
	amount0, amount1 *big.Int,
	sqrtPriceX96, sqrtLower, sqrtUpper *big.Int,
	quoteFunc QuoteFunc,
	minSwap0, minSwap1 *big.Int,
) (*SwapSuggestion, error) {
	if amount0 == nil || amount1 == nil {
		return nil, errors.New("nil amounts")
	}
	if sqrtPriceX96 == nil || sqrtLower == nil || sqrtUpper == nil {
		return nil, errors.New("nil sqrtPrice or range")
	}
	if quoteFunc == nil {
		return nil, errors.New("nil quoteFunc")
	}

	// ── Step 1: determine which token is in excess ───────────────────────────
	l0 := math.GetLiquidityForAmount0(sqrtPriceX96, sqrtUpper, amount0)
	l1 := math.GetLiquidityForAmount1(sqrtLower, sqrtPriceX96, amount1)

	if l0.Sign() == 0 && l1.Sign() == 0 {
		return nil, nil
	}

	var zeroForOne bool
	var fullAmount, minSwap *big.Int

	switch {
	case l0.Cmp(l1) > 0:
		zeroForOne = true
		fullAmount = new(big.Int).Set(amount0)
		minSwap = minSwap0
	case l1.Cmp(l0) > 0:
		zeroForOne = false
		fullAmount = new(big.Int).Set(amount1)
		minSwap = minSwap1
	default:
		return nil, nil // already balanced
	}

	// ── Step 2: pre-check full swap against price bounds
	hi, err := clampAmountByPrice(fullAmount, zeroForOne, sqrtLower, sqrtUpper, quoteFunc)
	if err != nil {
		return nil, err
	}

	// ── Step 3: binary search on [0, hi]
	lo := big.NewInt(0)
	var best *big.Int

	const maxIter = 64
	for i := 0; i < maxIter; i++ {
		diff := new(big.Int).Sub(hi, lo)
		if diff.Sign() <= 0 || diff.BitLen() <= 1 {
			break
		}

		mid := new(big.Int).Rsh(new(big.Int).Add(lo, hi), 1)
		out, sqrtAfter, err := quoteFunc(mid, zeroForOne)
		if err != nil {
			return nil, err
		}

		a0, a1 := postSwapAmounts(amount0, amount1, mid, out, zeroForOne)
		newL0 := math.GetLiquidityForAmount0(sqrtAfter, sqrtUpper, a0)
		newL1 := math.GetLiquidityForAmount1(sqrtLower, sqrtAfter, a1)

		var stillExcess bool
		if zeroForOne {
			stillExcess = newL0.Cmp(newL1) > 0
		} else {
			stillExcess = newL1.Cmp(newL0) > 0
		}

		if stillExcess {
			lo = mid
			best = new(big.Int).Set(mid)
		} else {
			hi = mid
		}
	}

	if best == nil {
		best = lo
	}

	// ── Step 4: apply minimum-swap threshold
	if best.Sign() <= 0 {
		return nil, nil
	}
	if minSwap != nil && best.Cmp(minSwap) < 0 {
		return nil, nil
	}

	// ── Step 5: build result using post-swap price
	out, sqrtAfter, err := quoteFunc(best, zeroForOne)
	if err != nil {
		return nil, err
	}
	postA0, postA1 := postSwapAmounts(amount0, amount1, best, out, zeroForOne)
	postLiq := math.GetLiquidityForAmounts(sqrtAfter, sqrtLower, sqrtUpper, postA0, postA1)

	return &SwapSuggestion{
		ZeroForOne:    zeroForOne,
		AmountIn:      best,
		EstimatedOut:  out,
		PostAmount0:   postA0,
		PostAmount1:   postA1,
		PostLiquidity: postLiq,
	}, nil
}

// clampAmountByPrice quotes the full excess amount
func clampAmountByPrice(
	fullAmount *big.Int,
	zeroForOne bool,
	sqrtLower, sqrtUpper *big.Int,
	quoteFunc QuoteFunc,
) (*big.Int, error) {
	_, sqrtAfter, err := quoteFunc(fullAmount, zeroForOne)
	if err != nil {
		return nil, err
	}

	if priceInBounds(zeroForOne, sqrtAfter, sqrtLower, sqrtUpper) {
		return new(big.Int).Set(fullAmount), nil
	}

	// Binary search for the maximum amountIn
	lo := big.NewInt(0)
	hi := new(big.Int).Set(fullAmount)

	for i := 0; i < 64; i++ {
		diff := new(big.Int).Sub(hi, lo)
		if diff.Sign() <= 0 || diff.BitLen() <= 1 {
			break
		}

		mid := new(big.Int).Rsh(new(big.Int).Add(lo, hi), 1)
		_, sqrtMid, err := quoteFunc(mid, zeroForOne)
		if err != nil {
			return nil, err
		}

		if priceInBounds(zeroForOne, sqrtMid, sqrtLower, sqrtUpper) {
			lo = mid
		} else {
			hi = mid
		}
	}
	return lo, nil
}

func priceInBounds(zeroForOne bool, sqrtPrice, sqrtLower, sqrtUpper *big.Int) bool {
	if zeroForOne {
		return sqrtPrice.Cmp(sqrtLower) >= 0
	}
	return sqrtPrice.Cmp(sqrtUpper) <= 0
}

func postSwapAmounts(amount0, amount1, amountIn, out *big.Int, zeroForOne bool) (*big.Int, *big.Int) {
	if zeroForOne {
		return new(big.Int).Sub(amount0, amountIn),
			new(big.Int).Add(amount1, out)
	}
	return new(big.Int).Add(amount0, out),
		new(big.Int).Sub(amount1, amountIn)
}
