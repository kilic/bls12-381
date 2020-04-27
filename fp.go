package bls

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

func fromBytes(in []byte) (*fe, error) {
	fe := &fe{}
	if len(in) != 48 {
		return nil, fmt.Errorf("input string should be equal 48 bytes")
	}
	fe.FromBytes(in)
	if !valid(fe) {
		return nil, fmt.Errorf("must be less than modulus")
	}
	mul(fe, fe, r2)
	return fe, nil
}

func from64Bytes(in []byte) (*fe, error) {
	if len(in) != 64 {
		return nil, fmt.Errorf("input string should be equal 64 bytes")
	}
	a0 := make([]byte, 48)
	copy(a0[16:48], in[:32])
	a1 := make([]byte, 48)
	copy(a1[16:48], in[32:])
	e0, err := fromBytes(a0)
	if err != nil {
		return nil, err
	}
	e1, err := fromBytes(a1)
	if err != nil {
		return nil, err
	}
	// F = 2 ^ 256 * R
	F := fe{
		0x75b3cd7c5ce820f,
		0x3ec6ba621c3edb0b,
		0x168a13d82bff6bce,
		0x87663c4bf8c449d2,
		0x15f34c83ddc8d830,
		0xf9628b49caa2e85,
	}

	mul(e0, e0, &F)
	add(e1, e1, e0)
	return e1, nil
}

func fromBig(in *big.Int) (*fe, error) {
	fe := new(fe).SetBig(in)
	if !valid(fe) {
		return nil, fmt.Errorf("invalid input string")
	}
	mul(fe, fe, r2)
	return fe, nil
}

func fromString(in string) (*fe, error) {
	fe, err := new(fe).SetString(in)
	if err != nil {
		return nil, err
	}
	if !valid(fe) {
		return nil, fmt.Errorf("invalid input string")
	}
	mul(fe, fe, r2)
	return fe, nil
}

func toBytes(e *fe) []byte {
	e2 := new(fe)
	fromMont(e2, e)
	return e2.Bytes()
}

func toBig(e *fe) *big.Int {
	e2 := new(fe)
	fromMont(e2, e)
	return e2.Big()
}

func toString(e *fe) (s string) {
	e2 := new(fe)
	fromMont(e2, e)
	return e2.String()
}

func valid(fe *fe) bool {
	return fe.Cmp(&modulus) == -1
}

func zero() *fe {
	return &fe{}
}

func one() *fe {
	return new(fe).Set(r1)
}

func newRand(r io.Reader) (*fe, error) {
	fe := new(fe)
	bi, err := rand.Int(r, modulus.Big())
	if err != nil {
		return nil, err
	}
	return fe.SetBig(bi), nil
}

func equal(a, b *fe) bool {
	return a.Equals(b)
}

func isZero(a *fe) bool {
	return a.IsZero()
}

func isOne(a *fe) bool {
	return a.Equals(one())
}

func toMont(c, a *fe) {
	mul(c, a, r2)
}

func fromMont(c, a *fe) {
	mul(c, a, &fe{1})
}

func exp(c, a *fe, e *big.Int) {
	z := new(fe).Set(r1)
	for i := e.BitLen(); i >= 0; i-- {
		mul(z, z, z)
		if e.Bit(i) == 1 {
			mul(z, z, a)
		}
	}
	c.Set(z)
}

func inverse(inv, e *fe) {
	if e.IsZero() {
		inv.SetZero()
		return
	}
	u := new(fe).Set(&modulus)
	v := new(fe).Set(e)
	s := &fe{1}
	r := &fe{0}
	var k int
	var z uint64
	var found = false
	// Phase 1
	for i := 0; i < 768; i++ {
		if v.IsZero() {
			found = true
			break
		}
		if u.IsEven() {
			u.div2(0)
			s.mul2()
		} else if v.IsEven() {
			v.div2(0)
			z += r.mul2()
		} else if u.Cmp(v) == 1 {
			lsubAssign(u, v)
			u.div2(0)
			laddAssign(r, s)
			s.mul2()
		} else {
			lsubAssign(v, u)
			v.div2(0)
			laddAssign(s, r)
			z += r.mul2()
		}
		k += 1
	}

	if !found {
		inv.SetZero()
		return
	}

	if k < 381 || k > 381+384 {
		inv.SetZero()
		return
	}

	if r.Cmp(&modulus) != -1 || z > 0 {
		lsubAssign(r, &modulus)
	}
	u.Set(&modulus)
	lsubAssign(u, r)

	// Phase 2
	for i := k; i < 384*2; i++ {
		double(u, u)
	}
	inv.Set(u)
	return
}

func sqrt(c, a *fe) (hasRoot bool) {
	u, v := new(fe).Set(a), new(fe)
	exp(c, a, pPlus1Over4)
	square(v, c)
	return equal(u, v)
}

func isQuadraticNonResidue(elem *fe) bool {
	result := new(fe)
	exp(result, elem, pMinus1Over2)
	return !equal(result, one())
}
