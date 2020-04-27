package bls

import (
	"encoding/hex"
	"fmt"
	"math/big"
)

type fe /***			***/ [6]uint64
type fe2 /**			***/ [2]fe
type fe6 /**			***/ [3]fe2
type fe12 /**			***/ [2]fe6

func (fe *fe) Bytes() []byte {
	out := make([]byte, 48)
	var a int
	for i := 0; i < 6; i++ {
		a = 48 - i*8
		out[a-1] = byte(fe[i])
		out[a-2] = byte(fe[i] >> 8)
		out[a-3] = byte(fe[i] >> 16)
		out[a-4] = byte(fe[i] >> 24)
		out[a-5] = byte(fe[i] >> 32)
		out[a-6] = byte(fe[i] >> 40)
		out[a-7] = byte(fe[i] >> 48)
		out[a-8] = byte(fe[i] >> 56)
	}
	return out
}

func (fe *fe) FromBytes(in []byte) *fe {
	size := 48
	l := len(in)
	if l >= size {
		l = size
	}
	padded := make([]byte, size)
	copy(padded[size-l:], in[:])
	var a int
	for i := 0; i < 6; i++ {
		a = size - i*8
		fe[i] = uint64(padded[a-1]) | uint64(padded[a-2])<<8 |
			uint64(padded[a-3])<<16 | uint64(padded[a-4])<<24 |
			uint64(padded[a-5])<<32 | uint64(padded[a-6])<<40 |
			uint64(padded[a-7])<<48 | uint64(padded[a-8])<<56
	}
	return fe
}

func (fe *fe) SetBig(a *big.Int) *fe {
	return fe.FromBytes(a.Bytes())
}

func (fe *fe) SetUint(a uint64) *fe {
	fe[0] = a
	fe[1] = 0
	fe[2] = 0
	fe[3] = 0
	fe[4] = 0
	fe[5] = 0
	return fe
}

func (fe *fe) SetString(s string) (*fe, error) {
	if s[:2] == "0x" {
		s = s[2:]
	}
	bytes, err := hex.DecodeString(s)
	if err != nil {
		return nil, err
	}
	return fe.FromBytes(bytes), nil
}

func (fe *fe) Set(fe2 *fe) *fe {
	fe[0] = fe2[0]
	fe[1] = fe2[1]
	fe[2] = fe2[2]
	fe[3] = fe2[3]
	fe[4] = fe2[4]
	fe[5] = fe2[5]
	return fe
}

func (fe *fe) SetZero() *fe {
	fe[0] = 0
	fe[1] = 0
	fe[2] = 0
	fe[3] = 0
	fe[4] = 0
	fe[5] = 0
	return fe
}

func (fe *fe) SetOne() *fe {
	fe.Set(r1)
	return fe
}

func (fe *fe) Big() *big.Int {
	return new(big.Int).SetBytes(fe.Bytes())
}

func (fe fe) String() (s string) {
	for i := 5; i >= 0; i-- {
		s = fmt.Sprintf("%s%16.16x", s, fe[i])
	}
	return "0x" + s
}

func (fe *fe) IsOdd() bool {
	var mask uint64 = 1
	return fe[0]&mask != 0
}

func (fe *fe) IsEven() bool {
	var mask uint64 = 1
	return fe[0]&mask == 0
}

func (fe *fe) IsZero() bool {
	return (fe[5] | fe[4] | fe[3] | fe[2] | fe[1] | fe[0]) == 0
}

func (fe *fe) IsOne() bool {
	return 1 == fe[0] && 0 == fe[1] && 0 == fe[2] && 0 == fe[3] && 0 == fe[4] && 0 == fe[5]
}

func (fe *fe) Cmp(fe2 *fe) int {
	for i := 5; i > -1; i-- {
		if fe[i] > fe2[i] {
			return 1
		} else if fe[i] < fe2[i] {
			return -1
		}
	}
	return 0
}

func (fe *fe) Equals(fe2 *fe) bool {
	return fe2[0] == fe[0] && fe2[1] == fe[1] && fe2[2] == fe[2] && fe2[3] == fe[3] && fe2[4] == fe[4] && fe2[5] == fe[5]
}

func (e *fe) signBE() bool {
	negZ, z := new(fe), new(fe)
	fromMont(z, e)
	neg(negZ, z)
	return negZ.Cmp(z) > -1
}

func (e *fe) sign() bool {
	r := new(fe)
	fromMont(r, e)
	return r[0]&1 == 0
}

func (fe *fe) div2(e uint64) {
	fe[0] = fe[0]>>1 | fe[1]<<63
	fe[1] = fe[1]>>1 | fe[2]<<63
	fe[2] = fe[2]>>1 | fe[3]<<63
	fe[3] = fe[3]>>1 | fe[4]<<63
	fe[4] = fe[4]>>1 | fe[5]<<63
	fe[5] = fe[5]>>1 | e<<63
}

func (fe *fe) mul2() uint64 {
	e := fe[5] >> 63
	fe[5] = fe[5]<<1 | fe[4]>>63
	fe[4] = fe[4]<<1 | fe[3]>>63
	fe[3] = fe[3]<<1 | fe[2]>>63
	fe[2] = fe[2]<<1 | fe[1]>>63
	fe[1] = fe[1]<<1 | fe[0]>>63
	fe[0] = fe[0] << 1
	return e
}

func (fe *fe2) signBE() bool {
	if !fe[1].IsZero() {
		return fe[1].signBE()
	}
	return fe[0].signBE()
}

func (e *fe2) sign() bool {
	r := new(fe)
	if !e[0].IsZero() {
		fromMont(r, &e[0])
		return r[0]&1 == 0
	}
	fromMont(r, &e[1])
	return r[0]&1 == 0
}
