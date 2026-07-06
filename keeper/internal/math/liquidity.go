package math

import "math/big"

var q96 = new(big.Int).Lsh(big.NewInt(1), 96)

func GetLiquidityForAmounts(
	sqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96 *big.Int,
	amount0, amount1 *big.Int,
) *big.Int {
	sqrtA, sqrtB := sortRatios(sqrtRatioAX96, sqrtRatioBX96)
	price := bigOrZero(sqrtPriceX96)

	if price.Cmp(sqrtA) <= 0 {
		return getLiquidityForAmount0(sqrtA, sqrtB, bigOrZero(amount0))
	}
	if price.Cmp(sqrtB) < 0 {
		l0 := getLiquidityForAmount0(price, sqrtB, bigOrZero(amount0))
		l1 := getLiquidityForAmount1(sqrtA, price, bigOrZero(amount1))
		if l0.Cmp(l1) < 0 {
			return l0
		}
		return l1
	}
	return getLiquidityForAmount1(sqrtA, sqrtB, bigOrZero(amount1))
}

func GetAmountsForLiquidity(
	sqrtPriceX96, sqrtRatioAX96, sqrtRatioBX96 *big.Int,
	liquidity *big.Int,
) (*big.Int, *big.Int) {
	sqrtA, sqrtB := sortRatios(sqrtRatioAX96, sqrtRatioBX96)
	price := bigOrZero(sqrtPriceX96)
	liq := bigOrZero(liquidity)

	if price.Cmp(sqrtA) <= 0 {
		return getAmount0ForLiquidity(sqrtA, sqrtB, liq), big.NewInt(0)
	}
	if price.Cmp(sqrtB) < 0 {
		return getAmount0ForLiquidity(price, sqrtB, liq),
			getAmount1ForLiquidity(sqrtA, price, liq)
	}
	return big.NewInt(0), getAmount1ForLiquidity(sqrtA, sqrtB, liq)
}

func GetLiquidityForAmount0(sqrtPriceAX96, sqrtPriceBX96, amount0 *big.Int) *big.Int {
	a, b := sortRatios(sqrtPriceAX96, sqrtPriceBX96)
	return getLiquidityForAmount0(a, b, bigOrZero(amount0))
}

func GetLiquidityForAmount1(sqrtPriceAX96, sqrtPriceBX96, amount1 *big.Int) *big.Int {
	a, b := sortRatios(sqrtPriceAX96, sqrtPriceBX96)
	return getLiquidityForAmount1(a, b, bigOrZero(amount1))
}

func getLiquidityForAmount0(sqrtA, sqrtB, amount0 *big.Int) *big.Int {
	sqrtA, sqrtB = sortRatios(sqrtA, sqrtB)
	intermediate := mulDiv(sqrtA, sqrtB, q96)
	return mulDiv(amount0, intermediate, new(big.Int).Sub(sqrtB, sqrtA))
}

func getLiquidityForAmount1(sqrtA, sqrtB, amount1 *big.Int) *big.Int {
	sqrtA, sqrtB = sortRatios(sqrtA, sqrtB)
	return mulDiv(amount1, q96, new(big.Int).Sub(sqrtB, sqrtA))
}

func getAmount0ForLiquidity(sqrtA, sqrtB, liquidity *big.Int) *big.Int {
	sqrtA, sqrtB = sortRatios(sqrtA, sqrtB)
	numerator := mulDiv(liquidity, new(big.Int).Sub(sqrtB, sqrtA), sqrtB)
	return mulDiv(numerator, q96, sqrtA)
}

func getAmount1ForLiquidity(sqrtA, sqrtB, liquidity *big.Int) *big.Int {
	sqrtA, sqrtB = sortRatios(sqrtA, sqrtB)
	return mulDiv(liquidity, new(big.Int).Sub(sqrtB, sqrtA), q96)
}

func mulDiv(a, b, denom *big.Int) *big.Int {
	return new(big.Int).Div(new(big.Int).Mul(a, b), denom)
}

func sortRatios(a, b *big.Int) (*big.Int, *big.Int) {
	a, b = bigOrZero(a), bigOrZero(b)
	if a.Cmp(b) > 0 {
		return b, a
	}
	return a, b
}

func bigOrZero(v *big.Int) *big.Int {
	if v == nil {
		return big.NewInt(0)
	}
	return v
}
