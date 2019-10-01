package bls

import (
	"bytes"
	"encoding/hex"
	"testing"
)

var expectedSerializedG2, _ = hex.DecodeString("a6ef29e7241e1a1cc60fee328e3290c023d55a6701db500eefab7f91391a8b8726fd0024121e64637281f907137fe268187b4baca36388e96194b73a7d532f6eea6bc098778dbfd3404584613b5ba9da97d5602e31fdbe9270b863876529b254")

func TestHashG2WithDomain(t *testing.T) {
	g2Point := HashToG2WithDomain([32]byte{}, [8]byte{})
	g2Elems := NewG2(nil)
	compressedPoint := g2Elems.ToCompressed(g2Point)

	if !bytes.Equal(expectedSerializedG2, compressedPoint[:]) {
		t.Fatal("expected hash to match test")
	}
}
