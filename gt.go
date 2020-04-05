package bls12381

import (
	"fmt"
	"math/big"
)

// E is type for target group element
type E = fe12

// GT is type for target multiplicative group GT.
type GT struct {
	fp12 *fp12
}

// NewGT constructs new target group instance.
func NewGT() *GT {
	fp12 := newFp12(nil)
	return &GT{fp12}
}

// Q returns group order in big.Int.
func (g *GT) Q() *big.Int {
	return new(big.Int).Set(q)
}

// FromBytes expects 576 byte input and returns target group element
// FromBytes returns error if given element is not on correct subgroup.
func (g *GT) FromBytes(in []byte) (*E, error) {
	e, err := g.fp12.fromBytes(in)
	if err != nil {
		return nil, err
	}
	if !g.IsValid(e) {
		return e, fmt.Errorf("invalid element")
	}
	return e, nil
}

// ToBytes serializes target group element.
func (g *GT) ToBytes(e *E) []byte {
	return g.fp12.toBytes(e)
}

// IsValid checks whether given target group element is in correct subgroup.
func (g *GT) IsValid(e *E) bool {
	r := g.New()
	g.fp12.exp(r, e, q)
	return g.Equal(r, g.fp12.one())
}

// New initializes a new target group element which is equal to one
func (g *GT) New() *E {
	return g.One()
}

// One initializes a new target group element which is equal to one
func (g *GT) One() *E {
	return g.fp12.one()
}

// IsOne returns true if given element equals to one
func (g *GT) IsOne(e *E) bool {
	return g.Equal(e, g.fp12.one())
}

// Copy copies values of the second source element to first element
func (g *GT) Copy(a, b *E) {
	g.fp12.copy(a, b)
}

// Equal returns true if given two element is equal, otherwise returns false
func (g *GT) Equal(a, b *E) bool {
	return g.fp12.equal(a, b)
}

// Add adds two field element `a` and `b` and assigns the result to the element in first argument.
func (g *GT) Add(c, a, b *E) {
	g.fp12.add(c, a, b)
}

// Sub subtracts two field element `a` and `b`, and assigns the result to the element in first argument.
func (g *GT) Sub(c, a, b *E) {
	g.fp12.sub(c, a, b)
}

// Mul multiplies two field element `a` and `b` and assigns the result to the element in first argument.
func (g *GT) Mul(c, a, b *E) {
	g.fp12.mul(c, a, b)
}

// Square squares an element `a` and assigns the result to the element in first argument.
func (g *GT) Square(c, a *E) {
	g.fp12.cyclotomicSquare(c, a)
}

// Exp exponents an element `a` by a scalar `s` and assigns the result to the element in first argument.
func (g *GT) Exp(c, a *E, s *big.Int) {
	g.fp12.cyclotomicExp(c, a, s)
}

// Inverse inverses an element `a` and assigns the result to the element in first argument.
func (g *GT) Inverse(c, a *E) {
	g.fp12.inverse(c, a)
}
