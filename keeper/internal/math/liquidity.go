package math

import (
	"math/big"

	"github.com/holiman/uint256"
)

var q96 = new(uint256.Int).Lsh(uint256.NewInt(1), 96)

// ShouldSwap mirrors on-chain _basicSwap excess + min amount checks.
func ShouldSwap(
	amount0, amount1 *big.Int,
	sqrtPriceX96, sqrtLower, sqrtUpper *big.Int,
	minSwap0, minSwap1 *big.Int,
) bool {
	if (amount0 == nil || amount0.Sign() == 0) && (amount1 == nil || amount1.Sign() == 0) {
		return false
	}

	liquidity := getLiquidityForAmounts(sqrtPriceX96, sqrtLower, sqrtUpper, amount0, amount1)
	if liquidity.Sign() == 0 {
		return false
	}

	need0, need1 := getAmountsForLiquidity(sqrtPriceX96, sqrtLower, sqrtUpper, liquidity)

	if amount0.Cmp(need0) > 0 {
		excess0 := new(big.Int).Sub(amount0, need0)
		return excess0.Cmp(minSwap0) >= 0
	}
	if amount1.Cmp(need1) > 0 {
		excess1 := new(big.Int).Sub(amount1, need1)
		return excess1.Cmp(minSwap1) >= 0
	}
	return false
}

func getLiquidityForAmounts(
	sqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96 *big.Int,
	amount0, amount1 *big.Int,
) *big.Int {
	sqrtA, sqrtB := sortRatios(sqrtRatioAX96, sqrtRatioBX96)
	price, a, b := toU256(sqrtPriceX96), toU256(sqrtA), toU256(sqrtB)

	if price.Cmp(a) <= 0 {
		return getLiquidityForAmount0(a, b, toU256(amount0)).ToBig()
	}
	if price.Cmp(b) < 0 {
		l0 := getLiquidityForAmount0(price, b, toU256(amount0))
		l1 := getLiquidityForAmount1(a, price, toU256(amount1))
		if l0.Cmp(l1) < 0 {
			return l1.ToBig()
		}
		return l0.ToBig()
	}
	return getLiquidityForAmount1(a, b, toU256(amount1)).ToBig()
}

func getAmountsForLiquidity(
	sqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
) (*big.Int, *big.Int) {
	sqrtA, sqrtB := sortRatios(sqrtRatioAX96, sqrtRatioBX96)
	price, a, b := toU256(sqrtPriceX96), toU256(sqrtA), toU256(sqrtB)
	liq := toU256(liquidity)

	var amount0, amount1 uint256.Int
	if price.Cmp(a) <= 0 {
		amount0 = *getAmount0ForLiquidity(a, b, liq)
	} else if price.Cmp(b) < 0 {
		amount0 = *getAmount0ForLiquidity(price, b, liq)
		amount1 = *getAmount1ForLiquidity(a, price, liq)
	} else {
		amount1 = *getAmount1ForLiquidity(a, b, liq)
	}
	return amount0.ToBig(), amount1.ToBig()
}

func getLiquidityForAmount0(sqrtA, sqrtB, amount0 *uint256.Int) *uint256.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	intermediate := mulDiv(sqrtA, sqrtB, q96)
	return mulDiv(amount0, intermediate, new(uint256.Int).Sub(sqrtB, sqrtA))
}

func getLiquidityForAmount1(sqrtA, sqrtB, amount1 *uint256.Int) *uint256.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	return mulDiv(amount1, q96, new(uint256.Int).Sub(sqrtB, sqrtA))
}

func getAmount0ForLiquidity(sqrtA, sqrtB, liquidity *uint256.Int) *uint256.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	numerator := mulDiv(liquidity, new(uint256.Int).Sub(sqrtB, sqrtA), sqrtB)
	return mulDiv(numerator, q96, sqrtA)
}

func getAmount1ForLiquidity(sqrtA, sqrtB, liquidity *uint256.Int) *uint256.Int {
	if sqrtA.Cmp(sqrtB) > 0 {
		sqrtA, sqrtB = sqrtB, sqrtA
	}
	return mulDiv(liquidity, new(uint256.Int).Sub(sqrtB, sqrtA), q96)
}

func mulDiv(a, b, denom *uint256.Int) *uint256.Int {
	product := new(uint256.Int).Mul(a, b)
	return new(uint256.Int).Div(product, denom)
}

func sortRatios(a, b *big.Int) (*big.Int, *big.Int) {
	if a.Cmp(b) > 0 {
		return b, a
	}
	return a, b
}

func toU256(v *big.Int) *uint256.Int {
	if v == nil {
		return uint256.NewInt(0)
	}
	out, overflow := uint256.FromBig(v)
	if overflow {
		panic("uint256 overflow")
	}
	return out
}
