package blssig

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"io"
	"math/big"

	bls "github.com/kilic/bls12-381"
)

type PointG1 = bls.PointG1
type PointG2 = bls.PointG2
type PublicKey = bls.PointG1
type Signature = bls.PointG2
type SecretKey = *big.Int

var curveOrder, _ = new(big.Int).SetString("73eda753299d7d483339d80809a1d80553bda402fffe5bfeffffffff00000001", 16)

func Sign(msg [32]byte, domain [8]byte, privateKey *big.Int) *Signature {
	g2 := newG2()
	signature := hashWithDomain(g2, msg, domain)
	g2.MulScalar(signature, signature, privateKey)
	return signature
}

func Verify(msg [32]byte, domain [8]byte, signature *Signature, publicKey *PublicKey) bool {
	e := bls.NewBLSPairingEngine()
	return pairingCheck(e, hashWithDomain(e.G2, msg, domain), signature, publicKey)
}

func pairingCheck(e *bls.BLSPairingEngine, msgHash, signature *Signature, publicKey *PublicKey) bool {
	return e.PairingCheck(
		[]bls.PointG1{
			*e.G1.NegativeOne(),
			*publicKey,
		},
		[]bls.PointG2{
			*signature,
			*msgHash,
		},
	)
}

func VerifyAggregateCommon(msg [32]byte, domain [8]byte, publicKeys []*PublicKey, signature *Signature) bool {
	if len(publicKeys) == 0 {
		return false
	}
	e := bls.NewBLSPairingEngine()
	msgHash := hashWithDomain(e.G2, msg, domain)
	aggregated := &bls.PointG1{}
	e.G1.Copy(aggregated, publicKeys[0])
	for i := 1; i < len(publicKeys); i++ {
		e.G1.Add(aggregated, aggregated, publicKeys[i])
	}
	return pairingCheck(e, msgHash, signature, aggregated)
}

func VerifyAggregate(msg [][32]byte, domain [8]byte, publicKeys []*PublicKey, signature *Signature) bool {
	size := len(publicKeys)
	if size == 0 {
		return false
	}
	if size != len(msg) {
		return false
	}
	points := make([]bls.PointG1, size+1)
	twistPoints := make([]bls.PointG2, size+1)
	e := bls.NewBLSPairingEngine()
	e.G1.Copy(&points[0], e.G1.NegativeOne())
	e.G2.Copy(&twistPoints[0], signature)
	for i := 0; i < size; i++ {
		e.G1.Copy(&points[i+1], publicKeys[i])
		e.G2.Copy(&twistPoints[i+1], hashWithDomain(e.G2, msg[i], domain))
	}
	return e.PairingCheck(points, twistPoints)
}

func AggregatePublicKey(p1 *PublicKey, p2 *PublicKey) *PublicKey {
	return newG1().Add(p1, p1, p2)
}

func AggregatePublicKeys(keys []*PublicKey) *PointG1 {
	g := newG1()
	if len(keys) == 0 {
		return g.Zero()
	}
	aggregated := new(PublicKey).Set(keys[0])
	for _, p := range keys[1:] {
		g.Add(aggregated, aggregated, p)
	}
	return aggregated
}

func AggregateSignature(p1 *Signature, p2 *Signature) *Signature {
	return newG2().Add(p1, p1, p2)
}

func AggregateSignatures(keys []*Signature) *Signature {
	g := newG2()
	if len(keys) == 0 {
		return g.Zero()
	}
	aggregated := new(Signature).Set(keys[0])
	for _, p := range keys[1:] {
		g.Add(aggregated, aggregated, p)
	}
	return aggregated
}

func SecretKeyFromBytes(priv []byte) (SecretKey, error) {
	var y [32]byte
	copy(y[:], priv)
	k := new(big.Int).SetBytes(y[:])
	if curveOrder.Cmp(k) < 0 {
		return nil, errors.New("invalid private key")
	}
	return k, nil
}

func RandSecretKey(r io.Reader) (SecretKey, error) {
	k, err := rand.Int(r, curveOrder)
	if err != nil {
		return nil, err
	}
	return k, nil
}

func NewPublicKeyFromCompresssed(in []byte) (*PointG1, error) {
	return newG1FromCompressed(in)
}

func NewPublicKeyFromUncompresssed(in []byte) (*PointG1, error) {
	return newG1FromCompressed(in)
}

func PublicKeyToCompressed(p *PointG1) []byte {
	return g1ToCompressed(p)
}

func PublicKeyToUncompressed(p *PointG1) []byte {
	return g1ToCompressed(p)
}

func NewSignatureFromCompresssed(in []byte) (*PointG2, error) {
	return newG2FromCompressed(in)
}

func NewSignatureFromUncompresssed(in []byte) (*PointG2, error) {
	return newG2FromUncompressed(in)
}

func SignatureToCompressed(p *Signature) []byte {
	return g2ToCompressed(p)
}

func SignatureToUncompressed(p *Signature) []byte {
	return g2ToUncompressed(p)
}

func newG1FromCompressed(in []byte) (*PointG1, error) {
	return newG1().FromCompressed(in)
}

func newG1FromUncompressed(in []byte) (*PointG1, error) {
	return newG1().FromUncompressed(in)
}

func g1ToCompressed(p *PointG1) []byte {
	return newG1().ToCompressed(p)
}

func g1ToUncompressed(p *PointG1) []byte {
	return newG1().ToUncompressed(p)
}

func newG2FromCompressed(in []byte) (*PointG2, error) {
	return newG2().FromCompressed(in)
}

func newG2FromUncompressed(in []byte) (*PointG2, error) {
	return newG2().FromUncompressed(in)
}

func g2ToCompressed(p *PointG2) []byte {
	return newG2().ToCompressed(p)
}

func g2ToUncompressed(p *PointG2) []byte {
	return newG2().ToUncompressed(p)
}

func PublicKeyFromSecretKey(secret SecretKey) *PublicKey {
	p := &PointG1{}
	return newG1().MulScalar(p, &bls.G1One, secret)
}

func HashWithDomain(msg [32]byte, domain [8]byte) *PointG2 {
	return hashWithDomain(newG2(), msg, domain)
}

func hashWithDomain(g2 *bls.G2, msg [32]byte, domain [8]byte) *PointG2 {
	xReBytes := [41]byte{}
	xImBytes := [41]byte{}
	xBytes := make([]byte, 96)
	copy(xReBytes[:32], msg[:])
	copy(xReBytes[32:40], domain[:])
	xReBytes[40] = 0x01
	copy(xImBytes[:32], msg[:])
	copy(xImBytes[32:40], domain[:])
	xImBytes[40] = 0x02
	copy(xBytes[16:48], sha256Hash(xImBytes[:]))
	copy(xBytes[64:], sha256Hash(xReBytes[:]))
	return g2.MapToPoint(xBytes)
}

func sha256Hash(in []byte) []byte {
	h := sha256.New()
	h.Write(in)
	return h.Sum(nil)
}

func newG1() *bls.G1 {
	return bls.NewG1(nil)
}

func newG2() *bls.G2 {
	return bls.NewG2(nil)
}
