package bls12381

import (
	"fmt"
	"hash"
)

func hashToField(hasher hash.Hash, msg []byte, domain []byte, count int) ([]*fe, error) {
	randBytes, err := expandMsg(hasher, msg, domain, count*64)
	if err != nil {
		return nil, err
	}
	els := make([]*fe, count)
	// var err error
	for i := 0; i < count; i++ {
		els[i], err = from64Bytes(randBytes[i*64 : (i+1)*64])
		if err != nil {
			return nil, err
		}
	}
	return els, nil
}

func expandMsg(h hash.Hash, msg []byte, domain []byte, outLen int) ([]byte, error) {
	domainLen := uint8(len(domain))
	if domainLen > 255 {
		return nil, fmt.Errorf("invalid domain length")
	}
	// var err error
	// DST_prime = DST || I2OSP(len(DST), 1)
	// b_0 = H(Z_pad || msg || l_i_b_str || I2OSP(0, 1) || DST_prime)
	h.Reset()
	h.Write(make([]byte, h.BlockSize()))
	h.Write(msg)
	h.Write([]byte{uint8(outLen >> 8), uint8(outLen)})
	h.Write([]byte{0})
	h.Write(domain)
	h.Write([]byte{domainLen})
	b0 := h.Sum(nil)

	// b_1 = H(b_0 || I2OSP(1, 1) || DST_prime)
	h.Reset()
	h.Write(b0)
	h.Write([]byte{1})
	h.Write(domain)
	h.Write([]byte{domainLen})
	b1 := h.Sum(nil)

	// b_i = H(strxor(b_0, b_(i - 1)) || I2OSP(i, 1) || DST_prime)
	ell := (outLen + h.Size() - 1) / h.Size()
	bi := b1
	out := make([]byte, outLen)
	for i := 1; i < ell; i++ {
		h.Reset()
		// b_i = H(strxor(b_0, b_(i - 1)) || I2OSP(i, 1) || DST_prime)
		tmp := make([]byte, h.Size())
		for j := 0; j < h.Size(); j++ {
			tmp[j] = b0[j] ^ bi[j]
		}
		h.Write(tmp)
		h.Write([]byte{1 + uint8(i)})
		h.Write(domain)
		h.Write([]byte{domainLen})

		// b_1 || ... || b_(ell - 1)
		copy(out[(i-1)*h.Size():i*h.Size()], bi[:])
		bi = h.Sum(nil)
	}
	// b_ell
	copy(out[(ell-1)*h.Size():], bi[:])
	return out[:outLen], nil
}
