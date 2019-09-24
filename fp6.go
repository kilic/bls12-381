package bls

import (
	"fmt"
	"io"
	"math/big"
)

type Fe6 [3]Fe2

func NewFp6(f *Fp2) *Fp6 {
	t := [6]*Fe2{}
	for i := 0; i < 6; i++ {
		t[i] = f.Zero()
	}
	if f == nil {
		return &Fp6{NewFp2(nil), t}
	}
	return &Fp6{f, t}
}

type Fp6 struct {
	f *Fp2
	t [6]*Fe2
}

var fp6One = Fe6{fp2One, fp2Zero, fp2Zero}
var fp6Zero = Fe6{fp2Zero, fp2Zero, fp2Zero}

func (fp *Fp6) NewElement() *Fe6 {
	return &Fe6{}
}

func (fp *Fp6) NewElementFromBytes(c *Fe6, b []byte) error {
	if len(b) < 288 {
		return fmt.Errorf("input string should be larger than 288 bytes")
	}
	if err := fp.f.NewElementFromBytes(&c[2], b[:96]); err != nil {
		return err
	}
	if err := fp.f.NewElementFromBytes(&c[1], b[96:192]); err != nil {
		return err
	}
	if err := fp.f.NewElementFromBytes(&c[0], b[192:]); err != nil {
		return err
	}
	return nil
}

func (fp *Fp6) RandElement(a *Fe6, r io.Reader) (*Fe6, error) {
	if _, err := fp.f.RandElement(&a[0], r); err != nil {
		return nil, err
	}
	if _, err := fp.f.RandElement(&a[1], r); err != nil {
		return nil, err
	}
	if _, err := fp.f.RandElement(&a[2], r); err != nil {
		return nil, err
	}
	return a, nil
}

func (fp *Fp6) Zero() *Fe6 {
	return &Fe6{}
}

func (fp *Fp6) One() *Fe6 {
	return &Fe6{*fp.f.One()}
}

func (fp *Fp6) ToBytes(a *Fe6) []byte {
	out := make([]byte, 288)
	copy(out[:96], fp.f.ToBytes(&a[2]))
	copy(out[96:192], fp.f.ToBytes(&a[1]))
	copy(out[192:], fp.f.ToBytes(&a[0]))
	return out
}

func (fp *Fp6) IsZero(a *Fe6) bool {
	return fp.f.IsZero(&a[0]) && fp.f.IsZero(&a[1]) && fp.f.IsZero(&a[2])
}

func (fp *Fp6) Equal(a, b *Fe6) bool {
	return fp.f.Equal(&a[0], &b[0]) && fp.f.Equal(&a[1], &b[1]) && fp.f.Equal(&a[2], &b[2])
}

func (fp *Fp6) Copy(c, a *Fe6) *Fe6 {
	fp.f.Copy(&c[0], &a[0])
	fp.f.Copy(&c[1], &a[1])
	fp.f.Copy(&c[2], &a[2])
	return c
}

func (fp *Fp6) MulByNonResidue(c, a *Fe6) {
	t := fp.t
	fp.f.Copy(t[0], &a[0])
	fp.f.MulByNonResidue(&c[0], &a[2])
	fp.f.Copy(&c[2], &a[1])
	fp.f.Copy(&c[1], t[0])
}

func (fp *Fp6) Add(c, a, b *Fe6) {
	fp.f.Add(&c[0], &a[0], &b[0])
	fp.f.Add(&c[1], &a[1], &b[1])
	fp.f.Add(&c[2], &a[2], &b[2])
}

func (fp *Fp6) Double(c, a *Fe6) {
	fp.f.Double(&c[0], &a[0])
	fp.f.Double(&c[1], &a[1])
	fp.f.Double(&c[2], &a[2])
}

func (fp *Fp6) Sub(c, a, b *Fe6) {
	fp.f.Sub(&c[0], &a[0], &b[0])
	fp.f.Sub(&c[1], &a[1], &b[1])
	fp.f.Sub(&c[2], &a[2], &b[2])
}

func (fp *Fp6) Neg(c, a *Fe6) {
	fp.f.Neg(&c[0], &a[0])
	fp.f.Neg(&c[1], &a[1])
	fp.f.Neg(&c[2], &a[2])
}

func (fq *Fp6) Conjugate(c, a *Fe6) {
	fq.f.Copy(&c[0], &a[0])
	fq.f.Neg(&c[1], &a[1])
	fq.f.Copy(&c[2], &a[2])
}

func (fp *Fp6) Mul(c, a, b *Fe6) {
	t := fp.t
	fp.f.Mul(t[0], &a[0], &b[0])
	fp.f.Mul(t[1], &a[1], &b[1])
	fp.f.Mul(t[2], &a[2], &b[2])
	fp.f.Add(t[3], &a[1], &a[2])
	fp.f.Add(t[4], &b[1], &b[2])
	fp.f.Mul(t[3], t[3], t[4])
	fp.f.Add(t[4], t[1], t[2])
	fp.f.Sub(t[3], t[3], t[4])
	fp.f.MulByNonResidue(t[3], t[3])
	fp.f.Add(t[5], t[0], t[3])
	fp.f.Add(t[3], &a[0], &a[1])
	fp.f.Add(t[4], &b[0], &b[1])
	fp.f.Mul(t[3], t[3], t[4])
	fp.f.Add(t[4], t[0], t[1])
	fp.f.Sub(t[3], t[3], t[4])
	fp.f.MulByNonResidue(t[4], t[2])
	fp.f.Add(&c[1], t[3], t[4])
	fp.f.Add(t[3], &a[0], &a[2])
	fp.f.Add(t[4], &b[0], &b[2])
	fp.f.Mul(t[3], t[3], t[4])
	fp.f.Add(t[4], t[0], t[2])
	fp.f.Sub(t[3], t[3], t[4])
	fp.f.Add(&c[2], t[1], t[3])
	fp.f.Copy(&c[0], t[5])
}

func (fp *Fp6) Square(c, a *Fe6) {
	t := fp.t
	fp.f.Square(t[0], &a[0])
	fp.f.Mul(t[1], &a[0], &a[1])
	fp.f.Add(t[1], t[1], t[1])
	fp.f.Sub(t[2], &a[0], &a[1])
	fp.f.Add(t[2], t[2], &a[2])
	fp.f.Square(t[2], t[2])
	fp.f.Mul(t[3], &a[1], &a[2])
	fp.f.Add(t[3], t[3], t[3])
	fp.f.Square(t[4], &a[2])
	fp.f.MulByNonResidue(t[5], t[3])
	fp.f.Add(&c[0], t[0], t[5])
	fp.f.MulByNonResidue(t[5], t[4])
	fp.f.Add(&c[1], t[1], t[5])
	fp.f.Add(t[1], t[1], t[2])
	fp.f.Add(t[1], t[1], t[3])
	fp.f.Add(t[0], t[0], t[4])
	fp.f.Sub(&c[2], t[1], t[0])
}

func (fp *Fp6) Inverse(c, a *Fe6) {
	t := fp.t
	fp.f.Square(t[0], &a[0])
	fp.f.Mul(t[1], &a[1], &a[2])
	fp.f.MulByNonResidue(t[1], t[1])
	fp.f.Sub(t[0], t[0], t[1])
	fp.f.Square(t[1], &a[1])
	fp.f.Mul(t[2], &a[0], &a[2])
	fp.f.Sub(t[1], t[1], t[2])
	fp.f.Square(t[2], &a[2])
	fp.f.MulByNonResidue(t[2], t[2])
	fp.f.Mul(t[3], &a[0], &a[1])
	fp.f.Sub(t[2], t[2], t[3])
	fp.f.Mul(t[3], &a[2], t[2])
	fp.f.Mul(t[4], &a[1], t[1])
	fp.f.Add(t[3], t[3], t[4])
	fp.f.MulByNonResidue(t[3], t[3])
	fp.f.Mul(t[4], &a[0], t[0])
	fp.f.Add(t[3], t[3], t[4])
	fp.f.Inverse(t[3], t[3])
	fp.f.Mul(&c[0], t[0], t[3])
	fp.f.Mul(&c[1], t[2], t[3])
	fp.f.Mul(&c[2], t[1], t[3])
}

func (fp *Fp6) Div(c, a, b *Fe6) {
	t0 := fp.NewElement()
	fp.Inverse(t0, b)
	fp.Mul(c, a, t0)
}

func (fq *Fp6) Exp(c, a *Fe6, e *big.Int) {
	z := fq.One()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.Square(z, z)
		if e.Bit(i) == 1 {
			fq.Mul(z, z, a)
		}
	}
	fq.Copy(c, z)
}

func (fp *Fp6) MulByBaseField(c, a *Fe6, b *Fe2) {
	fp.f.Mul(&c[0], &a[0], b)
	fp.f.Mul(&c[1], &a[1], b)
	fp.f.Mul(&c[2], &a[2], b)
}

func (fp *Fp6) MulBy01(a *Fe6, c0, c1 *Fe2) {
	t := fp.t
	fp.f.Mul(t[0], &a[0], c0)
	fp.f.Mul(t[1], &a[1], c1)
	fp.f.Add(t[5], &a[1], &a[2])
	fp.f.Mul(t[2], c1, t[5])
	fp.f.Sub(t[2], t[2], t[1])
	fp.f.MulByNonResidue(t[2], t[2])
	fp.f.Add(t[2], t[2], t[0])
	fp.f.Add(t[5], &a[0], &a[2])
	fp.f.Mul(t[3], c0, t[5])
	fp.f.Sub(t[3], t[3], t[0])
	fp.f.Add(t[3], t[3], t[1])
	fp.f.Add(t[4], c0, c1)
	fp.f.Add(t[5], &a[0], &a[1])
	fp.f.Mul(t[4], t[4], t[5])
	fp.f.Sub(t[4], t[4], t[0])
	fp.f.Sub(t[4], t[4], t[1])
	fp.f.Copy(&a[0], t[2])
	fp.f.Copy(&a[1], t[4])
	fp.f.Copy(&a[2], t[3])
}

func (fp *Fp6) MulBy1(a *Fe6, c1 *Fe2) {
	t := fp.t
	fp.f.Mul(t[0], &a[1], c1)
	fp.f.Add(t[1], &a[1], &a[2])
	fp.f.Mul(t[1], t[1], c1)
	fp.f.Sub(t[1], t[1], t[0])
	fp.f.MulByNonResidue(t[1], t[1])
	fp.f.Add(t[2], &a[0], &a[1])
	fp.f.Mul(t[2], t[2], c1)
	fp.f.Sub(&a[1], t[2], t[0])
	fp.f.Copy(&a[0], t[1])
	fp.f.Copy(&a[2], t[0])
}

func (fp *Fp6) FrobeniusMap(c, a *Fe6, power uint) {
	fp.f.FrobeniousMap(&c[0], &a[0], power)
	fp.f.FrobeniousMap(&c[1], &a[1], power)
	fp.f.FrobeniousMap(&c[2], &a[2], power)
	fp.f.Mul(&c[1], &c[1], &frobeniusCoeffs61[power%6])
	fp.f.Mul(&c[2], &c[2], &frobeniusCoeffs62[power%6])
}
