package bls

import (
	"crypto/rand"
	"errors"
	"flag"
	"math/big"
	"testing"
)

var fuz int
var forceNonADXArch bool

func TestMain(m *testing.M) {
	_fuz := flag.Int("fuzz", 10, "# of iterations")
	adx := flag.Bool("noadx", false, "to enfoce non adx arch")
	flag.Parse()
	forceNonADXArch = *adx
	fuz = *_fuz
	setup()
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
