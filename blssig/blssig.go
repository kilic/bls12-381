package blssig

import (
	"crypto/sha256"
	"math/big"

	bls "github.com/kilic/bls12-381"
)

type PointG1 = bls.PointG1
type PublicKey = bls.PointG1
type PointG2 = bls.PointG2
type Signature = bls.PointG2

func Sign(msg [32]byte, domain [8]byte, privateKey *big.Int) *PointG2 {
	g2 := newG2()
	signature := hashToPointWithDomain(g2, msg, domain)
	g2.MulScalar(signature, signature, privateKey)
	return signature
}

func Verify(msg [32]byte, domain [8]byte, signature *PointG2, publicKey *PointG1) bool {
	e := bls.NewBLSPairingEngine()
	return pairingCheck(e, hashToPointWithDomain(e.G2, msg, domain), signature, publicKey)
}

func pairingCheck(e *bls.BLSPairingEngine, msgHash, signature *PointG2, publicKey *PointG1) bool {
	target := &bls.Fe12{}
	e.Pair(target,
		[]bls.PointG1{
			bls.G1NegativeOne,
			*publicKey,
		},
		[]bls.PointG2{
			*signature,
			*msgHash,
		},
	)
	return e.Fp12.Equal(&bls.Fp12One, target)
}

func VerifyAggregateCommon(msg [32]byte, domain [8]byte, publicKeys []*PointG1, signature *PointG2) bool {
	if len(publicKeys) == 0 {
		return false
	}
	e := bls.NewBLSPairingEngine()
	msgHash := hashToPointWithDomain(e.G2, msg, domain)
	aggregated := &bls.PointG1{}
	e.G1.Copy(aggregated, publicKeys[0])
	for i := 1; i < len(publicKeys); i++ {
		e.G1.Add(aggregated, aggregated, publicKeys[i])
	}
	return pairingCheck(e, msgHash, signature, aggregated)
}

func VerifyAggregate(msg [][32]byte, domain [8]byte, publicKeys []*PointG1, signature *PointG2) bool {
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
	e.G1.Copy(&points[0], &bls.G1NegativeOne)
	e.G2.Copy(&twistPoints[0], signature)
	for i := 0; i < size; i++ {
		e.G1.Copy(&points[i+1], publicKeys[i])
		e.G2.Copy(&twistPoints[i+1], hashToPointWithDomain(e.G2, msg[i], domain))
	}
	target := &bls.Fe12{}
	e.Pair(target, points, twistPoints)
	return e.Fp12.Equal(&bls.Fp12One, target)
}

func AggregatePublicKey(p1 *PointG1, p2 *PointG1) *PointG1 {
	return newG1().Add(p1, p1, p2)
}

func NewG1FromCompressed(in []byte) (*PointG1, error) {
	return newG1().FromCompressed(in)
}

func NewG1FromUncompressed(in []byte) (*PointG1, error) {
	return newG1().FromUncompressed(in)
}

func NewG2FromCompressed(in []byte) (*PointG2, error) {
	return newG2().FromCompressed(in)
}

func NewG2FromUncompressed(in []byte) (*PointG2, error) {
	return newG2().FromUncompressed(in)
}

func PublicKeyFromSecretKey(secret *big.Int) *PointG1 {
	p := &PointG1{}
	return newG1().MulScalar(p, &bls.G1One, secret)
}

func HashToPointWithDomain(msg [32]byte, domain [8]byte) *PointG2 {
	return hashToPointWithDomain(newG2(), msg, domain)
}

func hashToPointWithDomain(g2 *bls.G2, msg [32]byte, domain [8]byte) *PointG2 {
	xReBytes := [41]byte{}
	xImBytes := [41]byte{}
	xBytes := make([]byte, 96)
	copy(xReBytes[:32], msg[:])
	copy(xReBytes[32:40], domain[:])
	xImBytes[40] = 0x01
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
