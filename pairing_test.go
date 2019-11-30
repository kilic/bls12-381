package bls

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestPairing(t *testing.T) {
	e := NewBLSPairingEngine()
	t.Run("Single Expected", func(t *testing.T) {
		G, H := new(PointG1), new(PointG2)
		e.G1.Copy(G, e.G1.One())
		e.G2.Copy(H, e.G2.One())
		f1 := e.Pair(G, H)
		f2 := &fe12{}
		if err := e.fp12.newElementFromBytes(
			f2,
			bytes_(48,
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
			)); err != nil {
			t.Fatal(err)
		}
		if !e.fp12.equal(f1, f2) {
			t.Fatal("bad pairing")
		}
	})
}

func TestBilinearity(t *testing.T) {
	e := NewBLSPairingEngine()
	a, _ := rand.Int(rand.Reader, big.NewInt(100))
	b, _ := rand.Int(rand.Reader, big.NewInt(100))
	c := new(big.Int).Mul(a, b)
	G, H := new(PointG1), new(PointG2)
	e.G1.MulScalar(G, e.G1.One(), a)
	e.G2.MulScalar(H, e.G2.One(), b)

	var f1, f2 *fe12
	// e(a*G1, b*G2) = e(G1, G2)^c
	t.Run("Bilinearity-1", func(t *testing.T) {
		f1 = e.Pair(G, H)
		f2 = e.Pair(e.G1.One(), e.G2.One())
		e.fp12.exp(f2, f2, c)
		if !e.fp12.equal(f1, f2) {
			t.Errorf("bad pairing")
		}
	})
	// e(a*G1, b*G2) = e(c*G1, G2)
	t.Run("Bilinearity-2", func(t *testing.T) {
		G = e.G1.MulScalar(G, e.G1.One(), c)
		f2 = e.Pair(G, e.G2.One())
		if !e.fp12.equal(f1, f2) {
			t.Errorf("bad pairing")
		}
	})
	// e(a*G1, b*G2) = e(G1, c*G2)
	t.Run("Bilinearity-3", func(t *testing.T) {
		H = e.G2.MulScalar(H, e.G2.One(), c)
		f2 = e.Pair(e.G1.One(), H)
		if !e.fp12.equal(f1, f2) {
			t.Errorf("bad pairing")
		}
	})
}
func TestPairingMulti(t *testing.T) {
	e := NewBLSPairingEngine()
	G := &PointG1{}
	H := &PointG2{}

	g1RandPoint := func() (*PointG1, *big.Int) {
		s, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		p := e.G1.MulScalar(&PointG1{}, G, s)
		e.G1.Affine(p)
		return p, s
	}
	g2RandPoint := func() (*PointG2, *big.Int) {
		s, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		p := e.G2.MulScalar(&PointG2{}, H, s)
		e.G2.Affine(p)
		return p, s
	}
	pairSize := 50
	points := make([]PointG1, pairSize)
	twistPoints := make([]PointG2, pairSize)
	acc := new(big.Int)
	for i := 0; i < pairSize; i++ {
		aG, a := g1RandPoint()
		bH, b := g2RandPoint()
		e.G1.Copy(&points[i], aG)
		e.G2.Copy(&twistPoints[i], bH)
		acc.Add(acc, new(big.Int).Mul(a, b))
	}
	f0 := e.MultiPair(points, twistPoints)
	f1 := e.Pair(G, H)
	e.fp12.exp(f1, f1, acc)
	if !e.fp12.equal(f0, f1) {
		t.Fatalf("bad pairing")
	}
}

func BenchmarkPairing(t *testing.B) {
	e := NewBLSPairingEngine()
	G := &PointG1{}
	H := &PointG2{}
	e.G1.Copy(G, &g1One)
	e.G2.Copy(H, &g2One)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		e.Pair(G, H)
	}
}

func BenchmarkFinalExp(t *testing.B) {
	e := NewBLSPairingEngine()
	a := fe12{}
	e.fp12.randElement(&a, rand.Reader)
	t.Run("1", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			e.finalExp(&a)
		}
	})
	t.Run("2", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			e.finalExp2(&a)
		}
	})
}
