package math

import "math/big"

// FloorTick aligns tick down to tickSpacing (matches UniswapV3Strategy._floor).
func FloorTick(tick, spacing int32) int32 {
	if spacing == 0 {
		return tick
	}
	compressed := tick / spacing
	if tick < 0 && tick%spacing != 0 {
		compressed--
	}
	return compressed * spacing
}

// ComputeTickRange returns the liquidity range centered on currentTick.
func ComputeTickRange(currentTick, halfRangeTicks, tickSpacing int32) (tickLower, tickUpper int32) {
	center := FloorTick(currentTick, tickSpacing)
	tickLower = FloorTick(center-halfRangeTicks, tickSpacing)
	tickUpper = FloorTick(center+halfRangeTicks, tickSpacing)
	if tickLower >= tickUpper {
		tickUpper = tickLower + tickSpacing
	}
	return tickLower, tickUpper
}

// NeedRebalance mirrors UniswapV3Strategy._needRebalance.
func NeedRebalance(currentTick, tickLower, tickUpper int32, hasPosition bool) bool {
	if !hasPosition {
		return true
	}
	return currentTick < tickLower || currentTick > tickUpper
}

// Int24FromBigInt converts abigen int24 (*big.Int) to int32.
func Int24FromBigInt(v *big.Int) int32 {
	if v == nil {
		return 0
	}
	// int24 two's complement range
	mod := new(big.Int).Lsh(big.NewInt(1), 24)
	x := new(big.Int).Mod(v, mod)
	if x.Cmp(new(big.Int).Lsh(big.NewInt(1), 23)) >= 0 {
		x.Sub(x, mod)
	}
	return int32(x.Int64())
}
