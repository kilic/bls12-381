package bls

type BLSPairingEngine struct {
	G1   *G1
	G2   *G2
	Fp12 *Fp12
	Fp2  *Fp2
	Fp   *Fp
	t2   [10]*Fe2
	t12  [9]Fe12
}

func NewBLSPairingEngine() *BLSPairingEngine {
	fp := NewFp()
	fp2 := NewFp2(fp)
	fp6 := NewFp6(fp2)
	fp12 := NewFp12(fp6)
	g1 := NewG1(fp)
	g2 := NewG2(fp2)
	t2 := [10]*Fe2{}
	for i := 0; i < 10; i++ {
		t2[i] = fp2.Zero()
	}
	t12 := [9]Fe12{}
	return &BLSPairingEngine{
		Fp:   fp,
		Fp2:  fp2,
		Fp12: fp12,
		t2:   t2,
		t12:  t12,
		G1:   g1,
		G2:   g2,
	}
}

func (e *BLSPairingEngine) doublingStep(coeff *[3]Fe2, r *PointG2) {
	fp2 := e.Fp2
	t := e.t2
	fp2.Mul(t[0], &r[0], &r[1])
	fp2.MulByFq(t[0], t[0], twoInv)
	fp2.Square(t[1], &r[1])
	fp2.Square(t[2], &r[2])
	fp2.Copy(t[3], b2)
	fp2.Double(t[7], t[2])
	fp2.Add(t[7], t[7], t[2])
	fp2.Mul(t[3], t[3], t[7])
	fp2.Double(t[4], t[3])
	fp2.Add(t[4], t[4], t[3])
	fp2.Add(t[5], t[1], t[4])
	fp2.MulByFq(t[5], t[5], twoInv)
	fp2.Add(t[6], &r[1], &r[2])
	fp2.Square(t[6], t[6])
	fp2.Add(t[7], t[2], t[1])
	fp2.Sub(t[6], t[6], t[7])
	fp2.Sub(&coeff[0], t[3], t[1])
	fp2.Square(t[7], &r[0])
	fp2.Sub(t[4], t[1], t[4])
	fp2.Mul(&r[0], t[4], t[0])
	fp2.Square(t[2], t[3])
	fp2.Double(t[3], t[2])
	fp2.Add(t[3], t[3], t[2])
	fp2.Square(t[5], t[5])
	fp2.Sub(&r[1], t[5], t[3])
	fp2.Mul(&r[2], t[1], t[6])
	fp2.Double(t[0], t[7])
	fp2.Add(&coeff[1], t[0], t[7])
	fp2.Neg(&coeff[2], t[6])

}

func (e *BLSPairingEngine) additionStep(coeff *[3]Fe2, r, q *PointG2) {
	fp2 := e.Fp2
	t := e.t2
	fp2.Mul(t[0], &q[1], &r[2])
	fp2.Neg(t[0], t[0])
	fp2.Add(t[0], t[0], &r[1])
	fp2.Mul(t[1], &q[0], &r[2])
	fp2.Neg(t[1], t[1])
	fp2.Add(t[1], t[1], &r[0])
	fp2.Square(t[2], t[0])
	fp2.Square(t[3], t[1])
	fp2.Mul(t[4], t[1], t[3])
	fp2.Mul(t[2], &r[2], t[2])
	fp2.Mul(t[3], &r[0], t[3])
	fp2.Double(t[5], t[3])
	fp2.Sub(t[5], t[4], t[5])
	fp2.Add(t[5], t[5], t[2])
	fp2.Mul(&r[0], t[1], t[5])
	fp2.Sub(t[2], t[3], t[5])
	fp2.Mul(t[2], t[2], t[0])
	fp2.Mul(t[3], &r[1], t[4])
	fp2.Sub(&r[1], t[2], t[3])
	fp2.Mul(&r[2], &r[2], t[4])
	fp2.Mul(t[2], t[1], &q[1])
	fp2.Mul(t[3], t[0], &q[0])
	fp2.Sub(&coeff[0], t[3], t[2])
	fp2.Neg(&coeff[1], t[0])
	fp2.Copy(&coeff[2], t[1])
}

func (e *BLSPairingEngine) prepare(ellCoeffs *[70][3]Fe2, twistPoint *PointG2) {
	if e.G2.IsZero(twistPoint) {
		return
	}
	r := &PointG2{}
	e.G2.Copy(r, twistPoint)
	j := 0
	for i := int(x.BitLen() - 2); i >= 0; i-- {
		e.doublingStep(&ellCoeffs[j], r)
		if x.Bit(i) != 0 {
			j++
			ellCoeffs[j] = Fe6{}
			e.additionStep(&ellCoeffs[j], r, twistPoint)
		}
		j++
	}
}

// notice that this function expects: len(points) == len(twistPoints)
func (e *BLSPairingEngine) millerLoop(f *Fe12, points []PointG1, twistPoints []PointG2) {
	for i := 0; i <= len(points)-1; i++ {
		e.G1.Affine(&points[i])
		e.G2.Affine(&twistPoints[i])
	}
	ellCoeffs := make([][70][3]Fe2, len(points))
	for i := 0; i < len(points); i++ {
		if !e.G1.IsZero(&points[i]) && !e.G2.IsZero(&twistPoints[i]) {
			e.prepare(&ellCoeffs[i], &twistPoints[i])
		}
	}
	fp12 := e.Fp12
	fp2 := e.Fp2
	t := e.t2
	j := 0
	// ell := func(f *Fe12, coeffs *[3]Fe2, p *PointG1) {
	// 	t := [3]Fe2{}
	// 	fp2.MulByFq(&t[0], &coeffs[2], &p[1])
	// 	fp2.MulByFq(&t[1], &coeffs[1], &p[0])
	// 	fp12.MulBy014(f, &coeffs[0], &t[1], &t[0])
	// }
	//
	// solveLine := func(f *Fe12) {
	// 	for i := 0; i <= len(points)-1; i++ {
	// 		ell(f, &ellCoeffs[i][j], &points[i])
	// 	}
	// }
	fp12.Copy(f, &Fp12One)
	for i := int(x.BitLen() - 2); i >= 0; i-- {
		fp12.Square(f, f)
		//solveLine(f)
		for i := 0; i <= len(points)-1; i++ {
			// ell(f, &ellCoeffs[i][j], &points[i])
			fp2.MulByFq(t[0], &ellCoeffs[i][j][2], &points[i][1])
			fp2.MulByFq(t[1], &ellCoeffs[i][j][1], &points[i][0])
			fp12.MulBy014Assign(f, &ellCoeffs[i][j][0], t[1], t[0])
		}
		if x.Bit(i) != 0 {
			j++
			// solveLine(f)
			for i := 0; i <= len(points)-1; i++ {
				// ell(f, &ellCoeffRefs[i][j], &points[i])
				fp2.MulByFq(t[0], &ellCoeffs[i][j][2], &points[i][1])
				fp2.MulByFq(t[1], &ellCoeffs[i][j][1], &points[i][0])
				fp12.MulBy014Assign(f, &ellCoeffs[i][j][0], t[1], t[0])
			}
		}
		j++
	}
	fp12.Conjugate(f, f)
}

func (e *BLSPairingEngine) exp(c, a *Fe12) {
	fp12 := e.Fp12
	fp12.CyclotomicExp(c, a, x)
	fp12.Conjugate(c, c)
}

// assigned operation
func (e *BLSPairingEngine) finalExp(f *Fe12) {
	fp12 := e.Fp12
	t := e.t12
	fp12.FrobeniusMap(&t[0], f, 6)
	fp12.Inverse(&t[1], f)
	fp12.Mul(&t[2], &t[0], &t[1])
	fp12.Copy(&t[1], &t[2])
	fp12.FrobeniusMapAssign(&t[2], 2)
	fp12.MulAssign(&t[2], &t[1])
	fp12.CyclotomicSquare(&t[1], &t[2])
	fp12.Conjugate(&t[1], &t[1])
	e.exp(&t[3], &t[2])
	fp12.CyclotomicSquare(&t[4], &t[3])
	fp12.Mul(&t[5], &t[1], &t[3])
	e.exp(&t[1], &t[5])
	e.exp(&t[0], &t[1])
	e.exp(&t[6], &t[0])
	fp12.MulAssign(&t[6], &t[4])
	e.exp(&t[4], &t[6])
	fp12.Conjugate(&t[5], &t[5])
	fp12.MulAssign(&t[4], &t[5])
	fp12.MulAssign(&t[4], &t[2])
	fp12.Conjugate(&t[5], &t[2])
	fp12.MulAssign(&t[1], &t[2])
	fp12.FrobeniusMapAssign(&t[1], 3)
	fp12.MulAssign(&t[6], &t[5])
	fp12.FrobeniusMapAssign(&t[6], 1)
	fp12.MulAssign(&t[3], &t[0])
	fp12.FrobeniusMapAssign(&t[3], 2)
	fp12.MulAssign(&t[3], &t[1])
	fp12.MulAssign(&t[3], &t[6])
	fp12.Mul(f, &t[3], &t[4])
}

func (e *BLSPairingEngine) Pair(f *Fe12, points []PointG1, twistPoints []PointG2) {
	e.millerLoop(f, points, twistPoints)
	e.finalExp(f)
}

func (e *BLSPairingEngine) Equal(f *Fe12, points []PointG1, twistPoints []PointG2) {
	e.millerLoop(f, points, twistPoints)
	e.finalExp(f)
}
