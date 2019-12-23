package bls

import (
	"crypto/cipher"
	"encoding/hex"
	"io"

	"go.dedis.ch/kyber/v3"
	"go.dedis.ch/kyber/v3/group/mod"
)

// KyberG2 is a kyber.Point holding a G2 point on BLS12-381 curve
type KyberG2 struct {
	p *PointG2
}

func nullKyberG2() *KyberG2 {
	var p PointG2
	return newKyberG2(&p)
}
func newKyberG2(p *PointG2) *KyberG2 {
	return &KyberG2{p: p}
}

func (k *KyberG2) Equal(k2 kyber.Point) bool {
	return NewG2(nil).Equal(k.p, k2.(*KyberG2).p)
}

func (k *KyberG2) Null() kyber.Point {
	return newKyberG2(NewG2(nil).Zero())
}

func (k *KyberG2) Base() kyber.Point {
	return newKyberG2(NewG2(nil).One())
}

func (k *KyberG2) Pick(rand cipher.Stream) kyber.Point {
	panic("bls12-381: unsupported operation")
}

func (k *KyberG2) Set(q kyber.Point) kyber.Point {
	k.p.Set(q.(*KyberG2).p)
	return k
}

func (k *KyberG2) Clone() kyber.Point {
	var p PointG2
	p.Set(k.p)
	return newKyberG2(&p)
}

func (k *KyberG2) EmbedLen() int {
	panic("bls12-381: unsupported operation")
}

func (k *KyberG2) Embed(data []byte, rand cipher.Stream) kyber.Point {
	panic("bls12-381: unsupported operation")
}

func (k *KyberG2) Data() ([]byte, error) {
	panic("bls12-381: unsupported operation")
}

func (k *KyberG2) Add(a, b kyber.Point) kyber.Point {
	aa := a.(*KyberG2)
	bb := b.(*KyberG2)
	NewG2(nil).Add(k.p, aa.p, bb.p)
	return k
}

func (k *KyberG2) Sub(a, b kyber.Point) kyber.Point {
	aa := a.(*KyberG2)
	bb := b.(*KyberG2)
	NewG2(nil).Sub(k.p, aa.p, bb.p)
	return k
}

func (k *KyberG2) Neg(a kyber.Point) kyber.Point {
	aa := a.(*KyberG2)
	NewG2(nil).Neg(k.p, aa.p)
	return k
}

func (k *KyberG2) Mul(s kyber.Scalar, q kyber.Point) kyber.Point {
	if q == nil {
		q = nullKyberG2().Base()
	}
	NewG2(nil).MulScalar(k.p, q.(*KyberG2).p, &s.(*mod.Int).V)
	return k
}

func (k *KyberG2) MarshalBinary() ([]byte, error) {
	return NewG2(nil).ToCompressed(k.p), nil
}

func (k *KyberG2) UnmarshalBinary(buff []byte) error {
	var err error
	k.p, err = NewG2(nil).FromCompressed(buff)
	return err
}

func (k *KyberG2) MarshalTo(w io.Writer) (int, error) {
	buf, err := k.MarshalBinary()
	if err != nil {
		return 0, err
	}
	return w.Write(buf)
}

func (k *KyberG2) UnmarshalFrom(r io.Reader) (int, error) {
	buf := make([]byte, k.MarshalSize())
	n, err := io.ReadFull(r, buf)
	if err != nil {
		return n, err
	}
	return n, k.UnmarshalBinary(buf)
}

func (k *KyberG2) MarshalSize() int {
	return 48
}

func (k *KyberG2) String() string {
	b, _ := k.MarshalBinary()
	return "bls12-381.G1: " + hex.EncodeToString(b)
}

func (k *KyberG2) hash(m []byte) kyber.Point {
	panic("bls12-381: Hash() not implemented for G1")
}
