package bls12381

import (
	"errors"
	"math/big"
)

type nafNumber []int

var bigZero = big.NewInt(0)
var bigOne = big.NewInt(1)

func bigToWNAF(e *big.Int, w int) nafNumber {

	naf := nafNumber{}
	k := new(big.Int).Lsh(bigOne, uint(w))
	halfK := new(big.Int).Rsh(k, 1)

	ee := new(big.Int).Set(e)
	for ee.Cmp(bigZero) != 0 {

		if ee.Bit(0) == 1 {

			nafSign := new(big.Int)
			nafSign.Mod(ee, k)

			if nafSign.Cmp(halfK) >= 0 {
				nafSign.Sub(nafSign, k)
			}

			naf = append(naf, int(nafSign.Int64()))

			ee.Sub(ee, nafSign)

		} else {
			naf = append(naf, 0)
		}

		ee.Rsh(ee, 1)
	}
	return naf
}

func bigFromWNAF(naf nafNumber, w int) *big.Int {
	acc := new(big.Int)
	k := new(big.Int).Set(bigOne)
	for i := 0; i < len(naf); i++ {
		if naf[i] != 0 {
			z := new(big.Int).Mul(k, big.NewInt(int64(naf[i])))
			acc.Add(acc, z)
		}
		k.Lsh(k, 1)
	}
	return acc
}

func (e *Fr) toWNAF(w int) (nafNumber, error) {
	if w < 2 {
		return nil, errors.New("bad window size")
	}

	naf := nafNumber{}

	k, halfK, kMask := 1<<w, 1<<(w-1), (1<<w)-1

	ee := new(Fr).Set(e)

	for !ee.IsZero() {

		if ee.Bit(0) {

			nafSign := int(ee[0]) & kMask
			if nafSign >= halfK {
				nafSign = nafSign - k
			}

			naf = append(naf, int(nafSign))
			if nafSign < 0 {
				ee.Add(ee, &Fr{uint64(-nafSign)})
			} else {
				ee.Sub(ee, &Fr{uint64(nafSign)})
			}

		} else {
			naf = append(naf, 0)
		}
		ee.div2()
	}

	return naf, nil
}

func (e *Fr) fromWNAF(naf nafNumber, w int) *Fr {
	acc := new(Fr).Zero()
	k := new(Fr).One()
	for i := 0; i < len(naf); i++ {

		if naf[i] < 0 {

			z := new(Fr)
			z.Mul(k, &Fr{uint64(-naf[i])})
			acc.Sub(acc, z)

		} else if naf[i] > 0 {

			z := new(Fr)
			z.Mul(k, &Fr{uint64(naf[i])})
			acc.Add(acc, z)
		}

		k.mul2()
	}
	return acc
}
