package bls

import (
	"fmt"
	"io"
	"math/big"
)

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

func (fp *Fp2) copy(c, a *Fe2) {
	fp.f.copy(&c[0], &a[0])
	fp.f.copy(&c[1], &a[1])
}

func (fp *Fp2) lcopy(c, a *lfe2) {
	fp.f.lcopy(&c[0], &a[0])
	fp.f.lcopy(&c[1], &a[1])
}

func (fp *Fp2) copyMixed(c *lfe2, a *Fe2) {
	fp.f.copyMixed(&c[0], &a[0])
	fp.f.copyMixed(&c[1], &a[1])
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
	fp.f.ladd6(t[0], &a[0], &a[1])
	fp.f.sub6(t[1], &a[0], &a[1])
	fp.f.ldouble6(t[2], &a[0])
	fp.f.Mul(&c[0], t[0], t[1])
	fp.f.Mul(&c[1], t[2], &a[1])
}

func (fp *Fp2) MulByNonResidue(c, a *Fe2) {
	t := fp.t
	fp.f.Sub(t[0], &a[0], &a[1])
	fp.f.Add(&c[1], &a[0], &a[1])
	fp.f.Copy(&c[0], t[0])
}

func (fp *Fp2) mulByNonResidue(c, a *Fe2) {
	t := fp.t
	fp.f.Sub(t[0], &a[0], &a[1])
	fp.f.Add(&c[1], &a[0], &a[1])
	fp.f.Copy(&c[0], t[0])
}

var l0, l1, l2, l3, l4, l5 lfe

func (fp *Fp2) mulByNonResidue6(c, a *Fe2) {
	t := fp.t
	fp.f.sub6(t[0], &a[0], &a[1])
	fp.f.add6(&c[1], &a[0], &a[1])
	fp.f.copy(&c[0], t[0])
}

func (fp *Fp2) mulByNonResidue12(a *lfe2) {
	fp.f.sub12(&l0, &a[0], &a[1])
	fp.f.add12(&a[1], &a[0], &a[1])
	fp.f.lcopy(&a[0], &l0)
}

func (fp *Fp2) mulByNonResidue12unsafe(c, a *lfe2) {
	fp.f.sub12(&c[0], &a[0], &a[1])
	fp.f.add12(&c[1], &a[0], &a[1])
}

func (fp *Fp2) add6(c, a, b *Fe2) {
	fp.f.add6(&c[0], &a[0], &b[0])
	fp.f.add6(&c[1], &a[1], &b[1])
}

func (fp *Fp2) ladd6(c, a, b *Fe2) {
	fp.f.ladd6(&c[0], &a[0], &b[0])
	fp.f.ladd6(&c[1], &a[1], &b[1])
}

func (fp *Fp2) add12(c, a, b *lfe2) {
	fp.f.add12(&c[0], &a[0], &b[0])
	fp.f.add12(&c[1], &a[1], &b[1])
}

func (fp *Fp2) ladd12(c, a, b *lfe2) {
	fp.f.ladd12(&c[0], &a[0], &b[0])
	fp.f.ladd12(&c[1], &a[1], &b[1])
}

func (fp *Fp2) double6(c, a *Fe2) {
	fp.f.double6(&c[0], &a[0])
	fp.f.double6(&c[1], &a[1])
}

func (fp *Fp2) ldouble6(c, a *Fe2) {
	fp.f.ldouble6(&c[0], &a[0])
	fp.f.ldouble6(&c[1], &a[1])
}

func (fp *Fp2) double12(c, a *lfe2) {
	fp.f.double12(&c[0], &a[0])
	fp.f.double12(&c[1], &a[1])
}

func (fp *Fp2) ldouble12(c, a *lfe2) {
	fp.f.ldouble12(&c[0], &a[0])
	fp.f.ldouble12(&c[1], &a[1])
}

func (fp *Fp2) sub6(c, a, b *Fe2) {
	fp.f.sub6(&c[0], &a[0], &b[0])
	fp.f.sub6(&c[1], &a[1], &b[1])
}

func (fp *Fp2) lsub6(c, a, b *Fe2) {
	fp.f.lsub6(&c[0], &a[0], &b[0])
	fp.f.lsub6(&c[1], &a[1], &b[1])
}

func (fp *Fp2) sub12(c, a, b *lfe2) {
	fp.f.sub12(&c[0], &a[0], &b[0])
	fp.f.sub12(&c[1], &a[1], &b[1])
}

func (fp *Fp2) lsub12(c, a, b *lfe2) {
	fp.f.lsub12(&c[0], &a[0], &b[0])
	fp.f.lsub12(&c[1], &a[1], &b[1])
}

func (fp *Fp2) submixed12(c, a, b *lfe2) {
	fp.f.sub12(&c[0], &a[0], &b[0])
	fp.f.lsub12(&c[1], &a[1], &b[1])
}

func (fp *Fp2) mul(c, a, b *Fe2) {
	t := fp.t
	fp.f.lmul(&l1, &a[0], &b[0])
	fp.f.lmul(&l2, &a[1], &b[1])
	fp.f.ladd6(t[0], &a[0], &a[1])
	fp.f.ladd6(t[1], &b[0], &b[1])
	fp.f.lmul(&l3, t[0], t[1])
	fp.f.lsub12(&l3, &l3, &l1)
	fp.f.lsub12(&l3, &l3, &l2)
	fp.f.mont(&c[1], &l3)
	fp.f.sub12(&l1, &l1, &l2)
	fp.f.mont(&c[0], &l1)
}

func (fp *Fp2) lmul(c *lfe2, a, b *Fe2) {
	t := fp.t
	fp.f.lmul(&l1, &a[0], &b[0])
	fp.f.lmul(&l2, &a[1], &b[1])
	fp.f.ladd6(t[0], &a[0], &a[1])
	fp.f.ladd6(t[1], &b[0], &b[1])
	fp.f.lmul(&l3, t[0], t[1])
	fp.f.lsub12(&l3, &l3, &l1)
	fp.f.lsub12(&c[1], &l3, &l2)
	fp.f.sub12(&c[0], &l1, &l2)
	// fp.f.lsub12opt1h2(&c[0], &l1, &l2)
}

func (fp *Fp2) lsquare(c *lfe2, a *Fe2) {
	t := fp.t
	fp.f.ladd6(t[0], &a[0], &a[1])
	fp.f.lsub6(t[1], &a[0], &a[1])
	fp.f.lmul(&c[0], t[0], t[1])
	fp.f.ldouble6(t[2], &a[0])
	fp.f.lmul(&c[1], t[2], &a[1])
}

func (fp *Fp2) square(c, a *Fe2) {
	t := fp.t
	fp.f.ladd6(t[0], &a[0], &a[1])
	fp.f.lsub6(t[1], &a[0], &a[1])
	fp.f.mul(&c[0], t[0], t[1])
	fp.f.ldouble6(t[2], &a[0])
	fp.f.mul(&c[1], t[2], &a[1])
}

func (fp *Fp2) mont(c *Fe2, a *lfe2) {
	fp.f.mont(&c[0], &a[0])
	fp.f.mont(&c[1], &a[1])
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
	fp.f.mul(&c[0], &a[0], b)
	fp.f.mul(&c[1], &a[1], b)
}

func (fp *Fp2) mulByFq(c, a *Fe2, b *Fe) {
	fp.f.mul(&c[0], &a[0], b)
	fp.f.mul(&c[1], &a[1], b)
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
	fp.f.copy(&c[0], &a[0])
	fp.f.mul(&c[1], &a[1], &frobeniusCoeffs2[power%2])
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
