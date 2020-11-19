package bls12381

import (
	"math/big"
)

// Guide to Pairing Based Cryptography
// 6.3.2. Decompositions for the k = 12 BLS Family

// glvQ1 = x^2 * R / q
var glvQ1 = &Fr{0x63f6e522f6cfee30, 0x7c6becf1e01faadd, 0x1, 0}
var glvQ1Big = bigFromHex("0x017c6becf1e01faadd63f6e522f6cfee30")

// glvQ2 = R / q = 2
var glvQ2 = &Fr{0x02, 0, 0, 0}
var glvQ2Big = bigFromHex("0x02")

// glvB1 = x^2 - 1 = 0xac45a4010001a40200000000ffffffff
var glvB1 = &Fr{0x00000000ffffffff, 0xac45a4010001a402, 0, 0}
var glvB1Big = bigFromHex("0xac45a4010001a40200000000ffffffff")

// glvB2 = x^2 = 0xac45a4010001a4020000000100000000
var glvB2 = &Fr{0x0000000100000000, 0xac45a4010001a402, 0, 0}
var glvB2Big = bigFromHex("0xac45a4010001a4020000000100000000")

// glvLambda1 = x^2 - 1
var glvLambda1 = &Fr{0x00000000ffffffff, 0xac45a4010001a402, 0, 0}
var glvLambda1Big = bigFromHex("0xac45a4010001a40200000000ffffffff")

// var glvLambda1Alt = &Fr{0xfffffffe00000001, 0xa7780001fffcb7fc, 0x3339d80809a1d804, 0x73eda753299d7d48}

// halfR = 2**256 / 2
var halfR = &wideFr{0, 0, 0, 0x8000000000000000, 0, 0, 0}
var halfRBig = bigFromHex("0x8000000000000000000000000000000000000000000000000000000000000000")

// r128 = 2**128 - 1
var r128 = &Fr{0xffffffffffffffff, 0xffffffffffffffff, 0, 0}

// glvPhi = 0x1a0111ea397fe699ec02408663d4de85aa0d857d89759ad4897d29650fb85f9b409427eb4f49fffd8bfd00000000aaac
var glvPhi = &fe{0xcd03c9e48671f071, 0x5dab22461fcda5d2, 0x587042afd3851b95, 0x8eb60ebe01bacb9e, 0x03f97d6e83d050d2, 0x18f0206554638741}

// var glvPhiAlt = &fe{ 0x30f1361b798a64e8, 0xf3b8ddab7ece5a2a, 0x16a8ca3ac61577f7, 0xc26a2ff874fd029b, 0x3636b76660701c6e, 0x051ba4ab241b6160, }

type glvVector struct {
	m1   *Fr
	m2   *Fr
	neg1 bool
	neg2 bool
}

func decompose(m *Fr) *glvVector {
	// Guide to Pairing Based Cryptography
	// 6.3.2. Decompositions for the k = 12 BLS Family

	// alpha1 = round(x^2 * m  / r)
	alpha1 := x2mr(m)
	// alpha2 = round(m / r)
	alpha2 := mr(m)

	z1, z2 := new(Fr), new(Fr)

	// z1 = (x^2 - 1) * round(x^2 * m  / r)
	z1.Mul(alpha1, glvB1)
	// z2 = x^2 * round(m / r)
	z2.Mul(alpha2, glvB2)

	a1, a2 := new(Fr), new(Fr)
	// a1 = m - z1 - alpha2
	a1.Sub(m, z1)
	a1.Sub(a1, alpha2)

	// a2 = z2 - alpha1
	a2.Sub(z2, alpha1)

	v := &glvVector{}
	if a1.Cmp(r128) == 1 {
		a1.Neg(a1)
		v.neg1 = true
	}
	v.m1 = new(Fr).Set(a1)
	if a2.Cmp(r128) == 1 {
		a2.Neg(a2)
		v.neg2 = true
	}
	v.m2 = new(Fr).Set(a2)
	return v
}

func decomposeBig(m *big.Int) (*big.Int, *big.Int) {
	// Guide to Pairing Based Cryptography
	// 6.3.2. Decompositions for the k = 12 BLS Family

	// alpha1 = round(x^2 * m  / r)
	alpha1 := new(big.Int).Mul(m, glvQ1Big)
	alpha1.Add(alpha1, halfRBig)
	alpha1.Rsh(alpha1, fourWordBitSize)

	// alpha2 = round(m / r)
	alpha2 := new(big.Int).Mul(m, glvQ2Big)
	alpha2.Add(alpha2, halfRBig)
	alpha2.Rsh(alpha2, fourWordBitSize)

	z1, z2 := new(big.Int), new(big.Int)
	// z1 = (x^2 - 1) * round(x^2 * m  / r)
	z1.Mul(alpha1, glvB1Big).Mod(z1, qBig)
	// z2 = x^2 * round(m / r)
	z2.Mul(alpha2, glvB2Big).Mod(z2, qBig)

	a1, a2 := new(big.Int), new(big.Int)

	// a1 = m - z1 - alpha2
	a1.Sub(m, z1)
	a1.Sub(a1, alpha2)

	// a2 = z2 - alpha1
	a2.Sub(z2, alpha1)

	return a1, a2
}

func x2mr(m *Fr) *Fr {
	a := new(wideFr)
	// m * x^2 * R / q
	a.mul(m, glvQ1)
	// round(x^2 * m  / q)
	return a.round()
}

func mr(m *Fr) *Fr {
	a := new(wideFr)
	// m * R / q
	a.mul(m, glvQ2)
	// round(m * R / q)
	return a.round()
}
