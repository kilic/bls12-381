package bls

import (
	"fmt"
	"io"
	"math/big"
)

type fp2Temp struct {
	t [4]*fe
}

type fp2 struct {
	fp2Temp
}

func newFp2Temp() fp2Temp {
	t := [4]*fe{}
	for i := 0; i < len(t); i++ {
		t[i] = &fe{}
	}
	return fp2Temp{t}
}

func newFp2() *fp2 {
	t := newFp2Temp()
	return &fp2{t}
}

func (e *fp2) fromBytes(in []byte) (*fe2, error) {
	if len(in) != 96 {
		return nil, fmt.Errorf("input string should be larger than 96 bytes")
	}
	c1, err := fromBytes(in[:48])
	if err != nil {
		return nil, err
	}
	c0, err := fromBytes(in[48:])
	if err != nil {
		return nil, err
	}
	return &fe2{*c0, *c1}, nil
}

func (e *fp2) toBytes(a *fe2) []byte {
	out := make([]byte, 96)
	copy(out[:48], toBytes(&a[1]))
	copy(out[48:], toBytes(&a[0]))
	return out
}

func (e *fp2) new() *fe2 {
	return e.zero()
}

func (e *fp2) zero() *fe2 {
	return &fe2{}
}

func (e *fp2) one() *fe2 {
	return &fe2{*one(), *zero()}
}

func (e *fp2) rand(r io.Reader) (*fe2, error) {
	a0, err := newRand(r)
	if err != nil {
		return nil, err
	}
	a1, err := newRand(r)
	if err != nil {
		return nil, err
	}
	return &fe2{*a0, *a1}, nil
}

func (e *fp2) fromMont(c, a *fe2) {
	fromMont(&c[0], &a[0])
	fromMont(&c[1], &a[1])
}

func (e *fp2) isZero(a *fe2) bool {
	return isZero(&a[0]) && isZero(&a[1])
}

func (e *fp2) isOne(a *fe2) bool {
	return isOne(&a[0]) && isZero(&a[1])
}

func (e *fp2) equal(a, b *fe2) bool {
	return equal(&a[0], &b[0]) && equal(&a[1], &b[1])
}

func (e *fp2) copy(c, a *fe2) {
	c[0].Set(&a[0])
	c[1].Set(&a[1])
}

func (e *fp2) add(c, a, b *fe2) {
	add(&c[0], &a[0], &b[0])
	add(&c[1], &a[1], &b[1])
}

func (e *fp2) addAssign(a, b *fe2) {
	addAssign(&a[0], &b[0])
	addAssign(&a[1], &b[1])
}

func (e *fp2) ladd(c, a, b *fe2) {
	ladd(&c[0], &a[0], &b[0])
	ladd(&c[1], &a[1], &b[1])
}

func (e *fp2) double(c, a *fe2) {
	double(&c[0], &a[0])
	double(&c[1], &a[1])
}

func (e *fp2) doubleAssign(a *fe2) {
	doubleAssign(&a[0])
	doubleAssign(&a[1])
}

// ldouble doubles field element `a` and sets the result `c` without modular reduction
func (e *fp2) ldouble(c, a *fe2) {
	ldouble(&c[0], &a[0])
	ldouble(&c[1], &a[1])
}

func (e *fp2) sub(c, a, b *fe2) {
	sub(&c[0], &a[0], &b[0])
	sub(&c[1], &a[1], &b[1])
}

func (e *fp2) subAssign(c, a *fe2) {
	subAssign(&c[0], &a[0])
	subAssign(&c[1], &a[1])
}

func (e *fp2) neg(c, a *fe2) {
	neg(&c[0], &a[0])
	neg(&c[1], &a[1])
}

func (e *fp2) conjugate(c, a *fe2) {
	c[0].Set(&a[0])
	neg(&c[1], &a[1])
}

func (e *fp2) mul(c, a, b *fe2) {
	t := e.t
	mul(t[1], &a[0], &b[0])
	mul(t[2], &a[1], &b[1])
	add(t[0], &a[0], &a[1])
	add(t[3], &b[0], &b[1])
	sub(&c[0], t[1], t[2])
	addAssign(t[1], t[2])
	mulAssign(t[0], t[3])
	sub(&c[1], t[0], t[1])
}

func (e *fp2) mulAssign(a, b *fe2) {
	t := e.t
	mul(t[1], &a[0], &b[0])
	mul(t[2], &a[1], &b[1])
	add(t[0], &a[0], &a[1])
	add(t[3], &b[0], &b[1])
	sub(&a[0], t[1], t[2])
	addAssign(t[1], t[2])
	mulAssign(t[0], t[3])
	sub(&a[1], t[0], t[1])
}

func (e *fp2) square(c, a *fe2) {
	t := e.t
	ladd(t[0], &a[0], &a[1])
	sub(t[1], &a[0], &a[1])
	ldouble(t[2], &a[0])
	mul(&c[0], t[0], t[1])
	mul(&c[1], t[2], &a[1])
}

func (e *fp2) squareAssign(a *fe2) {
	t := e.t
	ladd(t[0], &a[0], &a[1])
	sub(t[1], &a[0], &a[1])
	ldouble(t[2], &a[0])
	mul(&a[0], t[0], t[1])
	mul(&a[1], t[2], &a[1])
}

func (e *fp2) mulByNonResidue(c, a *fe2) {
	t := e.t
	sub(t[0], &a[0], &a[1])
	add(&c[1], &a[0], &a[1])
	c[0].Set(t[0])
}

func (e *fp2) mulByB(c, a *fe2) {
	t := e.t
	double(t[0], &a[0])
	double(t[1], &a[1])
	doubleAssign(t[0])
	doubleAssign(t[1])
	sub(&c[0], t[0], t[1])
	add(&c[1], t[0], t[1])
}

func (e *fp2) inverse(c, a *fe2) {
	t := e.t
	square(t[0], &a[0])
	square(t[1], &a[1])
	addAssign(t[0], t[1])
	inverse(t[0], t[0])
	mul(&c[0], &a[0], t[0])
	mulAssign(t[0], &a[1])
	neg(&c[1], t[0])
}

func (e *fp2) mulByFq(c, a *fe2, b *fe) {
	mul(&c[0], &a[0], b)
	mul(&c[1], &a[1], b)
}

func (e *fp2) exp(c, a *fe2, s *big.Int) {
	z := e.one()
	for i := s.BitLen() - 1; i >= 0; i-- {
		e.square(z, z)
		if s.Bit(i) == 1 {
			e.mul(z, z, a)
		}
	}
	e.copy(c, z)
}

func (e *fp2) div(c, a, b *fe2) {
	t0 := e.new()
	e.inverse(t0, b)
	e.mul(c, a, t0)
}

func (e *fp2) frobeniousMap(c, a *fe2, power uint) {
	c[0].Set(&a[0])
	if power%2 == 1 {
		neg(&c[1], &a[1])
		return
	}
	c[1].Set(&a[1])
}

func (e *fp2) frobeniousMapAssign(a *fe2, power uint) {
	if power%2 == 1 {
		neg(&a[1], &a[1])
		return
	}
}

func (e *fp2) sqrt(c, a *fe2) bool {
	u, x0, a1, alpha := &fe2{}, &fe2{}, &fe2{}, &fe2{}
	e.copy(u, a)
	e.exp(a1, a, pMinus3Over4)
	e.square(alpha, a1)
	e.mul(alpha, alpha, a)
	e.mul(x0, a1, a)
	if e.equal(alpha, negativeOne2) {
		neg(&c[0], &x0[1])
		c[1].Set(&x0[0])
		return true
	}
	e.add(alpha, alpha, e.one())
	e.exp(alpha, alpha, pMinus1Over2)
	e.mul(c, alpha, x0)
	e.square(alpha, c)
	return e.equal(alpha, u)
}

// swuMap and isogenyMap methods are used for
// implementation of Simplified Shallue-van de Woestijne-Ulas Method
// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-05#section-6.6.2

func (e *fp2) isogenyMap(x, y *fe2) {
	params := isogenyConstantsG2
	degree := 3
	xNum, xDen, yNum, yDen := new(fe2), new(fe2), new(fe2), new(fe2)
	e.copy(xNum, params[0][degree])
	e.copy(xDen, params[1][degree])
	e.copy(yNum, params[2][degree])
	e.copy(yDen, params[3][degree])
	for i := degree - 1; i >= 0; i-- {
		e.mul(xNum, xNum, x)
		e.mul(xDen, xDen, x)
		e.mul(yNum, yNum, x)
		e.mul(yDen, yDen, x)
		e.add(xNum, xNum, params[0][i])
		e.add(xDen, xDen, params[1][i])
		e.add(yNum, yNum, params[2][i])
		e.add(yDen, yDen, params[3][i])
	}
	e.inverse(xDen, xDen)
	e.inverse(yDen, yDen)
	e.mul(xNum, xNum, xDen)
	e.mul(yNum, yNum, yDen)
	e.mul(yNum, yNum, y)
	e.copy(x, xNum)
	e.copy(y, yNum)
}

func (e *fp2) swuMap(u *fe2) (*fe2, *fe2, bool) {
	params := swuParamsForG2
	var tv [4]*fe2
	for i := 0; i < 4; i++ {
		tv[i] = e.new()
	}
	// 1.  tv1 = Z * u^2
	e.square(tv[0], u)
	e.mul(tv[0], tv[0], params.z)
	// 2.  tv2 = tv1^2
	e.square(tv[1], tv[0])
	// 3.   x1 = tv1 + tv2
	x1 := e.new()
	e.add(x1, tv[0], tv[1])
	// 4.   x1 = inv0(x1)
	e.inverse(x1, x1)
	// 5.   e1 = x1 == 0
	e1 := e.isZero(x1)
	// 6.   x1 = x1 + 1
	e.add(x1, x1, e.one())
	// 7.   x1 = CMOV(x1, c2, e1)    # If (tv1 + tv2) == 0, set x1 = -1 / Z
	if e1 {
		e.copy(x1, params.zInv)
	}
	// 8.   x1 = x1 * c1      # x1 = (-B / A) * (1 + (1 / (Z^2 * u^4 + Z * u^2)))
	e.mul(x1, x1, params.minusBOverA)
	// 9.  gx1 = x1^2
	gx1 := e.new()
	e.square(gx1, x1)
	// 10. gx1 = gx1 + A
	e.add(gx1, gx1, params.a) // TODO: a is zero we can ommit
	// 11. gx1 = gx1 * x1
	e.mul(gx1, gx1, x1)
	// 12. gx1 = gx1 + B             # gx1 = g(x1) = x1^3 + A * x1 + B
	e.add(gx1, gx1, params.b)
	// 13.  x2 = tv1 * x1            # x2 = Z * u^2 * x1
	x2 := e.new()
	e.mul(x2, tv[0], x1)
	// 14. tv2 = tv1 * tv2
	e.mul(tv[1], tv[0], tv[1])
	// 15. gx2 = gx1 * tv2           # gx2 = (Z * u^2)^3 * gx1
	gx2 := e.new()
	e.mul(gx2, gx1, tv[1])
	// 16.  e2 = is_square(gx1)
	// is quadratic non-residue
	isQuadraticNonResidue := func(elem *fe2) bool {
		// https://github.com/leovt/constructible/wiki/Taking-Square-Roots-in-quadratic-extension-Fields
		c0, c1 := new(fe), new(fe)
		square(c0, &elem[0])
		square(c1, &elem[1])
		mul(c1, c1, nonResidue1)
		neg(c1, c1)
		add(c1, c1, c0)
		return isQuadraticNonResidue(c1)
	}
	e2 := !isQuadraticNonResidue(gx1)
	// 17.   x = CMOV(x2, x1, e2)    # If is_square(gx1), x = x1, else x = x2
	x := e.new()
	if e2 {
		e.copy(x, x1)
	} else {
		e.copy(x, x2)
	}
	// 18.  y2 = CMOV(gx2, gx1, e2)  # If is_square(gx1), y2 = gx1, else y2 = gx2
	y2 := e.new()
	if e2 {
		e.copy(y2, gx1)
	} else {
		e.copy(y2, gx2)
	}
	// 19.   y = sqrt(y2)
	y := e.new()
	if hasSquareRoot := e.sqrt(y, y2); !hasSquareRoot {
		return nil, nil, false
	}
	// 20.  e3 = sgn0(u) == sgn0(y)  # Fix sign of y
	uSign := u.sign()
	ySign := y.sign()
	if ((uSign == 1 && ySign == -1) || (uSign == -1 && ySign == 1)) || ((uSign == 0 && ySign == -1) || (uSign == -1 && ySign == 0)) {
		e.neg(y, y)
	}
	return x, y, true
}
