package math

import (
	"math/big"
)

const maxTick = 887272

// GetSqrtRatioAtTick returns sqrt(1.0001^tick) * 2^96 (Uniswap TickMath).
func GetSqrtRatioAtTick(tick int32) *big.Int {
	var absTick uint32
	if tick < 0 {
		absTick = uint32(-tick)
	} else {
		absTick = uint32(tick)
	}
	if absTick > maxTick {
		panic("tick out of range")
	}

	ratio := new(big.Int)
	if absTick&0x1 != 0 {
		ratio.SetString("fffcb933bd6fad37aa2d162d1a594001", 16)
	} else {
		ratio.SetString("100000000000000000000000000000000", 16)
	}

	multShift := func(hex string) {
		m := new(big.Int)
		m.SetString(hex, 16)
		ratio.Mul(ratio, m)
		ratio.Rsh(ratio, 128)
	}

	if absTick&0x2 != 0 {
		multShift("fff97272373d413259a46990580e213a")
	}
	if absTick&0x4 != 0 {
		multShift("fff2e50f5f656932ef12357cf3c7fdcc")
	}
	if absTick&0x8 != 0 {
		multShift("ffe5caca7e10e4e61c3624eaa0941cd0")
	}
	if absTick&0x10 != 0 {
		multShift("ffcb9843d60f6159c9db58835c926644")
	}
	if absTick&0x20 != 0 {
		multShift("ff973b41fa98c081472e6896dfb254c0")
	}
	if absTick&0x40 != 0 {
		multShift("ff2ea16466c96a3843ec78b326b52861")
	}
	if absTick&0x80 != 0 {
		multShift("fe5dee046a99a2a811c461f1969c3053")
	}
	if absTick&0x100 != 0 {
		multShift("fcbe86c7900a88aedcffc83b479aa3a4")
	}
	if absTick&0x200 != 0 {
		multShift("f987a7253ac413176f2b074cf7815e54")
	}
	if absTick&0x400 != 0 {
		multShift("f3392b0822b70005940c7a398e4b70f3")
	}
	if absTick&0x800 != 0 {
		multShift("e7159475a2c29b7443b29c7fa6e889d9")
	}
	if absTick&0x1000 != 0 {
		multShift("d097f3bdfd2022b8845ad8f792aa5825")
	}
	if absTick&0x2000 != 0 {
		multShift("a9f746462d870fdf8a65dc1f90e061e5")
	}
	if absTick&0x4000 != 0 {
		multShift("70d869a156d2a1b890bb3df62baf32f7")
	}
	if absTick&0x8000 != 0 {
		multShift("31be135f97d08fd981231505542fcfa6")
	}
	if absTick&0x10000 != 0 {
		multShift("9aa508b5b7a84e1c677de54f3e99bc9")
	}
	if absTick&0x20000 != 0 {
		multShift("5d6af8dedb81196699c329225ee604")
	}
	if absTick&0x40000 != 0 {
		multShift("2216e584f5fa1ea926041bedfe98")
	}
	if absTick&0x80000 != 0 {
		multShift("48a170391f7dc42444e8fa2")
	}

	if tick > 0 {
		maxU256 := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
		ratio.Div(maxU256, ratio)
	}

	shift32 := new(big.Int).Lsh(big.NewInt(1), 32)
	result := new(big.Int).Rsh(ratio, 32)
	if new(big.Int).Mod(ratio, shift32).Sign() != 0 {
		result.Add(result, big.NewInt(1))
	}
	return result
}
