package bls12381

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestHashToField(t *testing.T) {

	msg := []byte("hello world")
	domain := []byte("asdfqwerzxcv")
	count := 5
	hasher := sha256.New()
	els, err := hashToField(hasher, msg, domain, count)
	if err != nil {
		t.Fatal(err)
	}
	_ = els
	for i := 0; i < count; i++ {
		fmt.Printf("%x\n", toBytes(els[i]))
	}
}
