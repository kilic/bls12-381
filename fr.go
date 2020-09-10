package bls12381

import (
	"crypto/rand"
	"errors"
	"io"
	"math/big"
)

const frByteSize = 32
const frNumberOfLimbs = 4

// Fr is scalar field element representation
type Fr [4]uint64

func NewFr() *Fr {
	return &Fr{}
}

func (e *Fr) Set(e2 *Fr) *Fr {
	e[0] = e2[0]
	e[1] = e2[1]
	e[2] = e2[2]
	e[3] = e2[3]
	return e
}

func (e *Fr) zero() *Fr {
	e[0] = 0
	e[1] = 0
	e[2] = 0
	e[3] = 0
	return e
}

func (e *Fr) one() *Fr {
	e.Set(sr1)
	return e
}

func (e *Fr) FromBytes(in []byte) *Fr {
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
		e.zero()
		return nil
	}

	c = in.Cmp(q)
	if c == 0 {
		e.zero()
		return nil
	} else if c == -1 {
		words := in.Bits()
		e[0] = uint64(words[0])
		e[1] = uint64(words[1])
		e[2] = uint64(words[2])
		e[3] = uint64(words[3])
		return nil
	}
	_in := new(big.Int).Mod(in, q)
	words := _in.Bits()
	e[0] = uint64(words[0])
	e[1] = uint64(words[1])
	e[2] = uint64(words[2])
	e[3] = uint64(words[3])
	return nil
}

func (e *Fr) ToBytes() []byte {
	out := NewFr()
	out.Set(e)
	out.fromMont()
	return out.bytes()
}

func (e *Fr) Big() *big.Int {
	return new(big.Int).SetBytes(e.ToBytes())
}

func (e *Fr) big() *big.Int {
	return new(big.Int).SetBytes(e.bytes())
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
	return e.equal(sr1)
}

func (e *Fr) equal(e2 *Fr) bool {
	return e2[0] == e[0] && e2[1] == e[1] && e2[2] == e[2] && e2[3] == e[3]
}

func (e *Fr) Rand(r io.Reader) (*Fr, error) {
	bi, err := rand.Int(r, q)
	if err != nil {
		return nil, err
	}
	_ = e.setBig(bi)
	return e, nil
}

func (e *Fr) toMont() {
	e.Mul(e, sr2)
}

func (e *Fr) fromMont() {
	e.Mul(e, &Fr{1})
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
}

func (e *Fr) Square(a *Fr) {
	squareFR(e, a)
}

func (e *Fr) Neg(a *Fr) {
	negFR(e, a)
}

func (e *Fr) Exp(a *Fr, ee *big.Int) {
	z := new(Fr).one()
	for i := ee.BitLen(); i >= 0; i-- {
		z.Square(z)
		if ee.Bit(i) == 1 {
			z.Mul(z, a)
		}
	}
	e.Set(z)
}
