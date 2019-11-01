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
	iter := flag.Int("iter", 10, "# of iterations")
	bmi2 := flag.Bool("nobmi2", false, "to enfoce non bmi2 arch")
	flag.Parse()
	n = *iter
	enforceNonBMI2 = *bmi2
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
