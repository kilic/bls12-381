package bls

import (
	"fmt"
	"math"
	"math/big"
)

type PointG1 [3]fe

func (p *PointG1) Set(p2 *PointG1) *PointG1 {
	p[0].Set(&p2[0])
	p[1].Set(&p2[1])
	p[2].Set(&p2[2])
	return p
}

type tempG1 struct {
	t [9]*fe
}

type G1 struct {
	tempG1
}

func NewG1() *G1 {
	cfgArch()
	t := newTempG1()
	return &G1{t}
}

func newTempG1() tempG1 {
	t := [9]*fe{}
	for i := 0; i < 9; i++ {
		t[i] = &fe{}
	}
	return tempG1{t}
}

func (g *G1) Q() *big.Int {
	return new(big.Int).Set(q)
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
	x, err := fromBytes(in[:48])
	if err != nil {
		return nil, err
	}
	y, err := fromBytes(in[48:])
	if err != nil {
		return nil, err
	}
	z := one()
	p := &PointG1{*x, *y, *z}
	if !g.IsOnCurve(p) {
		return nil, fmt.Errorf("point is not on curve")
	}
	if !g.InCorrectSubgroup(p) {
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
	copy(out[:48], toBytes(&p[0]))
	copy(out[48:], toBytes(&p[1]))
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
	x, err := fromBytes(in[:])
	if err != nil {
		return nil, err
	}
	// solve curve equation
	y := &fe{}
	square(y, x)
	mul(y, y, x)
	add(y, y, b)
	if ok := sqrt(y, y); !ok {
		return nil, fmt.Errorf("point is not on curve")
	}
	// select lexicographically, should be in normalized form
	negY, negYn, yn := &fe{}, &fe{}, &fe{}
	fromMont(yn, y)
	neg(negY, y)
	neg(negYn, yn)
	if yn.Cmp(negYn) > -1 != a {
		y.Set(negY)
	}
	z := one()
	p := &PointG1{*x, *y, *z}
	if !g.InCorrectSubgroup(p) {
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
		copy(out[:], toBytes(&p[0]))
		y, negY := new(fe).Set(&p[1]), new(fe)
		fromMont(y, y)
		neg(negY, y)
		if y.Cmp(negY) > 0 {
			out[0] |= 1 << 5
		}
	}
	out[0] |= 1 << 7
	return out
}

func (g *G1) New() *PointG1 {
	return g.Zero()
}

func (g *G1) fromRawUnchecked(in []byte) (*PointG1, error) {
	p0, err := fromBytes(in[:48])
	if err != nil {
		return nil, err
	}
	p1, err := fromBytes(in[48:])
	if err != nil {
		panic(err)
	}
	p2 := one()
	return &PointG1{*p0, *p1, *p2}, nil
}

func (g *G1) InCorrectSubgroup(p *PointG1) bool {
	tmp := &PointG1{}
	g.MulScalar(tmp, p, q)
	return g.IsZero(tmp)
}

func (g *G1) Zero() *PointG1 {
	return &PointG1{
		*zero(),
		*one(),
		*zero(),
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
	return isZero(&p[2])
}

func (g *G1) Equal(p1, p2 *PointG1) bool {
	if g.IsZero(p1) {
		return g.IsZero(p2)
	}
	if g.IsZero(p2) {
		return g.IsZero(p1)
	}
	t := g.t
	square(t[0], &p1[2])
	square(t[1], &p2[2])
	mul(t[2], t[0], &p2[0])
	mul(t[3], t[1], &p1[0])
	mul(t[0], t[0], &p1[2])
	mul(t[1], t[1], &p2[2])
	mul(t[1], t[1], &p1[1])
	mul(t[0], t[0], &p2[1])
	return equal(t[0], t[1]) && equal(t[2], t[3])
}

func (g *G1) IsOnCurve(p *PointG1) bool {
	if g.IsZero(p) {
		return true
	}
	t := g.t
	square(t[0], &p[1])
	square(t[1], &p[0])
	mul(t[1], t[1], &p[0])
	square(t[2], &p[2])
	square(t[3], t[2])
	mul(t[2], t[2], t[3])
	mul(t[2], b, t[2])
	add(t[1], t[1], t[2])
	return equal(t[0], t[1])
}

func (g *G1) IsAffine(p *PointG1) bool {
	return equal(&p[2], one())
}

func (g *G1) Affine(p *PointG1) *PointG1 {
	if g.IsZero(p) {
		return p
	}
	if !g.IsAffine(p) {
		t := g.t
		inverse(t[0], &p[2])
		square(t[1], t[0])
		mul(&p[0], &p[0], t[1])
		mul(t[0], t[0], t[1])
		mul(&p[1], &p[1], t[0])
		p[2].Set(one())
	}
	return p
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
	square(t[7], &p1[2])
	mul(t[1], &p2[0], t[7])
	mul(t[2], &p1[2], t[7])
	mul(t[0], &p2[1], t[2])
	square(t[8], &p2[2])
	mul(t[3], &p1[0], t[8])
	mul(t[4], &p2[2], t[8])
	mul(t[2], &p1[1], t[4])
	if equal(t[1], t[3]) {
		if equal(t[0], t[2]) {
			return g.Double(r, p1)
		} else {
			return g.Copy(r, infinity)
		}
	}
	sub(t[1], t[1], t[3])
	double(t[4], t[1])
	square(t[4], t[4])
	mul(t[5], t[1], t[4])
	sub(t[0], t[0], t[2])
	double(t[0], t[0])
	square(t[6], t[0])
	sub(t[6], t[6], t[5])
	mul(t[3], t[3], t[4])
	double(t[4], t[3])
	sub(&r[0], t[6], t[4])
	sub(t[4], t[3], &r[0])
	mul(t[6], t[2], t[5])
	double(t[6], t[6])
	mul(t[0], t[0], t[4])
	sub(&r[1], t[0], t[6])
	add(t[0], &p1[2], &p2[2])
	square(t[0], t[0])
	sub(t[0], t[0], t[7])
	sub(t[0], t[0], t[8])
	mul(&r[2], t[0], t[1])
	return r
}

func (g *G1) Double(r, p *PointG1) *PointG1 {
	if g.IsZero(p) {
		g.Copy(r, p)
		return r
	}
	t := g.t
	square(t[0], &p[0])
	square(t[1], &p[1])
	square(t[2], t[1])
	add(t[1], &p[0], t[1])
	square(t[1], t[1])
	sub(t[1], t[1], t[0])
	sub(t[1], t[1], t[2])
	double(t[1], t[1])
	double(t[3], t[0])
	add(t[0], t[3], t[0])
	square(t[4], t[0])
	double(t[3], t[1])
	sub(&r[0], t[4], t[3])
	sub(t[1], t[1], &r[0])
	double(t[2], t[2])
	double(t[2], t[2])
	double(t[2], t[2])
	mul(t[0], t[0], t[1])
	sub(t[1], t[0], t[2])
	mul(t[0], &p[1], &p[2])
	r[1].Set(t[1])
	double(&r[2], t[0])
	return r
}

func (g *G1) Neg(r, p *PointG1) *PointG1 {
	r[0].Set(&p[0])
	r[2].Set(&p[2])
	neg(&r[1], &p[1])
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

func (g *G1) MultiExp(r *PointG1, points []*PointG1, powers []*big.Int) (*PointG1, error) {
	if len(points) != len(powers) {
		return nil, fmt.Errorf("point and scalar vectors should be in same length")
	}
	var c uint32 = 3
	if len(powers) >= 32 {
		c = uint32(math.Ceil(math.Log10(float64(len(powers)))))
	}
	bucketSize, numBits := (1<<c)-1, uint32(g.Q().BitLen())
	windows := make([]*PointG1, numBits/c+1)
	bucket := make([]*PointG1, bucketSize)
	acc, sum := g.New(), g.New()
	for i := 0; i < bucketSize; i++ {
		bucket[i] = g.New()
	}
	mask := (uint64(1) << c) - 1
	j := 0
	var cur uint32
	for cur <= numBits {
		g.Copy(acc, g.Zero())
		bucket = make([]*PointG1, (1<<c)-1)
		for i := 0; i < len(bucket); i++ {
			bucket[i] = g.New()
		}
		for i := 0; i < len(powers); i++ {
			s0 := powers[i].Uint64()
			index := uint(s0 & mask)
			if index != 0 {
				g.Add(bucket[index-1], bucket[index-1], points[i])
			}
			powers[i] = new(big.Int).Rsh(powers[i], uint(c))
		}

		g.Copy(sum, g.Zero())
		for i := len(bucket) - 1; i >= 0; i-- {
			g.Add(sum, sum, bucket[i])
			g.Add(acc, acc, sum)
		}
		windows[j] = g.New()
		g.Copy(windows[j], acc)
		j++
		cur += c
	}
	g.Copy(acc, g.Zero())
	for i := len(windows) - 1; i >= 0; i-- {
		for j := uint32(0); j < c; j++ {
			g.Double(acc, acc)
		}
		g.Add(acc, acc, windows[i])
	}
	g.Copy(r, acc)
	return r, nil
}
