package bls

import (
	"fmt"
	"io"
	"math/big"
)

type Fe12 [2]Fe6

type Fp12 struct {
	f  *Fp6
	t  [4]*Fe6
	t2 [9]*Fe2
}

func NewFp12(f *Fp6) *Fp12 {
	t := [4]*Fe6{}
	t2 := [9]*Fe2{}
	for i := 0; i < 4; i++ {
		t[i] = f.Zero()
	}
	for i := 0; i < 9; i++ {
		t2[i] = &Fe2{}
	}
	if f == nil {
		return &Fp12{NewFp6(nil), t, t2}
	}
	return &Fp12{f, t, t2}
}

var Fp12One = Fe12{Fp6One, Fp6Zero}
var Fp12Zero = Fe12{Fp6Zero, Fp6Zero}

func (fp *Fp12) NewElement() *Fe12 {
	return &Fe12{}
}

func (fp *Fp12) NewElementFromBytes(f *Fe12, b []byte) error {
	if len(b) < 576 {
		return fmt.Errorf("input string should be larger than 576 bytes")
	}
	if err := fp.f.NewElementFromBytes(&f[1], b[:288]); err != nil {
		return err
	}
	if err := fp.f.NewElementFromBytes(&f[0], b[288:]); err != nil {
		return err
	}
	return nil
}

func (fp *Fp12) RandElement(a *Fe12, r io.Reader) (*Fe12, error) {
	if _, err := fp.f.RandElement(&a[0], r); err != nil {
		return nil, err
	}
	if _, err := fp.f.RandElement(&a[1], r); err != nil {
		return nil, err
	}
	return a, nil
}

func (fp *Fp12) Zero() *Fe12 {
	return &Fe12{}
}

func (fp *Fp12) One() *Fe12 {
	return &Fe12{*fp.f.One()}
}

func (fp *Fp12) ToBytes(a *Fe12) []byte {
	out := make([]byte, 576)
	copy(out[:288], fp.f.ToBytes(&a[1]))
	copy(out[288:], fp.f.ToBytes(&a[0]))
	return out
}

func (fp *Fp12) IsZero(a *Fe12) bool {
	return fp.f.IsZero(&a[0]) && fp.f.IsZero(&a[1])
}

func (fp *Fp12) Equal(a, b *Fe12) bool {
	return fp.f.Equal(&a[0], &b[0]) && fp.f.Equal(&a[1], &b[1])
}

func (fp *Fp12) Copy(c, a *Fe12) *Fe12 {
	fp.f.Copy(&c[0], &a[0])
	fp.f.Copy(&c[1], &a[1])
	return c
}

func (fp *Fp12) Add(c, a, b *Fe12) {
	fp.f.Add(&c[0], &a[0], &b[0])
	fp.f.Add(&c[1], &a[1], &b[1])

}

func (fp *Fp12) Double(c, a *Fe12) {
	fp.f.Double(&c[0], &a[0])
	fp.f.Double(&c[1], &a[1])

}

func (fp *Fp12) Sub(c, a, b *Fe12) {
	fp.f.Sub(&c[0], &a[0], &b[0])
	fp.f.Sub(&c[1], &a[1], &b[1])

}

func (fp *Fp12) Neg(c, a *Fe12) {
	fp.f.Neg(&c[0], &a[0])
	fp.f.Neg(&c[1], &a[1])
}

func (fq *Fp12) Conjugate(c, a *Fe12) {
	fq.f.Copy(&c[0], &a[0])
	fq.f.Neg(&c[1], &a[1])
}

// func (fp *Fp12) Mul(c, a, b *Fe12) {
// 	t := fp.t
// 	fp.f.Mul(t[1], &a[0], &b[0])
// 	fp.f.Mul(t[2], &a[1], &b[1])
// 	fp.f.Add(t[0], t[1], t[2])
// 	fp.f.MulByNonResidue(t[2], t[2])
// 	fp.f.Add(t[3], t[1], t[2])
// 	fp.f.Add(t[1], &a[0], &a[1])
// 	fp.f.Add(t[2], &b[0], &b[1])
// 	fp.f.Mul(t[1], t[1], t[2])
// 	fp.f.Copy(&c[0], t[3])
// 	fp.f.Sub(&c[1], t[1], t[0])
// }

var lx0, lx1, lx2, lx3 lfe6

func (fp *Fp12) mul(c, a, b *Fe12) {
	t := fp.t
	fp.f.lmul(&lx0, &a[0], &b[0])
	fp.f.lmul(&lx1, &a[1], &b[1])
	fp.f.add6(t[0], &a[0], &a[1])
	fp.f.add6(t[1], &b[0], &b[1])
	fp.f.lmul(&lx2, t[0], t[1])
	fp.f.add12(&lx3, &lx0, &lx1)
	fp.f.sub12(&lx2, &lx2, &lx3)
	fp.f.mont(&c[1], &lx2)
	fp.f.mulByNonResidue12(&lx2, &lx1)
	fp.f.add12(&lx2, &lx0, &lx2)
	fp.f.mont(&c[0], &lx2)
}

func (fp *Fp12) mulAssign(a, b *Fe12) {
	t := fp.t
	fp.f.lmul(&lx0, &a[0], &b[0])
	fp.f.lmul(&lx1, &a[1], &b[1])
	fp.f.add6(t[0], &a[0], &a[1])
	fp.f.add6(t[1], &b[0], &b[1])
	fp.f.lmul(&lx2, t[0], t[1])
	fp.f.add12(&lx3, &lx0, &lx1)
	fp.f.sub12(&lx2, &lx2, &lx3)
	fp.f.mont(&a[1], &lx2)
	fp.f.mulByNonResidue12(&lx2, &lx1)
	fp.f.add12(&lx2, &lx0, &lx2)
	fp.f.mont(&a[0], &lx2)
}

// func (fp *Fp12) MulAssign(a, b *Fe12) {
// 	t := fp.t
// 	fp.f.Mul(t[1], &a[0], &b[0])
// 	fp.f.Mul(t[2], &a[1], &b[1])
// 	fp.f.Add(t[0], t[1], t[2])
// 	fp.f.MulByNonResidue(t[2], t[2])
// 	fp.f.Add(t[3], t[1], t[2])
// 	fp.f.Add(t[1], &a[0], &a[1])
// 	fp.f.Add(t[2], &b[0], &b[1])
// 	fp.f.Mul(t[1], t[1], t[2])
// 	fp.f.Copy(&a[0], t[3])
// 	fp.f.Sub(&a[1], t[1], t[0])
// }

func (fp *Fp12) Square(c, a *Fe12) {
	t := fp.t
	fp.f.Mul(t[0], &a[0], &a[1])
	fp.f.Double(t[3], t[0])
	fp.f.MulByNonResidue(t[1], t[0])
	fp.f.Add(t[0], t[1], t[0])
	fp.f.MulByNonResidue(t[1], &a[1])
	fp.f.Add(t[1], t[1], &a[0])
	fp.f.Add(t[2], &a[0], &a[1])
	fp.f.Mul(t[2], t[1], t[2])
	fp.f.Sub(&c[0], t[2], t[0])
	fp.f.Copy(&c[1], t[3])
}

// fp.f.f.Mul(t[0], &a0, &a1)
// fp.f.f.Add(t[1], &a0, &a1)
// fp.f.f.MulByNonResidue(t[2], &a1)
// fp.f.f.Add(t[2], t[2], &a0)
// fp.f.f.MulByNonResidue(t[3], t[0])
// fp.f.f.Mul(t[4], t[1], t[2])
// fp.f.f.Sub(t[4], t[4], t[0])
// fp.f.f.Sub(t[4], t[4], t[3])
// fp.f.f.Double(t[5], t[0])

func (fp *Fp12) fp4Square(c0, c1, a0, a1 *Fe2) {
	t := fp.t2
	fp2 := fp.f.f
	// fp2.add6(t[0], a0, a1)
	// fp2.lmul(&lt0, a0, a1)
	// fp2.mulByNonResidue6(t[1], a1)
	// fp2.add6(t[1], t[1], a0)
	// fp2.mulByNonResidue12unsafe(&lt1, &lt0)
	// fp2.lmul(&lt2, t[0], t[1])
	// fp2.sub12(&lt2, &lt2, &lt0)
	// fp2.sub12(&lt2, &lt2, &lt1)
	// fp2.double12(&lt0, &lt0)
	// fp2.mont(c0, &lt2)
	// fp2.mont(c1, &lt0)

	fp2.add6(t[0], a0, a1)
	fp2.mul(t[8], a0, a1)
	fp2.mulByNonResidue6(t[1], a1)
	fp2.add6(t[1], t[1], a0)
	fp2.mulByNonResidue6(t[7], t[8])
	fp2.mul(t[6], t[0], t[1])
	fp2.sub6(t[6], t[6], t[8])
	fp2.sub6(c0, t[6], t[7])
	fp2.double6(c1, t[8])
}

func (fp *Fp12) CyclotomicSquare(c, a *Fe12) {
	t := fp.t2
	fp2 := fp.f.f
	fp.fp4Square(t[3], t[4], &a[0][0], &a[1][1])
	//
	fp2.sub6(t[2], t[3], &a[0][0])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[0][0], t[2], t[3])
	//
	fp2.add6(t[2], t[4], &a[1][1])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[1][1], t[2], t[4])

	fp.fp4Square(t[3], t[4], &a[1][0], &a[0][2])
	fp.fp4Square(t[5], t[6], &a[0][1], &a[1][2])
	//
	fp2.sub6(t[2], t[3], &a[0][1])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[0][1], t[2], t[3])
	//
	fp2.add6(t[2], t[4], &a[1][2])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[1][2], t[2], t[4])
	//
	fp2.mulByNonResidue(t[3], t[6])
	fp2.add6(t[2], t[3], &a[1][0])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[1][0], t[2], t[3])
	//
	//
	fp2.sub6(t[2], t[5], &a[0][2])
	fp2.double6(t[2], t[2])
	fp2.add6(&c[0][2], t[2], t[5])
}

// func (fp *Fp12) CyclotomicSquare(c, a *Fe12) {
// 	t := fp.t2

// 	fp.f.f.Mul(t[0], &a[0][0], &a[1][1])
// 	fp.f.f.Add(t[1], &a[0][0], &a[1][1])
// 	fp.f.f.MulByNonResidue(t[2], &a[1][1])
// 	fp.f.f.Add(t[2], t[2], &a[0][0])
// 	fp.f.f.MulByNonResidue(t[3], t[0])
// 	fp.f.f.Mul(t[4], t[1], t[2])
// 	fp.f.f.Sub(t[4], t[4], t[0])
// 	fp.f.f.Sub(t[4], t[4], t[3])
// 	fp.f.f.Double(t[5], t[0])

// 	fp.f.f.Mul(t[0], &a[1][0], &a[0][2])
// 	fp.f.f.Add(t[1], &a[1][0], &a[0][2])
// 	fp.f.f.MulByNonResidue(t[2], &a[0][2])
// 	fp.f.f.Add(t[2], t[2], &a[1][0])
// 	fp.f.f.MulByNonResidue(t[3], t[0])
// 	fp.f.f.Mul(t[6], t[1], t[2])
// 	fp.f.f.Sub(t[6], t[6], t[0])
// 	fp.f.f.Sub(t[6], t[6], t[3])
// 	fp.f.f.Double(t[7], t[0])

// 	fp.f.f.Mul(t[0], &a[0][1], &a[1][2])
// 	fp.f.f.Add(t[1], &a[0][1], &a[1][2])
// 	fp.f.f.MulByNonResidue(t[2], &a[1][2])
// 	fp.f.f.Add(t[2], t[2], &a[0][1])
// 	fp.f.f.MulByNonResidue(t[3], t[0])
// 	fp.f.f.Mul(t[8], t[1], t[2])
// 	fp.f.f.Sub(t[8], t[8], t[0])
// 	fp.f.f.Sub(t[8], t[8], t[3])
// 	fp.f.f.Double(t[0], t[0])

// 	fp.f.f.MulByNonResidue(t[0], t[0])
// 	fp.f.f.Sub(t[1], t[4], &a[0][0])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[4])
// 	fp.f.f.Copy(&c[0][0], t[1])
// 	fp.f.f.Add(t[1], t[5], &a[1][1])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[5])
// 	fp.f.f.Copy(&c[1][1], t[1])
// 	fp.f.f.Add(t[1], t[0], &a[1][0])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[0])
// 	fp.f.f.Copy(&c[1][0], t[1])
// 	fp.f.f.Sub(t[1], t[8], &a[0][2])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[8])
// 	fp.f.f.Copy(&c[0][2], t[1])
// 	fp.f.f.Sub(t[1], t[6], &a[0][1])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[6])
// 	fp.f.f.Copy(&c[0][1], t[1])
// 	fp.f.f.Add(t[1], t[7], &a[1][2])
// 	fp.f.f.Double(t[1], t[1])
// 	fp.f.f.Add(t[1], t[1], t[7])
// 	fp.f.f.Copy(&c[1][2], t[1])
// }

func (fp *Fp12) Inverse(c, a *Fe12) {
	t := fp.t
	fp.f.Square(t[0], &a[0])
	fp.f.Square(t[1], &a[1])
	fp.f.MulByNonResidue(t[1], t[1])
	fp.f.Sub(t[1], t[0], t[1])
	fp.f.Inverse(t[0], t[1])
	fp.f.Mul(&c[0], &a[0], t[0])
	fp.f.Mul(t[0], &a[1], t[0])
	fp.f.Neg(&c[1], t[0])
}

func (fp *Fp12) Div(c, a, b *Fe12) {
	t0 := fp.NewElement()
	fp.Inverse(t0, b)
	fp.mul(c, a, t0)
}

func (fq *Fp12) Exp(c, a *Fe12, e *big.Int) {
	z := fq.One()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.Square(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.Copy(c, z)
}

func (fq *Fp12) CyclotomicExp(c, a *Fe12, e *big.Int) {
	z := fq.One()
	for i := e.BitLen() - 1; i >= 0; i-- {
		fq.CyclotomicSquare(z, z)
		if e.Bit(i) == 1 {
			fq.mul(z, z, a)
		}
	}
	fq.Copy(c, z)
}

func (fp *Fp12) MulBy034Assign(a *Fe12, c0, c3, c4 *Fe2) {
	o := &Fe2{}
	t := fp.t
	fp.f.MulByBaseField(t[0], &a[0], c0)
	fp.f.Copy(t[1], &a[1])
	fp.f.mulBy01(t[1], c3, c4)
	fp.f.f.Add(o, c0, c3)
	fp.f.Add(t[2], &a[1], &a[0])
	fp.f.mulBy01(t[2], o, c4)
	fp.f.Sub(t[2], t[2], t[0])
	fp.f.Sub(&a[1], t[2], t[1])
	fp.f.MulByNonResidue(t[1], t[1])
	fp.f.Add(&a[0], t[0], t[1])
}

func (fp *Fp12) MulBy014Assign(a *Fe12, c0, c1, c4 *Fe2) {
	o := &Fe2{}
	t := fp.t
	fp.f.Copy(t[0], &a[0])
	fp.f.mulBy01(t[0], c0, c1)
	fp.f.Copy(t[1], &a[1])
	fp.f.mulBy1(t[1], c4)
	fp.f.f.Add(o, c1, c4)
	fp.f.Add(&a[1], &a[1], &a[0])
	fp.f.mulBy01(&a[1], c0, o)
	fp.f.Sub(&a[1], &a[1], t[0])
	fp.f.Sub(&a[1], &a[1], t[1])
	fp.f.MulByNonResidue(t[1], t[1])
	fp.f.Add(&a[0], t[1], t[0])
}

func (fp *Fp12) FrobeniusMap(c, a *Fe12, power uint) {
	fp.f.FrobeniusMap(&c[0], &a[0], power)
	fp.f.FrobeniusMap(&c[1], &a[1], power)
	fp.f.MulByBaseField(&c[1], &c[1], &frobeniusCoeffs12[power%12])
}

func (fp *Fp12) FrobeniusMapAssign(a *Fe12, power uint) {
	fp.f.FrobeniusMap(&a[0], &a[0], power)
	fp.f.FrobeniusMap(&a[1], &a[1], power)
	fp.f.MulByBaseField(&a[1], &a[1], &frobeniusCoeffs12[power%12])
}
