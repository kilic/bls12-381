package bls

import (
	"fmt"
	"io"
	"math/big"
)

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

var Fp6One = Fe6{Fp2One, Fp2Zero, Fp2Zero}
var Fp6Zero = Fe6{Fp2Zero, Fp2Zero, Fp2Zero}

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

var lt0, lt1, lt2, lt3, lt4, lt5, lt6, l20, l21, l22, l23, l24 lfe2
var v0, v1, v2 lfe2

func (fp *Fp6) mont(c *Fe6, a *lfe6) {
	fp.f.mont(&c[0], &a[0])
	fp.f.mont(&c[1], &a[1])
	fp.f.mont(&c[2], &a[2])
}

func (fp *Fp6) mul(c *lfe6, a, b *Fe6) {

	t := fp.t
	fp1 := fp.f.f
	fp2 := fp.f
	//	vi = ai * bi
	//	-- vi0 ∈ [0, 2^(N-2) * p + p^2]
	//	vi0 ∈ [0, 2^N * p]
	//	vi1 ∈ [0, 2p^2]
	fp2.mul(&v0, &a[0], &b[0])
	fp2.mul(&v1, &a[1], &b[1])
	fp2.mul(&v2, &a[2], &b[2])

	// c0
	//
	// t0 = (a1 + a2)
	// t1 = (b1 + b2)
	// t0, t1 ∈ [0,2p]
	fp2.ladd6(t[0], &a[1], &a[2])
	fp2.ladd6(t[1], &b[1], &b[2])
	//	lt0	=	(t0 * t1)
	//			= (a1 + a2) * (b1 + b2)
	//	lt00	∈	[0, 2^N * p]
	//	lt01	∈	[0, 8p^2]
	fp2.mul(&lt0, t[0], t[1])
	//	lt0	= (t0 * t1) - v1 - v2
	//			= (a1 + a2) * (b1 + b2) - v1 - v2
	//			= (a1 * b1) + (a2 * b2)
	//	lt00 ∈ [0, 2^N * p]
	//	lt01 ∈ [0, 4 p^2]
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v1[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v1[1])
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v2[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v2[1])

	//  lt1 = lt0 * α
	//	lt10 = lt00 - lt01
	//	lt11 = lt00 + lt01
	//	lt00, lt01 ∈ [0, 2^N * p]
	fp1.lsub12opt2(&lt1[0], &lt0[0], &lt0[1])
	fp1.ladd12opt2(&lt1[1], &lt0[0], &lt0[1])

	//	c0 = lt1 + v0
	fp2.ladd12opt2(&c[0], &lt1, &v0)

	// c1
	//
	// t0 = (a0 + a1)
	// t1 = (b0 + b1)
	// t0, t1 ∈ [0,2p]
	fp2.ladd6(t[0], &a[0], &a[1])
	fp2.ladd6(t[1], &b[0], &b[1])
	//	lt0	=	(t0 * t1)
	//			= (a0 + a1) * (b0 + b1)
	//	lt00	∈	[0, 2^N * p]
	//	lt01	∈	[0, 8p^2]
	fp2.mul(&lt0, t[0], t[1])
	//	lt0	= (t0 * t1) - v0 - v1
	//			= (a0 + a1) * (b0 + b1) - v0 - v1
	//			= (a0 * b1) + (a0 * b1)
	//	lt00 ∈ [0, 2^N * p]
	//	lt01 ∈ [0, 4 p^2]
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v0[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v0[1])
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v1[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v1[1])

	//  lt1 = v2 * α
	//	lt10 = v20 + v21 // opt1 in first phase might help here
	//	lt11 = v20 - v21
	//	lt00, lt01 ∈ [0, 2^N * p]
	fp1.lsub12opt2(&lt1[0], &v2[0], &v2[1])
	fp1.ladd12opt2(&lt1[1], &v2[0], &v2[1])

	//	c1 = lt0 + lt1
	fp2.ladd12opt2(&c[1], &lt1, &lt0)

	// c2
	//
	// t0 = (a0 + a2)
	// t1 = (b0 + b2)
	// t0, t1 ∈ [0,2p]
	fp2.ladd6(t[0], &a[0], &a[2])
	fp2.ladd6(t[1], &b[0], &b[2])
	//	lt0	=	(t0 * t1)
	//			= (a0 + a2) * (b0 + b2)
	//	lt00	∈	[0, 2^N * p]
	//	lt01	∈	[0, 8p^2]
	fp2.mul(&lt0, t[0], t[1])
	//	lt0	= (t0 * t1) - v0 - v2
	//			= (a0 + a2) * (b0 + b2) - v0 - v2
	//			= (a0 * b2) + (a0 * b2)
	//	lt00 ∈ [0, 2^N * p]
	//	lt01 ∈ [0, 4 p^2]
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v0[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v0[1])
	fp1.lsub12opt2(&lt0[0], &lt0[0], &v2[0])
	fp1.lsub12(&lt0[1], &lt0[1], &v2[1])

	//	c2 = lt0 + v1
	fp2.ladd12opt2(&c[2], &lt0, &v1)
	fp1.ladd12opt2(&c[2][0], &lt0[0], &v1[0])
	fp1.ladd12(&c[2][1], &lt0[1], &v1[1])
}

var l60 lfe6

func (fp *Fp6) mulr(c, a, b *Fe6) {
	fp.mul(&l60, a, b)
	fp.mont(c, &l60)
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

func (fp *Fp6) square(lc *lfe6, a *Fe6) {
	t := fp.t
	fp1 := fp.f.f
	fp2 := fp.f
	//	lt0 = a0 * b1
	//	lt00 ∈ [0, 2^N * p]
	//	lt01 ∈ [0, 2p^2]
	fp2.mul(&lt0, &a[0], &a[1])
	//	lt00 ∈ [0, 2^N * p]
	//	lt01 ∈ [0, 4p^2]
	fp1.ldouble12opt2(&lt0[0], &lt0[0])
	fp1.ldouble12(&lt0[1], &lt0[1])
	// v4 = 2(a0 * a1)

	// v5
	fp2.square(&lt1, &a[2])
	// α * v5
	fp1.lsub12opt2(&lt2[0], &lt1[0], &lt1[1])
	fp1.ladd12opt2(&lt2[1], &lt1[0], &lt1[1])

	//  c = 2(a0 * a1) + v4 @ lt2
	fp2.ladd12opt2(&lt2, &lt2, &lt0)

	fp2.lsub12opt2(&lt0, &lt0, &lt1) // v2
	fp2.square(&lt1, &a[0])          // v3
	fp2.sub6(t[0], &a[0], &a[1])
	fp2.add6(t[0], t[0], &a[2])
	fp2.square(&lt3, t[0]) // v4

	fp2.mul(&lt4, &a[1], &a[2])
	fp1.ldouble12opt2(&lt4[0], &lt4[0])
	fp1.ldouble12(&lt4[1], &lt4[1]) // v5 //*

	// α * v5
	fp1.lsub12opt2(&lt5[0], &lt4[0], &lt4[1])
	fp1.ladd12opt2(&lt5[1], &lt4[0], &lt4[1])

	fp2.ladd12opt2(&lc[0], &lt5, &lt1)
	fp2.lcopy(&lc[1], &lt2)

	fp2.lsub12opt2(&lt3, &lt3, &lt1)   // v4 - v3
	fp2.ladd12opt2(&lt3, &lt3, &lt0)   // v4 - v3 + v2
	fp2.ladd12opt2(&lc[2], &lt3, &lt4) // v4 - v3 + v2 + v5
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

func (fp *Fp6) lMulBy01(a *Fe6, c0, c1 *Fe2) {
	t := fp.t
	fp.f.mul(&v0, &a[0], c0)
	fp.f.mul(&v1, &a[1], c1)

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
