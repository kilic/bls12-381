package bls

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestPairing(t *testing.T) {
	bls := NewBLSPairingEngine()
	g1Zero, g2Zero, g1One, g2One := bls.G1.Zero(), bls.G2.Zero(), bls.G1.One(), bls.G2.One()
	generatorPair := Pairs{Pair{*g1One, *g2One}}
	expected, err := bls.fp12.fromBytes(
		fromHex(48,
			"0x0f41e58663bf08cf068672cbd01a7ec73baca4d72ca93544deff686bfd6df543d48eaa24afe47e1efde449383b676631",
			"0x04c581234d086a9902249b64728ffd21a189e87935a954051c7cdba7b3872629a4fafc05066245cb9108f0242d0fe3ef",
			"0x03350f55a7aefcd3c31b4fcb6ce5771cc6a0e9786ab5973320c806ad360829107ba810c5a09ffdd9be2291a0c25a99a2",
			"0x11b8b424cd48bf38fcef68083b0b0ec5c81a93b330ee1a677d0d15ff7b984e8978ef48881e32fac91b93b47333e2ba57",
			"0x06fba23eb7c5af0d9f80940ca771b6ffd5857baaf222eb95a7d2809d61bfe02e1bfd1b68ff02f0b8102ae1c2d5d5ab1a",
			"0x19f26337d205fb469cd6bd15c3d5a04dc88784fbb3d0b2dbdea54d43b2b73f2cbb12d58386a8703e0f948226e47ee89d",
			"0x018107154f25a764bd3c79937a45b84546da634b8f6be14a8061e55cceba478b23f7dacaa35c8ca78beae9624045b4b6",
			"0x01b2f522473d171391125ba84dc4007cfbf2f8da752f7c74185203fcca589ac719c34dffbbaad8431dad1c1fb597aaa5",
			"0x193502b86edb8857c273fa075a50512937e0794e1e65a7617c90d8bd66065b1fffe51d7a579973b1315021ec3c19934f",
			"0x1368bb445c7c2d209703f239689ce34c0378a68e72a6b3b216da0e22a5031b54ddff57309396b38c881c4c849ec23e87",
			"0x089a1c5b46e5110b86750ec6a532348868a84045483c92b7af5af689452eafabf1a8943e50439f1d59882a98eaa0170f",
			"0x1250ebd871fc0a92a7b2d83168d0d727272d441befa15c503dd8e90ce98db3e7b6d194f60839c508a84305aaca1789b6",
		))
	if err != nil {
		t.Fatal(err)
	}
	t.Run("Single Expected", func(t *testing.T) {
		pairs := Pairs{Pair{*g1One, *g2One}}
		f1 := bls.pair(pairs)
		if !bls.fp12.equal(f1, expected) {
			t.Fatal("bad pairing")
		}
	})

	randomPair := func() (Pair, *big.Int) {
		var pair Pair
		a, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		bls.G1.MulScalar(&pair.g1, g1One, a)
		bls.G1.Affine(&pair.g1)

		b, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		bls.G2.MulScalar(&pair.g2, g2One, b)
		bls.G2.Affine(&pair.g2)
		c := new(big.Int).Mul(a, b)
		return pair, c
	}

	t.Run("Empty Pair", func(t *testing.T) {
		f := bls.pair(Pairs{})
		if !bls.fp12.equal(bls.fp12.one(), f) {
			t.Fatalf("bad pairing")
		}
	})

	t.Run("Multi Pair", func(t *testing.T) {
		pairSize := 50
		pairs := make(Pairs, pairSize)
		acc := new(big.Int)
		var c *big.Int
		for i := 0; i < pairSize; i++ {
			pairs[i], c = randomPair()
			acc.Add(acc, c)
		}
		f0 := bls.pair(pairs)
		f1 := bls.pair(generatorPair)
		bls.fp12.exp(f1, f1, acc)
		if !bls.fp12.equal(f0, f1) {
			t.Fatalf("bad pairing")
		}
	})

	t.Run("Bilinearity", func(t *testing.T) {
		var f1, f2 *fe12
		// e(a*G1, b*G2) = e(G1, G2)^c
		randomPair, c := randomPair()
		pair := Pairs{randomPair}
		t.Run("Bilinearity-1", func(t *testing.T) {
			f1 = bls.pair(generatorPair)
			f2 = bls.pair(pair)
			bls.fp12.exp(f1, f1, c)
			if !bls.fp12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(c*G1, G2)
		t.Run("Bilinearity-2", func(t *testing.T) {
			bls.G1.MulScalar(&generatorPair[0].g1, g1One, c)
			f2 = bls.pair(generatorPair)
			if !bls.fp12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
		// e(a*G1, b*G2) = e(G1, c*G2)
		t.Run("Bilinearity-3", func(t *testing.T) {
			bls.G1.Copy(&generatorPair[0].g1, g1One)
			bls.G2.MulScalar(&generatorPair[0].g2, g2One, c)
			f2 = bls.pair(generatorPair)
			if !bls.fp12.equal(f1, f2) {
				t.Errorf("bad pairing")
			}
		})
	})

	t.Run("Non-Degeneracy", func(t *testing.T) {
		pairs := Pairs{
			Pair{*g1Zero, *g2One},
			Pair{*g1One, *g2Zero},
			Pair{*g1Zero, *g2Zero},
		}
		if !bls.fp12.equal(bls.fp12.one(), bls.pair(pairs)) {
			t.Errorf("bad pairing")
		}
		pairs = Pairs{
			Pair{*g1Zero, *g2One},
			Pair{*g1One, *g2Zero},
			Pair{*g1Zero, *g2Zero},
			Pair{*g1One, *g2One},
		}
		if !bls.fp12.equal(expected, bls.pair(pairs)) {
			t.Errorf("bad pairing")
		}
	})

	t.Run("Pairing Check", func(t *testing.T) {
		g1NegOne := bls.G1.NegativeOne()
		g2NegOne := bls.G2.New()
		bls.G2.Neg(g2NegOne, g2One)
		pairs := Pairs{
			Pair{*g1One, *g2One},
			Pair{*g1NegOne, *g2One},
			Pair{*g1One, *g2One},
			Pair{*g1One, *g2NegOne},
		}
		if !bls.fp12.equal(bls.fp12.one(), bls.pair(pairs)) {
			t.Errorf("bad pairing")
		}
	})
}

func BenchmarkPairing(t *testing.B) {
	bls := NewBLSPairingEngine()
	g1One, g2One := bls.G1.One(), bls.G2.One()
	t.ResetTimer()
	generatorPair := Pairs{Pair{*g1One, *g2One}}
	for i := 0; i < t.N; i++ {
		bls.pair(generatorPair)
	}
}

func BenchmarkFinalExp(t *testing.B) {
	bls := NewBLSPairingEngine()
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bls.finalExp(bls.fp12.one())
	}

}
