package bls

import (
	"fmt"
	"math/big"
)

type PointG2 [3]Fe2

func (p *PointG2) Set(p2 *PointG2) *PointG2 {
	p[0][0].Set(&p2[0][0])
	p[1][1].Set(&p2[1][1])
	p[2][0].Set(&p2[2][0])
	p[0][1].Set(&p2[0][1])
	p[1][0].Set(&p2[1][0])
	p[2][1].Set(&p2[2][1])
	return p
}

// VerifyMsgWithDomain verifies that a message with a specified domain
// against a g1 point.
func (g *PointG2) VerifyMsgWithDomain(msg [32]byte, p *PointG1, domain [8]byte) bool {
	e := NewBLSPairingEngine()
	target := &Fe12{}
	e.Pair(target,
		[]PointG1{
			G1NegativeOne,
			*p,
		},
		[]PointG2{
			*g,
			*HashToG2WithDomain(msg, domain),
		},
	)
	return e.Fp12.Equal(&Fp12One, target)
}

type G2 struct {
	f *Fp2
	t [9]*Fe2
}

func NewG2(f *Fp2) *G2 {
	t := [9]*Fe2{}
	for i := 0; i < 9; i++ {
		t[i] = f.Zero()
	}
	if f == nil {
		return &G2{
			f: NewFp2(nil),
			t: t,
		}
	}
	return &G2{
		f: f,
		t: t,
	}
}

func (g *G2) FromUncompressed(uncompressed []byte) (*PointG2, error) {
	if len(uncompressed) < 192 {
		return nil, fmt.Errorf("input string should be equal or larger than 192")
	}
	var in [192]byte
	copy(in[:], uncompressed[:192])
	if in[0]&(1<<7) != 0 {
		return nil, fmt.Errorf("compression flag should be zero")
	}
	if in[0]&(1<<5) != 0 {
		return nil, fmt.Errorf("sort flag should be zero")
	}
	in[0] &= 0x1f
	x, y := &Fe2{}, &Fe2{}
	if err := g.f.NewElementFromBytes(x, in[:96]); err != nil {
		return nil, err
	}
	if err := g.f.NewElementFromBytes(y, in[96:]); err != nil {
		return nil, err
	}
	if in[0]&(1<<6) == 1 {
		if !(g.f.IsZero(x) && g.f.IsZero(y)) {
			return nil, fmt.Errorf("input string should be zero when infinity flag is set")
		}
	}
	p := &PointG2{}
	g.f.Copy(&p[0], x)
	g.f.Copy(&p[1], y)
	g.f.Copy(&p[2], &Fp2One)
	if !g.IsOnCurve(p) {
		return nil, fmt.Errorf("point is not on curve")
	}
	if !g.isTorsionFree(p) {
		return nil, fmt.Errorf("point is not on correct subgroup")
	}
	return p, nil
}

func (g *G2) ToUncompressed(p *PointG2) []byte {
	out := make([]byte, 192)
	g.Affine(p)
	if g.IsZero(p) {
		out[0] |= 1 << 6
	}
	copy(out[:96], g.f.ToBytes(&p[0]))
	copy(out[96:], g.f.ToBytes(&p[1]))
	return out
}

func (g *G2) FromCompressed(compressed []byte) (*PointG2, error) {
	if len(compressed) < 96 {
		return nil, fmt.Errorf("input string should be equal or larger than 96")
	}
	var in [96]byte
	copy(in[:], compressed[:])
	if in[0]&(1<<7) == 0 {
		return nil, fmt.Errorf("bad compression")
	}
	if in[0]&(1<<6) != 0 {
		in[0] &= 0x3f
		for i := 0; i < 48; i++ {
			if in[i] != 0 {
				return nil, fmt.Errorf("bad infinity compression")
			}
		}
		return g.Zero(), nil
	}
	a := in[0]&(1<<5) != 0
	in[0] &= 0x1f
	x := &Fe2{}
	if err := g.f.NewElementFromBytes(x, in[:]); err != nil {
		return nil, err
	}
	// solve curve equation
	y := &Fe2{}
	g.f.Square(y, x)
	g.f.Mul(y, y, x)
	g.f.Add(y, y, b2)
	if ok := g.f.Sqrt(y, y); !ok {
		return nil, fmt.Errorf("point is not on curve")
	}
	// select lexicographically, should be in mont reduced form
	negYn, negY, yn := &Fe2{}, &Fe2{}, &Fe2{}
	g.f.f.Demont(&yn[0], &y[0])
	g.f.f.Demont(&yn[1], &y[1])
	g.f.Neg(negY, y)
	g.f.Neg(negYn, yn)
	if (yn[1].Cmp(&negYn[1]) > 0 != a) || (yn[1].IsZero() && yn[0].Cmp(&negYn[0]) > 0 != a) {
		g.f.Copy(y, negY)
	}
	p := &PointG2{}
	g.f.Copy(&p[0], x)
	g.f.Copy(&p[1], y)
	g.f.Copy(&p[2], &Fp2One)
	if !g.isTorsionFree(p) {
		return nil, fmt.Errorf("point is not on correct subgroup")
	}
	return p, nil
}

func (g *G2) ToCompressed(p *PointG2) []byte {
	out := make([]byte, 96)
	g.Affine(p)
	if g.IsZero(p) {
		out[0] |= 1 << 6
	} else {
		copy(out[:], g.f.ToBytes(&p[0]))
		y, negY := &Fe2{}, &Fe2{}
		g.f.Copy(y, &p[1])
		g.f.f.Demont(&y[0], &y[0])
		g.f.f.Demont(&y[1], &y[1])
		g.f.Neg(negY, y)
		if (y[1].Cmp(&negY[1]) > 0) || (y[1].IsZero() && y[1].Cmp(&negY[1]) > 0) {
			out[0] |= 1 << 5
		}
	}
	out[0] |= 1 << 7
	return out
}

func (g *G2) fromRawUnchecked(in []byte) *PointG2 {
	p := &PointG2{}
	if err := g.f.NewElementFromBytes(&p[0], in[:96]); err != nil {
		panic(err)
	}
	if err := g.f.NewElementFromBytes(&p[1], in[96:]); err != nil {
		panic(err)
	}
	g.f.Copy(&p[2], &Fp2One)
	return p
}

func (g *G2) isTorsionFree(p *PointG2) bool {
	tmp := &PointG2{}
	g.MulScalar(tmp, p, q)
	return g.IsZero(tmp)
}

func (g *G2) Zero() *PointG2 {
	return &PointG2{
		*g.f.Zero(),
		*g.f.One(),
		*g.f.Zero(),
	}
}

func (g *G2) Copy(dst *PointG2, src *PointG2) *PointG2 {
	return dst.Set(src)
}

func (g *G2) IsZero(p *PointG2) bool {
	return g.f.IsZero(&p[2])
}

func (g *G2) Equal(p1, p2 *PointG2) bool {
	if g.IsZero(p1) {
		return g.IsZero(p2)
	}
	if g.IsZero(p2) {
		return g.IsZero(p1)
	}
	t := g.t
	g.f.Square(t[0], &p1[2])
	g.f.Square(t[1], &p2[2])
	g.f.Mul(t[2], t[0], &p2[0])
	g.f.Mul(t[3], t[1], &p1[0])
	g.f.Mul(t[0], t[0], &p1[2])
	g.f.Mul(t[1], t[1], &p2[2])
	g.f.Mul(t[1], t[1], &p1[1])
	g.f.Mul(t[0], t[0], &p2[1])
	return g.f.Equal(t[0], t[1]) && g.f.Equal(t[2], t[3])
}

func (g *G2) IsOnCurve(p *PointG2) bool {
	if g.IsZero(p) {
		return true
	}
	t := g.t
	g.f.Square(t[0], &p[1])
	g.f.Square(t[1], &p[0])
	g.f.Mul(t[1], t[1], &p[0])
	g.f.Square(t[2], &p[2])
	g.f.Square(t[3], t[2])
	g.f.Mul(t[2], t[2], t[3])
	g.f.Mul(t[2], b2, t[2])
	g.f.Add(t[1], t[1], t[2])
	return g.f.Equal(t[0], t[1])
}

func (g *G2) IsAffine(p *PointG2) bool {
	return g.f.Equal(&p[2], &Fp2One)
}

func (g *G2) Affine(p *PointG2) {
	if g.IsZero(p) {
		return
	}
	if !g.IsAffine(p) {
		t := g.t
		g.f.Inverse(t[0], &p[2])
		g.f.Square(t[1], t[0])
		g.f.Mul(&p[0], &p[0], t[1])
		g.f.Mul(t[0], t[0], t[1])
		g.f.Mul(&p[1], &p[1], t[0])
		g.f.Copy(&p[2], g.f.One())
	}
}

func (g *G2) Add(r, p1, p2 *PointG2) *PointG2 {
	if g.IsZero(p1) {
		g.Copy(r, p2)
		return r
	}
	if g.IsZero(p2) {
		g.Copy(r, p1)
		return r
	}
	t := g.t
	g.f.Square(t[7], &p1[2])
	g.f.Mul(t[1], &p2[0], t[7])
	g.f.Mul(t[2], &p1[2], t[7])
	g.f.Mul(t[0], &p2[1], t[2])
	g.f.Square(t[8], &p2[2])
	g.f.Mul(t[3], &p1[0], t[8])
	g.f.Mul(t[4], &p2[2], t[8])
	g.f.Mul(t[2], &p1[1], t[4])
	if g.f.Equal(t[1], t[3]) {
		if g.f.Equal(t[0], t[2]) {
			return g.Double(r, p1)
		} else {
			return g.Copy(r, infinity2)
		}
	}
	g.f.Sub(t[1], t[1], t[3])
	g.f.Double(t[4], t[1])
	g.f.Square(t[4], t[4])
	g.f.Mul(t[5], t[1], t[4])
	g.f.Sub(t[0], t[0], t[2])
	g.f.Double(t[0], t[0])
	g.f.Square(t[6], t[0])
	g.f.Sub(t[6], t[6], t[5])
	g.f.Mul(t[3], t[3], t[4])
	g.f.Double(t[4], t[3])
	g.f.Sub(&r[0], t[6], t[4])
	g.f.Sub(t[4], t[3], &r[0])
	g.f.Mul(t[6], t[2], t[5])
	g.f.Double(t[6], t[6])
	g.f.Mul(t[0], t[0], t[4])
	g.f.Sub(&r[1], t[0], t[6])
	g.f.Add(t[0], &p1[2], &p2[2])
	g.f.Square(t[0], t[0])
	g.f.Sub(t[0], t[0], t[7])
	g.f.Sub(t[0], t[0], t[8])
	g.f.Mul(&r[2], t[0], t[1])
	return r
}

func (g *G2) Double(r, p *PointG2) *PointG2 {
	if g.IsZero(p) {
		g.Copy(r, p)
		return r
	}
	t := g.t
	g.f.Square(t[0], &p[0])
	g.f.Square(t[1], &p[1])
	g.f.Square(t[2], t[1])
	g.f.Add(t[1], &p[0], t[1])
	g.f.Square(t[1], t[1])
	g.f.Sub(t[1], t[1], t[0])
	g.f.Sub(t[1], t[1], t[2])
	g.f.Double(t[1], t[1])
	g.f.Double(t[3], t[0])
	g.f.Add(t[0], t[3], t[0])
	g.f.Square(t[4], t[0])
	g.f.Double(t[3], t[1])
	g.f.Sub(&r[0], t[4], t[3])
	g.f.Sub(t[1], t[1], &r[0])
	g.f.Double(t[2], t[2])
	g.f.Double(t[2], t[2])
	g.f.Double(t[2], t[2])
	g.f.Mul(t[0], t[0], t[1])
	g.f.Sub(t[1], t[0], t[2])
	g.f.Mul(t[0], &p[1], &p[2])
	g.f.Copy(&r[1], t[1])
	g.f.Double(&r[2], t[0])
	return r
}

func (g *G2) Neg(r, p *PointG2) *PointG2 {
	g.f.Copy(&r[0], &p[0])
	g.f.Neg(&r[1], &p[1])
	g.f.Copy(&r[2], &p[2])
	return r
}

func (g *G2) Sub(c, a, b *PointG2) *PointG2 {
	d := &PointG2{}
	g.Neg(d, b)
	g.Add(c, a, d)
	return c
}

// negates second operand
func (g *G2) SubUnsafe(c, a, b *PointG2) *PointG2 {
	g.Neg(b, b)
	g.Add(c, a, b)
	return c
}

func (g *G2) MulScalar(c, p *PointG2, e *big.Int) *PointG2 {
	q, n := &PointG2{}, &PointG2{}
	g.Copy(n, p)
	l := e.BitLen()
	for i := 0; i < l; i++ {
		if e.Bit(i) == 1 {
			g.Add(q, q, n)
		}
		g.Double(n, n)
	}
	return g.Copy(c, q)
}

func (g *G2) MulByCofactor(c, p *PointG2) {
	g.MulScalar(c, p, cofactorG2)
}

func (g *G2) MapToPoint(in []byte) *PointG2 {
	x, y := &Fe2{}, &Fe2{}
	fp2 := g.f
	fp := fp2.f
	err := fp2.NewElementFromBytes(x, in)
	if err != nil {
		panic(err)
	}
	for {
		fp2.Square(y, x)
		fp2.Mul(y, y, x)
		fp2.Add(y, y, b2)
		if ok := fp2.Sqrt(y, y); ok {
			// favour negative y
			negYn, negY, yn := &Fe2{}, &Fe2{}, &Fe2{}
			fp.Demont(&yn[0], &y[0])
			fp.Demont(&yn[1], &y[1])
			fp2.Neg(negY, y)
			fp2.Neg(negYn, yn)
			if yn[1].Cmp(&negYn[1]) > 0 || (yn[1].IsZero() && yn[0].Cmp(&negYn[0]) > 0) {
				fp2.Copy(y, y)
			} else {
				fp2.Copy(y, negY)
			}
			p := &PointG2{*x, *y, Fp2One}
			g.MulByCofactor(p, p)
			return p
		}
		fp2.Add(x, x, &Fp2One)
	}
}

// func (g *G2) MultiExp(r *PointG2, points []*PointG2, powers []*big.Int) (*PointG2, error) {
// 	if len(points) != len(powers) {
// 		return nil, fmt.Errorf("point and scalar vectors should be in same length")
// 	}
// 	var c uint = 3
// 	if len(powers) > 32 {
// 		c = uint(math.Ceil(math.Log10(float64(len(powers)))))
// 	}
// 	bucket_size, numBits := (1<<c)-1, q.BitLen()
// 	windows := make([]PointG2, numBits/int(c)+1)
// 	bucket := make([]PointG2, bucket_size)
// 	acc, sum, zero := g.Zero(), g.Zero(), g.Zero()
// 	s := new(big.Int)
// 	for i, m := 0, 0; i <= numBits; i, m = i+int(c), m+1 {
// 		for i := 0; i < bucket_size; i++ {
// 			g.Copy(&bucket[i], zero)
// 		}
// 		for j := 0; j < len(powers); j++ {
// 			s = powers[j]
// 			index := s.Uint64() & uint64(bucket_size)
// 			if index != 0 {
// 				g.Add(&bucket[index-1], &bucket[index-1], points[j])
// 			}
// 			s.Rsh(s, c)
// 		}
// 		g.Copy(acc, zero)
// 		g.Copy(sum, zero)
// 		for k := bucket_size - 1; k >= 0; k-- {
// 			g.Add(sum, sum, &bucket[k])
// 			g.Add(acc, acc, sum)
// 		}
// 		g.Copy(&windows[m], acc)
// 	}
// 	g.Copy(acc, zero)
// 	for i := len(windows) - 1; i >= 0; i-- {
// 		for j := 0; j < int(c); j++ {
// 			g.Double(acc, acc)
// 		}
// 		g.Add(acc, acc, &windows[i])
// 	}
// 	g.Copy(r, acc)
// 	return r, nil
// }
