package bls

import (
	"crypto/rand"
	"errors"
	"flag"
	"math/big"
	"testing"
)

var n int

func TestMain(m *testing.M) {
	iter := flag.Int("iter", 10, "# of iterationss")
	flag.Parse()
	n = *iter
	m.Run()
}

func randScalar(max *big.Int) *big.Int {
	a, err := rand.Int(rand.Reader, max)
	if err != nil {
		panic(errors.New(""))
	}
	return a
}

func randScalars(max *big.Int, size int) []*big.Int {
	var scalars []*big.Int
	for i := 0; i < size; i++ {
		a, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(errors.New(""))
		}
		scalars = append(scalars, a)
	}
	return scalars
}
