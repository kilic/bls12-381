package bls

import (
	"crypto/cipher"
	"crypto/sha256"
	"encoding/hex"
	"io"

	"github.com/drand/kyber"
	"github.com/drand/kyber/group/mod"
)

var domainG2 = [8]byte{1, 9, 6, 9, 9, 6, 9, 2}

// BLS_SIG_BLS12381G2-SHA256-SSWU-RO-_NUL_
var csuite = []byte{66, 76, 83, 95, 83, 73, 71, 95, 66, 76, 83, 49, 50, 51, 56, 49, 71, 50, 45, 83, 72, 65, 50, 53, 54, 45, 83, 83, 87, 85, 45, 82, 79, 45, 95, 78, 85, 76, 95}

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
	var dst, src [32]byte
	rand.XORKeyStream(dst[:], src[:])
	return k.Hash(dst[:])
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
	return 96
}

func (k *KyberG2) String() string {
	b, _ := k.MarshalBinary()
	return "bls12-381.G1: " + hex.EncodeToString(b)
}

/*func hashWithDomainG2(g2 *G2, msg [32]byte, domain [8]byte) *PointG2 {*/
//xReBytes := [41]byte{}
//xImBytes := [41]byte{}
//xBytes := make([]byte, 96)
//copy(xReBytes[:32], msg[:])
//copy(xReBytes[32:40], domain[:])
//xReBytes[40] = 0x01
//copy(xImBytes[:32], msg[:])
//copy(xImBytes[32:40], domain[:])
//xImBytes[40] = 0x02
//copy(xBytes[16:48], sha256Hash(xImBytes[:]))
//copy(xBytes[64:], sha256Hash(xReBytes[:]))
//return g2.MapToPoint(xBytes)
/*}*/
func (k *KyberG2) Hash(m []byte) kyber.Point {
	if len(m) != 32 {
		m = sha256Hash(m)
	}
	var s [32]byte
	copy(s[:], m)
	pg2 := hashWithDomainG2(NewG2(nil), s, domainG2)
	//buff := blsmap.MapToG2(m, nil)
	//_ = k.UnmarshalBinary(buff)
	return newKyberG2(pg2)
}

func sha256Hash(in []byte) []byte {
	h := sha256.New()
	h.Write(in)
	return h.Sum(nil)
}
