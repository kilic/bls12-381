package bls12381

import (
	"math/bits"
)

type scalar []uint64

func (repr scalar) new(x uint64) scalar {
	n := make(scalar, len(repr))
	n[0] = x
	return n
}

func (repr scalar) isZero() bool {
	for _, i := range repr {
		if i != 0 {
			return false
		}
	}
	return true
}

func (repr scalar) isOdd() bool {
	if len(repr) == 0 {
		return false
	}
	return repr[0]&1 == 1
}

func (repr scalar) div2() {
	var t uint64
	for i := len(repr) - 1; i >= 0; i-- {
		t2 := repr[i] << 63
		repr[i] >>= 1
		repr[i] |= t
		t = t2
	}
}

func (repr scalar) sbb(b scalar) {
	var borrow uint64
	for i := 0; i < len(repr); i++ {
		repr[i], borrow = bits.Sub64(repr[i], b[i], borrow)
	}
}
func (repr scalar) adc(b scalar) {
	var carry uint64
	for i := 0; i < len(repr); i++ {
		repr[i], carry = bits.Add64(repr[i], b[i], carry)
	}
}

func wnaf(e scalar, window uint) []int64 {
	if e.isZero() {
		return []int64{}
	}
	max := int64(1 << window)
	midpoint := int64(1 << (window - 1))
	modulusMask := uint64(1<<window) - 1

	var b scalar
	var out []int64
	for !e.isZero() {
		var z int64
		if e.isOdd() {
			maskedBits := int64(e[0] & modulusMask)
			if maskedBits > midpoint {
				z = maskedBits - max
				b = e.new(uint64(0 - z))
				e.adc(b)
			} else {
				z = maskedBits
				b := e.new(uint64(z))
				e.sbb(b)
			}
		} else {
			z = 0
		}
		out = append(out, z)
		e.div2()
	}
	return out
}
