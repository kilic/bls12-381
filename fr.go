package bls12381

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

const frByteSize = 32
const frBitSize = 255
const frNumberOfLimbs = 4

// Fr is scalar field element representation
type Fr [4]uint64

func NewFr() *Fr {
	return &Fr{}
}

func (e *Fr) Rand(r io.Reader) (*Fr, error) {
	bi, err := rand.Int(r, qBig)
	if err != nil {
		return nil, err
	}
	_ = e.setBig(bi)
	return e, nil
}

func (e *Fr) Set(e2 *Fr) *Fr {
	e[0] = e2[0]
	e[1] = e2[1]
	e[2] = e2[2]
	e[3] = e2[3]
	return e
}

func (e *Fr) Zero() *Fr {
	e[0] = 0
	e[1] = 0
	e[2] = 0
	e[3] = 0
	return e
}

func (e *Fr) One() *Fr {
	e.Set(&Fr{1})
	return e
}

func (e *Fr) RedOne() *Fr {
	e.Set(sr1)
	return e
}

func (e *Fr) FromBytes(in []byte) *Fr {
	e.setBytes(in)
	return e
}

func (e *Fr) RedFromBytes(in []byte) *Fr {
	e.setBytes(in)
	e.toMont()
	return e
}

func (e *Fr) setBytes(in []byte) {
	u := new(big.Int).SetBytes(in)
	_ = e.setBig(u)
}

func (e *Fr) setBig(in *big.Int) error {
	zero := new(big.Int)
	c := in.Cmp(zero)
	if c == -1 {
		return errors.New("cannot set negative element")
	} else if c == 0 {
		e.Zero()
		return nil
	}

	c = in.Cmp(qBig)
	if c == 0 {
		e.Zero()
		return nil
	} else if c == -1 {
		words := in.Bits()
		e[0] = uint64(words[0])
		e[1] = uint64(words[1])
		e[2] = uint64(words[2])
		e[3] = uint64(words[3])
		return nil
	}
	_in := new(big.Int).Mod(in, qBig)
	words := _in.Bits()
	e[0] = uint64(words[0])
	e[1] = uint64(words[1])
	e[2] = uint64(words[2])
	e[3] = uint64(words[3])
	return nil
}

func (e *Fr) ToBytes() []byte {
	return NewFr().Set(e).bytes()
}

func (e *Fr) RedToBytes() []byte {
	out := NewFr().Set(e)
	out.fromMont()
	return out.bytes()
}

func (e *Fr) ToBig() *big.Int {
	return new(big.Int).SetBytes(e.ToBytes())
}

func (e *Fr) RedToBig() *big.Int {
	return new(big.Int).SetBytes(e.RedToBytes())
}

func (e *Fr) bytes() []byte {
	out := make([]byte, frByteSize)
	var a int
	for i := 0; i < frNumberOfLimbs; i++ {
		a = frByteSize - i*8
		out[a-1] = byte(e[i])
		out[a-2] = byte(e[i] >> 8)
		out[a-3] = byte(e[i] >> 16)
		out[a-4] = byte(e[i] >> 24)
		out[a-5] = byte(e[i] >> 32)
		out[a-6] = byte(e[i] >> 40)
		out[a-7] = byte(e[i] >> 48)
		out[a-8] = byte(e[i] >> 56)
	}
	return out
}

func (e *Fr) isZero() bool {
	return (e[3] | e[2] | e[1] | e[0]) == 0
}

func (e *Fr) isOne() bool {
	return e.equal(&Fr{1})
}

func (e *Fr) isRedOne() bool {
	return e.equal(sr1)
}

func (e *Fr) equal(e2 *Fr) bool {
	return e2[0] == e[0] && e2[1] == e[1] && e2[2] == e[2] && e2[3] == e[3]
}

func (e *Fr) sliceUint64(from int) uint64 {
	if from < 64 {
		return e[0]>>from | e[1]<<(64-from)
	} else if from < 128 {
		return e[1]>>(from-64) | e[2]<<(128-from)
	} else if from < 192 {
		return e[2]>>(from-128) | e[3]<<(192-from)
	}
	return e[3] >> (from - 192)
}

func (e *Fr) Bit(at int) bool {
	if at < 64 {
		return (e[0]>>at)&1 == 1
	} else if at < 128 {
		return (e[1]>>(at-64))&1 == 1
	} else if at < 192 {
		return (e[2]>>(at-128))&1 == 1
	} else if at < 256 {
		return (e[3]>>(at-192))&1 == 1
	}
	return false
}

func (e *Fr) toMont() {
	e.RedMul(e, sr2)
}

func (e *Fr) fromMont() {
	e.RedMul(e, &Fr{1})
}

func (e *Fr) Add(a, b *Fr) {
	addFR(e, a, b)
}

func (e *Fr) Double(a *Fr) {
	doubleFR(e, a)
}

func (e *Fr) Sub(a, b *Fr) {
	subFR(e, a, b)
}

func (e *Fr) Mul(a, b *Fr) {
	mulFR(e, a, b)
	mulFR(e, e, sr2)
}

func (e *Fr) RedMul(a, b *Fr) {
	mulFR(e, a, b)
}

func (e *Fr) Square(a *Fr) {
	squareFR(e, a)
	mulFR(e, e, sr2)
}

func (e *Fr) RedSquare(a *Fr) {
	squareFR(e, a)
}

func (e *Fr) Neg(a *Fr) {
	negFR(e, a)
}

func (e *Fr) Exp(a *Fr, ee *big.Int) {
	z := new(Fr).RedOne()
	for i := ee.BitLen(); i >= 0; i-- {
		z.RedSquare(z)
		if ee.Bit(i) == 1 {
			z.RedMul(z, a)
		}
	}
	e.Set(z)
}

// func toWNAF(e *big.Int, w uint) nafNumber {
// 	z := new(big.Int).Set(e)
// 	naf := []nafSign{}
// 	W := new(big.Int).Lsh(bigOne, w)
// 	Wl := new(big.Int).Rsh(W, 1)
// 	for z.Cmp(bigZero) != 0 {
// 		if z.Bit(0) == 1 {
// 			nafBit := new(big.Int)
// 			nafBit.Mod(z, W)
// 			if nafBit.Cmp(Wl) >= 0 {
// 				nafBit.Sub(nafBit, W)
// 			}
// 			naf = append(naf, nafSign(nafBit.Int64()))
// 			z.Sub(z, nafBit)
// 		} else {
// 			naf = append(naf, nafZERO)
// 		}
// 		z.Rsh(z, 1)
// 	}
// 	return naf
// }

// func fromWNAF(naf nafNumber, w int) *big.Int {
// 	acc := new(big.Int)
// 	d := big.NewInt(1)
// 	d.Lsh(d, uint(len(naf)-1))
// 	for i := len(naf) - 1; i >= 0; i-- {
// 		if naf[i] == nafNEG {
// 			acc.Sub(acc, d)
// 		} else if naf[i] == nafPOS {
// 			acc.Add(acc, d)
// 		}
// 		d.Rsh(d, 1)
// 	}
// 	return acc
// }
