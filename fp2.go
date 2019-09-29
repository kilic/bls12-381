package bls

import (
	"fmt"
	"io"
	"math/big"
)

type Fe2 [2]Fe

type Fp2 struct {
	f *Fp
	t [4]*Fe
}

var Fp2One = Fe2{FpOne, FpZero}
var Fp2Zero = Fe2{FpZero, FpZero}

func NewFp2(f *Fp) *Fp2 {
	t := [4]*Fe{}
	for i := 0; i < 4; i++ {
		t[i] = f.Zero()
	}
	if f == nil {
		return &Fp2{NewFp(), t}
	}
	return &Fp2{f, t}
}

func (fp *Fp2) NewElement() *Fe2 {
	return &Fe2{}
}

func (fp *Fp2) NewElementFromBytes(fe *Fe2, b []byte) error {
	if len(b) < 96 {
		return fmt.Errorf("input string should be larger than 96 bytes")
	}
	if err := fp.f.NewElementFromBytes(&fe[1], b[:48]); err != nil {
		return err
	}
	if err := fp.f.NewElementFromBytes(&fe[0], b[48:]); err != nil {
		return err
	}
	return nil
}

func (fp *Fp2) RandElement(a *Fe2, r io.Reader) (*Fe2, error) {
	if _, err := fp.f.RandElement(&a[0], r); err != nil {
		return nil, err
	}
	if _, err := fp.f.RandElement(&a[1], r); err != nil {
		return nil, err
	}
	return a, nil
}

func (fp *Fp2) Zero() *Fe2 {
	return &Fe2{}
}

func (fp *Fp2) One() *Fe2 {
	return &Fe2{*fp.f.One(), *fp.f.Zero()}
}

func (fp *Fp2) ToBytes(a *Fe2) []byte {
	out := make([]byte, 96)
	copy(out[:48], fp.f.ToBytes(&a[1]))
	copy(out[48:], fp.f.ToBytes(&a[0]))
	return out
}

func (fp *Fp2) IsZero(a *Fe2) bool {
	return fp.f.IsZero(&a[0]) && fp.f.IsZero(&a[1])
}

func (fp *Fp2) Equal(a, b *Fe2) bool {
	return fp.f.Equal(&a[0], &b[0]) && fp.f.Equal(&a[1], &b[1])
}

func (fp *Fp2) Copy(c, a *Fe2) *Fe2 {
	fp.f.Copy(&c[0], &a[0])
	fp.f.Copy(&c[1], &a[1])
	return c
}

func (fp *Fp2) MulByNonResidue(c, a *Fe2) {
	t := fp.t
	fp.f.Sub(t[0], &a[0], &a[1])
	fp.f.Add(&c[1], &a[0], &a[1])
	fp.f.Copy(&c[0], t[0])
}

func (fp *Fp2) Add(c, a, b *Fe2) {
	fp.f.Add(&c[0], &a[0], &b[0])
	fp.f.Add(&c[1], &a[1], &b[1])
}

func (fp *Fp2) Double(c, a *Fe2) {
	fp.f.Double(&c[0], &a[0])
	fp.f.Double(&c[1], &a[1])
}

func (fp *Fp2) Sub(c, a, b *Fe2) {
	fp.f.Sub(&c[0], &a[0], &b[0])
	fp.f.Sub(&c[1], &a[1], &b[1])
}

func (fp *Fp2) Neg(c, a *Fe2) {
	fp.f.Neg(&c[0], &a[0])
	fp.f.Neg(&c[1], &a[1])
}

func (fp *Fp2) Conjugate(c, a *Fe2) {
	fp.f.Copy(&c[0], &a[0])
	fp.f.Neg(&c[1], &a[1])
}

func (fp *Fp2) Mul(c, a, b *Fe2) {
	t := fp.t
	fp.f.Mul(t[1], &a[0], &b[0])
	fp.f.Mul(t[2], &a[1], &b[1])
	fp.f.Sub(t[0], t[1], t[2])
	fp.f.Add(t[1], t[1], t[2])
	fp.f.Add(t[2], &a[0], &a[1])
	fp.f.Add(t[3], &b[0], &b[1])
	fp.f.Copy(&c[0], t[0])
	fp.f.Mul(t[0], t[2], t[3])
	fp.f.Sub(&c[1], t[0], t[1])
}

func (fp *Fp2) Square(c, a *Fe2) {
	t := fp.t
	fp.f.Add(t[0], &a[0], &a[1])
	fp.f.Sub(t[1], &a[0], &a[1])
	fp.f.Double(t[2], &a[0])
	fp.f.Mul(&c[0], t[0], t[1])
	fp.f.Mul(&c[1], t[2], &a[1])
}

func (fp *Fp2) Inverse(c, a *Fe2) {
	t := fp.t
	fp.f.Square(t[0], &a[0])
	fp.f.Square(t[1], &a[1])
	fp.f.Add(t[0], t[0], t[1])
	fp.f.Inverse(t[0], t[0])
	fp.f.Mul(&c[0], &a[0], t[0])
	fp.f.Mul(t[0], &a[1], t[0])
	fp.f.Neg(&c[1], t[0])
}

func (fp *Fp2) MulByFq(c, a *Fe2, b *Fe) {
	fp.f.Mul(&c[0], &a[0], b)
	fp.f.Mul(&c[1], &a[1], b)
}

func (fp *Fp2) Exp(c, a *Fe2, e *big.Int) {
	z := fp.One()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fp.Square(z, z)
		if e.Bit(i) == 1 {
			fp.Mul(z, z, a)
		}
	}
	fp.Copy(c, z)
}

func (fp *Fp2) Div(c, a, b *Fe2) {
	t0 := fp.NewElement()
	fp.Inverse(t0, b)
	fp.Mul(c, a, t0)
}

func (fp *Fp2) FrobeniousMap(c, a *Fe2, power uint) {
	fp.f.Copy(&c[0], &a[0])
	fp.f.Mul(&c[1], &a[1], &frobeniusCoeffs2[power%2])
}

func (fp *Fp2) Sqrt(c, a *Fe2) bool {
	u, x0, a1, alpha := &Fe2{}, &Fe2{}, &Fe2{}, &Fe2{}
	fp.Copy(u, a)
	fp.Exp(a1, a, pMinus3Over4)
	fp.Square(alpha, a1)
	fp.Mul(alpha, alpha, a)
	fp.Mul(x0, a1, a)
	if fp.Equal(alpha, negativeOne2) {
		fp.f.Neg(&c[0], &x0[1])
		fp.f.Copy(&c[1], &x0[0])
		return true
	}
	fp.Add(alpha, alpha, &Fp2One)
	fp.Exp(alpha, alpha, pMinus1Over2)
	fp.Mul(c, alpha, x0)
	fp.Square(alpha, c)
	return fp.Equal(alpha, u)
}
