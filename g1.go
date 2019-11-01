package bls

import (
	"fmt"
	"math/big"
)

type PointG1 [3]fe

func (p *PointG1) Set(p2 *PointG1) *PointG1 {
	p[0].Set(&p2[0])
	p[1].Set(&p2[1])
	p[2].Set(&p2[2])
	return p
}

type G1 struct {
	f *fp
	t [9]*fe
}

func NewG1(f *fp) *G1 {
	t := [9]*fe{}
	for i := 0; i < 9; i++ {
		t[i] = f.zero()
	}
	if f == nil {
		f = newFp()
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
	if in[0]&(1<<6) != 0 {
		for i, v := range in {
			if (i == 0 && v != 0x40) || (i != 0 && v != 0x00) {
				return nil, fmt.Errorf("input string should be zero when infinity flag is set")
			}
		}
		return g.Zero(), nil
	}
	in[0] &= 0x1f
	x, y := &fe{}, &fe{}
	if err := g.f.newElementFromBytes(x, in[:48]); err != nil {
		return nil, err
	}
	if err := g.f.newElementFromBytes(y, in[48:]); err != nil {
		return nil, err
	}
	p := &PointG1{}
	g.f.copy(&p[0], x)
	g.f.copy(&p[1], y)
	g.f.copy(&p[2], &fpOne)
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
	copy(out[:48], g.f.toBytes(&p[0]))
	copy(out[48:], g.f.toBytes(&p[1]))
	return out
}

func (g *G1) FromCompressed(compressed []byte) (*PointG1, error) {
	if len(compressed) < 48 {
		return nil, fmt.Errorf("input string should be equal or larger than 48")
	}
	var in [48]byte
	copy(in[:], compressed[:])
	if in[0]&(1<<7) == 0 {
		return nil, fmt.Errorf("compression flag should be set")
	}
	if in[0]&(1<<6) != 0 {
		// in[0] == (1 << 6) + (1 << 7)
		for i, v := range in {
			if (i == 0 && v != 0xc0) || (i != 0 && v != 0x00) {
				return nil, fmt.Errorf("input string should be zero when infinity flag is set")
			}
		}
		return g.Zero(), nil
	}
	a := in[0]&(1<<5) != 0
	in[0] &= 0x1f
	x := &fe{}
	if err := g.f.newElementFromBytes(x, in[:]); err != nil {
		return nil, err
	}
	// solve curve equation
	y := &fe{}
	g.f.square(y, x)
	g.f.mul(y, y, x)
	g.f.add(y, y, b)
	if ok := g.f.sqrt(y, y); !ok {
		return nil, fmt.Errorf("point is not on curve")
	}
	// select lexicographically, should be in normalized form
	negY, negYn, yn := &fe{}, &fe{}, &fe{}
	g.f.demont(yn, y)
	g.f.neg(negY, y)
	g.f.neg(negYn, yn)
	if yn.Cmp(negYn) > -1 != a {
		g.f.copy(y, negY)
	}
	p := &PointG1{}
	g.f.copy(&p[0], x)
	g.f.copy(&p[1], y)
	g.f.copy(&p[2], &fpOne)
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
		copy(out[:], g.f.toBytes(&p[0]))
		y, negY := &fe{}, &fe{}
		g.f.copy(y, &p[1])
		g.f.demont(y, y)
		g.f.neg(negY, y)
		if y.Cmp(negY) > 0 {
			out[0] |= 1 << 5
		}
	}
	out[0] |= 1 << 7
	return out
}

func (g *G1) fromRawUnchecked(in []byte) *PointG1 {
	p := &PointG1{}
	if err := g.f.newElementFromBytes(&p[0], in[:48]); err != nil {
		panic(err)
	}
	if err := g.f.newElementFromBytes(&p[1], in[48:]); err != nil {
		panic(err)
	}
	g.f.copy(&p[2], &fpOne)
	return p
}

func (g *G1) isTorsionFree(p *PointG1) bool {
	tmp := &PointG1{}
	g.MulScalar(tmp, p, q)
	return g.IsZero(tmp)
}

func (g *G1) Zero() *PointG1 {
	return &PointG1{
		*g.f.zero(),
		*g.f.one(),
		*g.f.zero(),
	}
}

func (g *G1) NegativeOne() *PointG1 {
	return g.Copy(&PointG1{}, &g1NegativeOne)
}

func (g *G1) One() *PointG1 {
	return g.Copy(&PointG1{}, &g1One)
}

func (g *G1) Copy(dst *PointG1, src *PointG1) *PointG1 {
	return dst.Set(src)
}

func (g *G1) IsZero(p *PointG1) bool {
	return g.f.isZero(&p[2])
}

func (g *G1) Equal(p1, p2 *PointG1) bool {
	if g.IsZero(p1) {
		return g.IsZero(p2)
	}
	if g.IsZero(p2) {
		return g.IsZero(p1)
	}
	t := g.t
	g.f.square(t[0], &p1[2])
	g.f.square(t[1], &p2[2])
	g.f.mul(t[2], t[0], &p2[0])
	g.f.mul(t[3], t[1], &p1[0])
	g.f.mul(t[0], t[0], &p1[2])
	g.f.mul(t[1], t[1], &p2[2])
	g.f.mul(t[1], t[1], &p1[1])
	g.f.mul(t[0], t[0], &p2[1])
	return g.f.equal(t[0], t[1]) && g.f.equal(t[2], t[3])
}

func (g *G1) IsOnCurve(p *PointG1) bool {
	if g.IsZero(p) {
		return true
	}
	t := g.t
	g.f.square(t[0], &p[1])
	g.f.square(t[1], &p[0])
	g.f.mul(t[1], t[1], &p[0])
	g.f.square(t[2], &p[2])
	g.f.square(t[3], t[2])
	g.f.mul(t[2], t[2], t[3])
	g.f.mul(t[2], b, t[2])
	g.f.add(t[1], t[1], t[2])
	return g.f.equal(t[0], t[1])
}

func (g *G1) IsAffine(p *PointG1) bool {
	return g.f.equal(&p[2], &fpOne)
}

func (g *G1) Affine(p *PointG1) {
	if g.IsZero(p) {
		return
	}
	if !g.IsAffine(p) {
		t := g.t
		g.f.inverse(t[0], &p[2])
		g.f.square(t[1], t[0])
		g.f.mul(&p[0], &p[0], t[1])
		g.f.mul(t[0], t[0], t[1])
		g.f.mul(&p[1], &p[1], t[0])
		g.f.copy(&p[2], g.f.one())
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
	g.f.square(t[7], &p1[2])
	g.f.mul(t[1], &p2[0], t[7])
	g.f.mul(t[2], &p1[2], t[7])
	g.f.mul(t[0], &p2[1], t[2])
	g.f.square(t[8], &p2[2])
	g.f.mul(t[3], &p1[0], t[8])
	g.f.mul(t[4], &p2[2], t[8])
	g.f.mul(t[2], &p1[1], t[4])
	if g.f.equal(t[1], t[3]) {
		if g.f.equal(t[0], t[2]) {
			return g.Double(r, p1)
		} else {
			return g.Copy(r, infinity)
		}
	}
	g.f.sub(t[1], t[1], t[3])
	g.f.double(t[4], t[1])
	g.f.square(t[4], t[4])
	g.f.mul(t[5], t[1], t[4])
	g.f.sub(t[0], t[0], t[2])
	g.f.double(t[0], t[0])
	g.f.square(t[6], t[0])
	g.f.sub(t[6], t[6], t[5])
	g.f.mul(t[3], t[3], t[4])
	g.f.double(t[4], t[3])
	g.f.sub(&r[0], t[6], t[4])
	g.f.sub(t[4], t[3], &r[0])
	g.f.mul(t[6], t[2], t[5])
	g.f.double(t[6], t[6])
	g.f.mul(t[0], t[0], t[4])
	g.f.sub(&r[1], t[0], t[6])
	g.f.add(t[0], &p1[2], &p2[2])
	g.f.square(t[0], t[0])
	g.f.sub(t[0], t[0], t[7])
	g.f.sub(t[0], t[0], t[8])
	g.f.mul(&r[2], t[0], t[1])
	return r
}

func (g *G1) Double(r, p *PointG1) *PointG1 {
	if g.IsZero(p) {
		g.Copy(r, p)
		return r
	}
	t := g.t
	g.f.square(t[0], &p[0])
	g.f.square(t[1], &p[1])
	g.f.square(t[2], t[1])
	g.f.add(t[1], &p[0], t[1])
	g.f.square(t[1], t[1])
	g.f.sub(t[1], t[1], t[0])
	g.f.sub(t[1], t[1], t[2])
	g.f.double(t[1], t[1])
	g.f.double(t[3], t[0])
	g.f.add(t[0], t[3], t[0])
	g.f.square(t[4], t[0])
	g.f.double(t[3], t[1])
	g.f.sub(&r[0], t[4], t[3])
	g.f.sub(t[1], t[1], &r[0])
	g.f.double(t[2], t[2])
	g.f.double(t[2], t[2])
	g.f.double(t[2], t[2])
	g.f.mul(t[0], t[0], t[1])
	g.f.sub(t[1], t[0], t[2])
	g.f.mul(t[0], &p[1], &p[2])
	g.f.copy(&r[1], t[1])
	g.f.double(&r[2], t[0])
	return r
}

func (g *G1) Neg(r, p *PointG1) *PointG1 {
	g.f.copy(&r[0], &p[0])
	g.f.neg(&r[1], &p[1])
	g.f.copy(&r[2], &p[2])
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
	g.MulScalar(c, p, cofactorG1)
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
