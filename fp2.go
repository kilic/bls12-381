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
	fp *fp
	fp2Temp
}

func newFp2Temp() fp2Temp {
	t := [4]*fe{}
	for i := 0; i < len(t); i++ {
		t[i] = &fe{}
	}
	return fp2Temp{t}
}

func newFp2(fp *fp) *fp2 {
	t := newFp2Temp()
	if fp == nil {
		return &fp2{newFp(), t}
	}
	return &fp2{fp, t}
}

func (e *fp2) fromBytes(in []byte) (*fe2, error) {
	if len(in) != 96 {
		return nil, fmt.Errorf("input string should be larger than 96 bytes")
	}
	fp := e.fp
	c0, err := fp.fromBytes(in[:48])
	if err != nil {
		return nil, err
	}
	c1, err := fp.fromBytes(in[48:])
	if err != nil {
		return nil, err
	}
	return &fe2{*c0, *c1}, nil
}

func (e *fp2) toBytes(a *fe2) []byte {
	fp := e.fp
	out := make([]byte, 96)
	copy(out[:48], fp.toBytes(&a[1]))
	copy(out[48:], fp.toBytes(&a[0]))
	return out
}

func (e *fp2) new() *fe2 {
	return e.zero()
}

func (e *fp2) zero() *fe2 {
	return &fe2{}
}

func (e *fp2) one() *fe2 {
	fp := e.fp
	return &fe2{*fp.one(), *fp.zero()}
}

func (e *fp2) rand(r io.Reader) (*fe2, error) {
	fp := e.fp
	a0, err := fp.rand(r)
	if err != nil {
		return nil, err
	}
	a1, err := fp.rand(r)
	if err != nil {
		return nil, err
	}
	return &fe2{*a0, *a1}, nil
}

func (e *fp2) fromMont(c, a *fe2) {
	fp := e.fp
	fp.fromMont(&c[0], &a[0])
	fp.fromMont(&c[1], &a[1])
}

func (e *fp2) isZero(a *fe2) bool {
	fp := e.fp
	return fp.isZero(&a[0]) && fp.isZero(&a[1])
}

func (e *fp2) isOne(a *fe2) bool {
	fp := e.fp
	return fp.isOne(&a[0]) && fp.isZero(&a[1])
}

func (e *fp2) equal(a, b *fe2) bool {
	fp := e.fp
	return fp.equal(&a[0], &b[0]) && fp.equal(&a[1], &b[1])
}

func (e *fp2) copy(c, a *fe2) {
	fp := e.fp
	fp.copy(&c[0], &a[0])
	fp.copy(&c[1], &a[1])
}

func (e *fp2) add(c, a, b *fe2) {
	fp := e.fp
	fp.add(&c[0], &a[0], &b[0])
	fp.add(&c[1], &a[1], &b[1])
}

func (e *fp2) addAssign(a, b *fe2) {
	fp := e.fp
	fp.addAssign(&a[0], &b[0])
	fp.addAssign(&a[1], &b[1])
}

func (e *fp2) ladd(c, a, b *fe2) {
	fp := e.fp
	fp.ladd(&c[0], &a[0], &b[0])
	fp.ladd(&c[1], &a[1], &b[1])
}

func (e *fp2) double(c, a *fe2) {
	fp := e.fp
	fp.double(&c[0], &a[0])
	fp.double(&c[1], &a[1])
}

func (e *fp2) doubleAssign(a *fe2) {
	fp := e.fp
	fp.doubleAssign(&a[0])
	fp.doubleAssign(&a[1])
}

// ldouble doubles field element `a` and sets the result `c` without modular reduction
func (e *fp2) ldouble(c, a *fe2) {
	fp := e.fp
	fp.ldouble(&c[0], &a[0])
	fp.ldouble(&c[1], &a[1])
}

func (e *fp2) sub(c, a, b *fe2) {
	fp := e.fp
	fp.sub(&c[0], &a[0], &b[0])
	fp.sub(&c[1], &a[1], &b[1])
}

func (e *fp2) subAssign(c, a *fe2) {
	fp := e.fp
	fp.subAssign(&c[0], &a[0])
	fp.subAssign(&c[1], &a[1])
}

// lsub subtracts field element `b` from `a` and sets the result `c` without modular reduction
func (e *fp2) lsub(c, a, b *fe2) {
	fp := e.fp
	fp.lsub(&c[0], &a[0], &b[0])
	fp.lsub(&c[1], &a[1], &b[1])
}

func (e *fp2) neg(c, a *fe2) {
	fp := e.fp
	fp.neg(&c[0], &a[0])
	fp.neg(&c[1], &a[1])
}

func (e *fp2) conjugate(c, a *fe2) {
	fp := e.fp
	fp.copy(&c[0], &a[0])
	fp.neg(&c[1], &a[1])
}

func (e *fp2) mul(c, a, b *fe2) {
	fp, t := e.fp, e.t
	fp.mul(t[1], &a[0], &b[0])
	fp.mul(t[2], &a[1], &b[1])
	fp.add(t[0], &a[0], &a[1])
	fp.add(t[3], &b[0], &b[1])
	fp.sub(&c[0], t[1], t[2])
	fp.addAssign(t[1], t[2])
	fp.mulAssign(t[0], t[3])
	fp.sub(&c[1], t[0], t[1])
}

func (e *fp2) mulAssign(a, b *fe2) {
	fp, t := e.fp, e.t
	fp.mul(t[1], &a[0], &b[0])
	fp.mul(t[2], &a[1], &b[1])
	fp.add(t[0], &a[0], &a[1])
	fp.add(t[3], &b[0], &b[1])
	fp.sub(&a[0], t[1], t[2])
	fp.addAssign(t[1], t[2])
	fp.mulAssign(t[0], t[3])
	fp.sub(&a[1], t[0], t[1])
}

func (e *fp2) square(c, a *fe2) {
	t, fp := e.t, e.fp
	fp.ladd(t[0], &a[0], &a[1])
	fp.sub(t[1], &a[0], &a[1])
	fp.ldouble(t[2], &a[0])
	fp.mul(&c[0], t[0], t[1])
	fp.mul(&c[1], t[2], &a[1])
}

func (e *fp2) squareAssign(a *fe2) {
	t, fp := e.t, e.fp
	fp.ladd(t[0], &a[0], &a[1])
	fp.sub(t[1], &a[0], &a[1])
	fp.ldouble(t[2], &a[0])
	fp.mul(&a[0], t[0], t[1])
	fp.mul(&a[1], t[2], &a[1])
}

func (e *fp2) mulByNonResidue(c, a *fe2) {
	t, fp := e.t, e.fp
	fp.sub(t[0], &a[0], &a[1])
	fp.add(&c[1], &a[0], &a[1])
	fp.copy(&c[0], t[0])
}

func (e *fp2) mulByB(c, a *fe2) {
	t, fp := e.t, e.fp
	fp.double(t[0], &a[0])
	fp.double(t[1], &a[1])
	fp.doubleAssign(t[0])
	fp.doubleAssign(t[1])
	fp.sub(&c[0], t[0], t[1])
	fp.add(&c[1], t[0], t[1])
}

func (e *fp2) inverse(c, a *fe2) {
	t, fp := e.t, e.fp
	fp.square(t[0], &a[0])
	fp.square(t[1], &a[1])
	fp.addAssign(t[0], t[1])
	fp.inverse(t[0], t[0])
	fp.mul(&c[0], &a[0], t[0])
	fp.mulAssign(t[0], &a[1])
	fp.neg(&c[1], t[0])
}

func (e *fp2) mulByFq(c, a *fe2, b *fe) {
	fp := e.fp
	fp.mul(&c[0], &a[0], b)
	fp.mul(&c[1], &a[1], b)
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
	fp := e.fp
	fp.copy(&c[0], &a[0])
	if power%2 == 1 {
		fp.neg(&c[1], &a[1])
		return
	}
	fp.copy(&c[1], &a[1])
}

func (e *fp2) frobeniousMapAssign(a *fe2, power uint) {
	fp := e.fp
	if power%2 == 1 {
		fp.neg(&a[1], &a[1])
		return
	}
}

func (e *fp2) sqrt(c, a *fe2) bool {
	fp := e.fp
	u, x0, a1, alpha := &fe2{}, &fe2{}, &fe2{}, &fe2{}
	e.copy(u, a)
	e.exp(a1, a, pMinus3Over4)
	e.square(alpha, a1)
	e.mul(alpha, alpha, a)
	e.mul(x0, a1, a)
	if e.equal(alpha, negativeOne2) {
		fp.neg(&c[0], &x0[1])
		fp.copy(&c[1], &x0[0])
		return true
	}
	e.add(alpha, alpha, e.one())
	e.exp(alpha, alpha, pMinus1Over2)
	e.mul(c, alpha, x0)
	e.square(alpha, c)
	return e.equal(alpha, u)
}
