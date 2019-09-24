package bls

import (
	"fmt"
	"math/big"
)

type PointG1 [3]Fe

func (p *PointG1) Set(p2 *PointG1) *PointG1 {
	p[0].Set(&p2[0])
	p[1].Set(&p2[1])
	p[2].Set(&p2[2])
	return p
}

type G1 struct {
	f *Fp
	t [9]*Fe
}

func NewG1(f *Fp) *G1 {
	t := [9]*Fe{}
	for i := 0; i < 9; i++ {
		t[i] = f.Zero()
	}
	return &G1{
		f: f,
		t: t,
	}
}

func (g *G1) FromUncompressed(uncompressed []byte) (*PointG1, error) {
	if len(uncompressed) < 96 {
		return nil, fmt.Errorf("input string should be equal or larger than 96")
	}
	var in [96]byte
	copy(in[:], uncompressed[:96])
	if in[0]&(1<<7) != 0 {
		return nil, fmt.Errorf("compression flag should be zero")
	}
	if in[0]&(1<<5) != 0 {
		return nil, fmt.Errorf("sort flag should be zero")
	}
	in[0] &= 0x1f
	x, y := &Fe{}, &Fe{}
	if err := g.f.NewElementFromBytes(x, in[:48]); err != nil {
		return nil, err
	}
	if err := g.f.NewElementFromBytes(y, in[48:]); err != nil {
		return nil, err
	}
	if in[0]&(1<<6) == 1 {
		if !(x.IsZero() && y.IsZero()) {
			return nil, fmt.Errorf("input string should be zero when infinity flag is set")
		}
	}
	p := &PointG1{}
	g.f.Copy(&p[0], x)
	g.f.Copy(&p[1], y)
	g.f.Copy(&p[2], &fpOne)
	if !g.IsOnCurve(p) {
		return nil, fmt.Errorf("point is not on curve")
	}
	if !g.isTorsionFree(p) {
		return nil, fmt.Errorf("point is not on correct subgroup")
	}
	return p, nil
}

func (g *G1) ToUncompressed(p *PointG1) []byte {
	out := make([]byte, 96)
	g.Affine(p)
	if g.IsZero(p) {
		out[0] |= 1 << 6
	}
	copy(out[:48], g.f.ToBytes(&p[0]))
	copy(out[48:], g.f.ToBytes(&p[1]))
	return out
}

func (g *G1) FromCompressed(compressed []byte) (*PointG1, error) {
	if len(compressed) < 48 {
		return nil, fmt.Errorf("input string should be equal or larger than 48")
	}
	var in [48]byte
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
	x := &Fe{}
	if err := g.f.NewElementFromBytes(x, in[:]); err != nil {
		return nil, err
	}
	// solve curve equation
	y := &Fe{}
	g.f.Square(y, x)
	g.f.Mul(y, y, x)
	g.f.Add(y, y, b)
	if ok := g.f.Sqrt(y, y); !ok {
		return nil, fmt.Errorf("point is not on curve")
	}
	// select lexicographically, should be in mont reduced form
	negY, negYn, yn := &Fe{}, &Fe{}, &Fe{}
	g.f.Demont(yn, y)
	g.f.Neg(negY, y)
	g.f.Neg(negYn, yn)
	if yn.Cmp(negYn) > -1 != a {
		g.f.Copy(y, negY)
	}
	p := &PointG1{}
	g.f.Copy(&p[0], x)
	g.f.Copy(&p[1], y)
	g.f.Copy(&p[2], &fpOne)
	if !g.isTorsionFree(p) {
		return nil, fmt.Errorf("point is not on correct subgroup")
	}
	return p, nil
}

func (g *G1) ToCompressed(p *PointG1) []byte {
	out := make([]byte, 48)
	g.Affine(p)
	if g.IsZero(p) {
		out[0] |= 1 << 6
	} else {
		copy(out[:], g.f.ToBytes(&p[0]))
		y, negY := &Fe{}, &Fe{}
		g.f.Copy(y, &p[1])
		g.f.Demont(y, y)
		g.f.Neg(negY, y)
		if y.Cmp(negY) > 0 {
			out[0] |= 1 << 5
		}
	}
	out[0] |= 1 << 7
	return out
}

func (g *G1) fromRawUnchecked(in []byte) *PointG1 {
	p := &PointG1{}
	if err := g.f.NewElementFromBytes(&p[0], in[:48]); err != nil {
		panic(err)
	}
	if err := g.f.NewElementFromBytes(&p[1], in[48:]); err != nil {
		panic(err)
	}
	g.f.Copy(&p[2], &fpOne)
	return p
}

func (g *G1) isTorsionFree(p *PointG1) bool {
	tmp := &PointG1{}
	g.MulScalar(tmp, p, q)
	return g.IsZero(tmp)
}

func (g *G1) Zero() *PointG1 {
	return &PointG1{
		*g.f.Zero(),
		*g.f.One(),
		*g.f.Zero(),
	}
}

func (g *G1) Copy(dst *PointG1, src *PointG1) *PointG1 {
	return dst.Set(src)
}

func (g *G1) IsZero(p *PointG1) bool {
	return g.f.IsZero(&p[2])
}

func (g *G1) Equal(p1, p2 *PointG1) bool {
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

func (g *G1) IsOnCurve(p *PointG1) bool {
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
	g.f.Mul(t[2], b, t[2])
	g.f.Add(t[1], t[1], t[2])
	return g.f.Equal(t[0], t[1])
}

func (g *G1) IsAffine(p *PointG1) bool {
	return g.f.Equal(&p[2], &fpOne)
}

func (g *G1) Affine(p *PointG1) {
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

func (g *G1) Add(r, p1, p2 *PointG1) *PointG1 {
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
			return g.Copy(r, infinity)
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

func (g *G1) Double(r, p *PointG1) *PointG1 {
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

func (g *G1) Neg(r, p *PointG1) *PointG1 {
	g.f.Copy(&r[0], &p[0])
	g.f.Neg(&r[1], &p[1])
	g.f.Copy(&r[2], &p[2])
	return r
}

func (g *G1) Sub(c, a, b *PointG1) *PointG1 {
	d := &PointG1{}
	g.Neg(d, b)
	g.Add(c, a, d)
	return c
}

// negates second operand
func (g *G1) SubUnsafe(c, a, b *PointG1) *PointG1 {
	g.Neg(b, b)
	g.Add(c, a, b)
	return c
}

func (g *G1) MulScalar(c, p *PointG1, e *big.Int) *PointG1 {
	q, n := &PointG1{}, &PointG1{}
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

func (g *G1) MulByCofactor(c, p *PointG1) {
	g.MulScalar(c, p, cofactor)
}

// func (g *G1) MultiExp(r *PointG1, points []*PointG1, powers []*big.Int) (*PointG1, error) {
// 	if len(points) != len(powers) {
// 		return nil, fmt.Errorf("point and scalar vectors should be in same length")
// 	}
// 	var c uint = 3
// 	if len(powers) > 32 {
// 		c = uint(math.Ceil(math.Log10(float64(len(powers)))))
// 	}
// 	bucket_size, numBits := (1<<c)-1, q.BitLen()
// 	windows := make([]PointG1, numBits/int(c)+1)
// 	bucket := make([]PointG1, bucket_size)
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
