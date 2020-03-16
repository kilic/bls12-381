package bls

import (
	"fmt"
	"math/big"
)

type E = fe12

type Gt struct {
	fp12 *fp12
}

func NewGt() *Gt {
	fp12 := newFp12(nil)
	return &Gt{fp12}
}

func (g *Gt) Q() *big.Int {
	return new(big.Int).Set(q)
}

func (g *Gt) FromBytes(in []byte) (*E, error) {
	e, err := g.fp12.fromBytes(in)
	if err != nil {
		return nil, err
	}
	if !g.IsValid(e) {
		return nil, fmt.Errorf("invalid element")
	}
	return e, nil
}

func (g *Gt) FromBytesUnchecked(in []byte) (*E, error) {
	e, err := g.fp12.fromBytes(in)
	if err != nil {
		return nil, err
	}
	return e, nil
}

func (g *Gt) ToBytes(e *E) []byte {
	return g.fp12.toBytes(e)
}

func (g *Gt) IsValid(e *E) bool {
	r := g.New()
	g.fp12.exp(r, e, q)
	return g.Equal(r, g.fp12.one())
}

func (g *Gt) New() *E {
	return g.One()
}

func (g *Gt) One() *E {
	return g.fp12.one()
}

func (g *Gt) IsOne(e *E) bool {
	return g.Equal(e, g.fp12.one())
}

func (g *Gt) Copy(a, b *E) {
	g.fp12.copy(a, b)
}

func (g *Gt) Equal(a, b *E) bool {
	return g.fp12.equal(a, b)
}

func (g *Gt) Add(c, a, b *E) {
	g.fp12.add(c, a, b)
}

func (g *Gt) Sub(c, a, b *E) {
	g.fp12.sub(c, a, b)
}

func (g *Gt) Mul(c, a, b *E) {
	g.fp12.mul(c, a, b)
}

func (g *Gt) Square(c, a *E) {
	g.fp12.cyclotomicSquare(c, a)
}

func (g *Gt) Exp(c, a *E, s *big.Int) {
	g.fp12.cyclotomicExp(c, a, s)
}
