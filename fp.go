package bls

import (
	"crypto/rand"
	"fmt"
	"io"
	"math/big"
)

type Fp struct {
}

var fpOne = *r1
var fpZero = Fe{0, 0, 0, 0, 0, 0}

func NewFp() *Fp {
	return &Fp{}
}

func (f *Fp) NewElementFromBytes(fe *Fe, in []byte) error {
	fe.FromBytes(in)
	if !f.Valid(fe) {
		return fmt.Errorf("invalid input string")
	}
	f.Mul(fe, fe, r2)
	return nil
}

func (f *Fp) NewElementFromUint(in uint64) (*Fe, error) {
	fe := &Fe{in}
	if in == 0 {
		return fe, nil
	}
	if !f.Valid(fe) {
		return nil, fmt.Errorf("invalid input string")
	}
	f.Mul(fe, fe, r2)
	return fe, nil
}

func (f *Fp) NewElementFromBig(in *big.Int) (*Fe, error) {
	fe := new(Fe).SetBig(in)
	if !f.Valid(fe) {
		return nil, fmt.Errorf("invalid input string")
	}
	f.Mul(fe, fe, r2)
	return fe, nil
}

func (f *Fp) NewElementFromString(in string) (*Fe, error) {
	fe, err := new(Fe).SetString(in)
	if err != nil {
		return nil, err
	}
	if !f.Valid(fe) {
		return nil, fmt.Errorf("invalid input string")
	}
	f.Mul(fe, fe, r2)
	return fe, nil
}

func (f *Fp) ToBytes(fe *Fe) []byte {
	fe2 := new(Fe)
	f.Demont(fe2, fe)
	return fe2.Bytes()
}

func (f *Fp) ToBig(fe *Fe) *big.Int {
	fe2 := new(Fe)
	f.Demont(fe2, fe)
	return fe2.Big()
}

func (f *Fp) ToString(fe *Fe) (s string) {
	fe2 := new(Fe)
	f.Demont(fe2, fe)
	return fe2.String()
}

func (f *Fp) Valid(fe *Fe) bool {
	return fe.Cmp(&modulus) == -1
}

func (f *Fp) Zero() *Fe {
	return new(Fe).SetUint(0)
}

func (f *Fp) One() *Fe {
	return new(Fe).Set(r1)
}

func (f *Fp) Copy(dst *Fe, src *Fe) *Fe {
	return dst.Set(src)
}

func (f *Fp) RandElement(fe *Fe, r io.Reader) (*Fe, error) {
	bi, err := rand.Int(r, modulus.Big())
	if err != nil {
		return nil, err
	}
	return fe.SetBig(bi), nil
}

func (f *Fp) Equal(a, b *Fe) bool {
	return a.Equals(b)
}

func (f *Fp) IsZero(a *Fe) bool {
	return a.IsZero()
}

func (f *Fp) Mont(c, a *Fe) {
	montmul(c, a, r2)
}

func (f *Fp) Demont(c, a *Fe) {
	montmul(c, a, &Fe{1})
}

func (f *Fp) Add(c, a, b *Fe) {
	add(c, a, b)
}

func (f *Fp) Double(c, a *Fe) {
	double(c, a)
}

func (f *Fp) Sub(c, a, b *Fe) {
	sub(c, a, b)
}

func (f *Fp) Neg(c, a *Fe) {
	if a.IsZero() {
		c.Set(a)
	} else {
		neg(c, a)
	}
}

func (f *Fp) Square(c, a *Fe) {
	montsquare(c, a)
}

func (f *Fp) Mul(c, a, b *Fe) {
	montmul(c, a, b)
}

func (f *Fp) Exp(c, a *Fe, e *big.Int) {
	z := new(Fe).Set(r1)
	for i := e.BitLen(); i >= 0; i-- {
		montsquare(z, z)
		if e.Bit(i) == 1 {
			montmul(z, z, a)
		}
	}
	c.Set(z)
}

func (f *Fp) Inverse(inv, fe *Fe) {
	f.InvMontUp(inv, fe)
}

func (f *Fp) InvMontUp(inv, fe *Fe) {
	u := new(Fe).Set(&modulus)
	v := new(Fe).Set(fe)
	s := &Fe{1}
	r := &Fe{0}
	var k int
	var z uint64
	var found = false
	// Phase 1
	for i := 0; i < 384*2; i++ {
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
			subn(u, v)
			u.div2(0)
			addn(r, s)
			s.mul2()
		} else {
			subn(v, u)
			v.div2(0)
			addn(s, r)
			z += r.mul2()
		}
		k += 1
	}
	if found && k > 384 {
		if r.Cmp(&modulus) != -1 || z > 0 {
			subn(r, &modulus)
		}
		u.Set(&modulus)
		subn(u, r)
		// Phase 2
		for i := k; i < 384*2; i++ {
			double(u, u)
		}
		inv.Set(u)
	} else {
		inv.Set(&Fe{0})
	}
}

func (f *Fp) InvMontDown(inv, fe *Fe) {
	u := new(Fe).Set(&modulus)
	v := new(Fe).Set(fe)
	s := &Fe{1}
	r := &Fe{0}
	var k int
	var z uint64
	var found = false
	// Phase 1
	for i := 0; i < 384*2; i++ {
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
			subn(u, v)
			u.div2(0)
			addn(r, s)
			s.mul2()
		} else {
			subn(v, u)
			v.div2(0)
			addn(s, r)
			z += r.mul2()
		}
		k += 1
	}
	if found && k > 384 {
		if r.Cmp(&modulus) != -1 || z > 0 {
			subn(r, &modulus)
		}
		u.Set(&modulus)
		subn(u, r)
		// Phase 2
		var e uint64
		for i := 0; i < k-384; i++ {
			if u.IsEven() {
				u.div2(0)
			} else {
				e = addn(u, &modulus)
				u.div2(e)
			}
		}
		inv.Set(u)
	} else {
		inv.Set(&Fe{0})
	}
}

func (f *Fp) InvEEA(inv, fe *Fe) {
	u := new(Fe).Set(fe)
	v := new(Fe).Set(&modulus)
	x1 := &Fe{1}
	x2 := &Fe{0}
	var e uint64
	for !u.IsOne() && !v.IsOne() {
		for u.IsEven() {
			u.div2(0)
			if x1.IsEven() {
				x1.div2(0)
			} else {
				e = addn(x1, &modulus)
				x1.div2(e)
			}
		}
		for v.IsEven() {
			v.div2(0)
			if x2.IsEven() {
				x2.div2(0)
			} else {
				e = addn(x2, &modulus)
				x2.div2(e)
			}
		}
		if u.Cmp(v) == -1 {
			subn(v, u)
			sub(x2, x2, x1)
		} else {
			subn(u, v)
			sub(x1, x1, x2)
		}
	}
	if u.IsOne() {
		inv.Set(x1)
		return
	}
	inv.Set(x2)
}

func (f *Fp) Sqrt(c, a *Fe) (hasRoot bool) {
	var u, v Fe
	f.Copy(&u, a)
	f.Exp(c, a, pPlus1Over4)
	f.Square(&v, c)
	return f.Equal(&u, &v)
}
