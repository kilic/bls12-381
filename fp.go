package bls12381

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

// swuMap is implementation of Simplified Shallue-van de Woestijne-Ulas Method
// defined at draft-irtf-cfrg-hash-to-curve-06.
func swuMap(u *fe) (*fe, *fe) {
	// https: //tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-05#section-6.6.2
	var params = swuParamsForG1
	var tv [4]*fe
	for i := 0; i < 4; i++ {
		tv[i] = new(fe)
	}
	// 1.  tv1 = Z * u^2
	square(tv[0], u)
	mul(tv[0], tv[0], params.z)
	// 2.  tv2 = tv1^2
	square(tv[1], tv[0])
	// 3.   x1 = tv1 + tv2
	x1 := new(fe)
	add(x1, tv[0], tv[1])
	// 4.   x1 = inv0(x1)
	inverse(x1, x1)
	// 5.   e1 = x1 == 0
	e1 := isZero(x1)
	// 6.   x1 = x1 + 1
	add(x1, x1, one())
	// 7.   x1 = CMOV(x1, c2, e1)    # If (tv1 + tv2) == 0, set x1 = -1 / Z
	if e1 {
		x1.Set(params.zInv)
	}
	// 8.   x1 = x1 * c1      # x1 = (-B / A) * (1 + (1 / (Z^2 * u^4 + Z * u^2)))
	mul(x1, x1, params.minusBOverA)
	// 9.  gx1 = x1^2
	gx1 := new(fe)
	square(gx1, x1)
	// 10. gx1 = gx1 + A
	add(gx1, gx1, params.a) // TODO: a is zero we can ommit
	// 11. gx1 = gx1 * x1
	mul(gx1, gx1, x1)
	// 12. gx1 = gx1 + B             # gx1 = g(x1) = x1^3 + A * x1 + B
	add(gx1, gx1, params.b)
	// 13.  x2 = tv1 * x1            # x2 = Z * u^2 * x1
	x2 := new(fe)
	mul(x2, tv[0], x1)
	// 14. tv2 = tv1 * tv2
	mul(tv[1], tv[0], tv[1])
	// 15. gx2 = gx1 * tv2           # gx2 = (Z * u^2)^3 * gx1
	gx2 := new(fe)
	mul(gx2, gx1, tv[1])
	// 16.  e2 = is_square(gx1)
	e2 := !isQuadraticNonResidue(gx1)
	// 17.   x = CMOV(x2, x1, e2)    # If is_square(gx1), x = x1, else x = x2
	x := new(fe)
	if e2 {
		x.Set(x1)
	} else {
		x.Set(x2)
	}
	// 18.  y2 = CMOV(gx2, gx1, e2)  # If is_square(gx1), y2 = gx1, else y2 = gx2
	y2 := new(fe)
	if e2 {
		y2.Set(gx1)
	} else {
		y2.Set(gx2)
	}
	// 19.   y = sqrt(y2)
	y := new(fe)
	sqrt(y, y2)
	// 20.  e3 = sgn0(u) == sgn0(y)  # Fix sign of y
	if y.sign()^u.sign() != 0 {
		neg(y, y)
	}
	return x, y
}

// isogenyMap applies 11-isogeny map for BLS12-381 G1 defined at draft-irtf-cfrg-hash-to-curve-06.
func isogenyMap(x, y *fe) {
	// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#appendix-C.2
	params := isogenyConstansG1
	degree := 15
	xNum, xDen, yNum, yDen := new(fe), new(fe), new(fe), new(fe)
	xNum.Set(params[0][degree])
	xDen.Set(params[1][degree])
	yNum.Set(params[2][degree])
	yDen.Set(params[3][degree])
	for i := degree - 1; i >= 0; i-- {
		mul(xNum, xNum, x)
		mul(xDen, xDen, x)
		mul(yNum, yNum, x)
		mul(yDen, yDen, x)
		add(xNum, xNum, params[0][i])
		add(xDen, xDen, params[1][i])
		add(yNum, yNum, params[2][i])
		add(yDen, yDen, params[3][i])
	}
	inverse(xDen, xDen)
	inverse(yDen, yDen)
	mul(xNum, xNum, xDen)
	mul(yNum, yNum, yDen)
	mul(yNum, yNum, y)
	x.Set(xNum)
	y.Set(yNum)
}
