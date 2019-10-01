package bls

import (
	"crypto/sha256"
)

// HashToG2WithDomain hashes 32 byte message and uint64 domain parameters to a G2 point.
// This method is pending the BLS standardisation process.
func HashToG2WithDomain(messageHash [32]byte, domain [8]byte) *PointG2 {
	hashedObj := hashWithDomain(messageHash, domain)
	g2Elems := NewG2(nil)
	return g2Elems.MapToPoint(hashedObj)
}

func hashWithDomain(messageHash [32]byte, domain [8]byte) []byte {
	xReBytes := [41]byte{}
	xImBytes := [41]byte{}
	xBytes := make([]byte, 96)
	copy(xReBytes[:32], messageHash[:])
	copy(xReBytes[32:40], domain[:])
	copy(xReBytes[40:41], []byte{0x01})
	copy(xImBytes[:32], messageHash[:])
	copy(xImBytes[32:40], domain[:])
	copy(xImBytes[40:41], []byte{0x02})
	copy(xBytes[16:48], sha256Hash(xImBytes[:]))
	copy(xBytes[64:], sha256Hash(xReBytes[:]))
	return xBytes
}

// sha256Hash using the sha256 hashing algorithm. It takes the
// given input and hashes it using sha256.
func sha256Hash(in []byte) []byte {
	h := sha256.New()
	h.Write(in)
	return h.Sum(nil)
}
