package bls

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestG2(t *testing.T) {
	g2 := NewG2(NewFp2(NewFp()))
	one := g2.fromRawUnchecked(bytes_(48,
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
	))
	randPoint := func() *PointG2 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g2.MulScalar(&PointG2{}, one, k)
	}
	zero := g2.Zero()
	var t0, t1 PointG2
	t.Run("Generator", func(t *testing.T) {
		if !g2.IsOnCurve(one) {
			t.Fatalf("generator is not on curve")
		}
	})
	t.Run("Serialization", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a := randPoint()
			uncompressed := g2.ToUncompressed(a)
			b, err := g2.FromUncompressed(uncompressed)
			if err != nil {
				t.Fatal(err)
			}
			if !g2.Equal(a, b) {
				t.Fatalf("bad encoding & decoding 1")
			}
			compressed := g2.ToCompressed(b)
			a, err = g2.FromCompressed(compressed)
			if err != nil {
				t.Fatal(err)
			}
			if !g2.Equal(a, b) {
				t.Fatalf("bad encoding & decoding 2")
			}

		}
	})
	t.Run("Addition", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a, b := randPoint(), randPoint()
			g2.Add(&t0, a, one)
			g2.Add(&t0, &t0, b)
			g2.Add(&t1, one, b)
			g2.Add(&t1, &t1, a)
			if !g2.Equal(&t0, &t1) || !g2.IsOnCurve(&t1) || !g2.IsOnCurve(&t0) {
				t.Fatalf("")
			}
			g2.Add(b, a, zero)
			if !g2.Equal(a, b) || !g2.IsOnCurve(b) {
				t.Fatalf("")
			}
		}
	})
	t.Run("Doubling", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a := randPoint()
			g2.Double(&t0, a)
			g2.Sub(&t0, &t0, a)
			if !g2.Equal(&t0, a) || !g2.IsOnCurve(&t0) {
				t.Fatalf("")
			}
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		for i := 0; i < n; i++ {
			s1, s2, s3 := randScalar(q), randScalar(q), randScalar(q)
			g2.MulScalar(&t0, one, s1)
			g2.MulScalar(&t0, &t0, s2)
			s3.Mul(s1, s2)
			g2.MulScalar(&t1, one, s3)
			if !g2.Equal(&t0, &t1) || !g2.IsOnCurve(&t1) || !g2.IsOnCurve(&t0) {
				t.Errorf("")
			}
			a := randPoint()
			g2.MulScalar(a, a, big.NewInt(0))
			if !g2.Equal(a, zero) || !g2.IsOnCurve(a) {
				t.Errorf("")
			}
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a, s := randPoint(), randScalar(q)
			var p PointG2
			g2.MulScalar(&p, a, new(big.Int).ModInverse(s, q))
			g2.MulScalar(&p, &p, s)
			if !g2.Equal(&p, a) || !g2.IsOnCurve(&p) || !g2.IsOnCurve(a) {
				t.Errorf("")
				return
			}
		}
	})
}

func BenchmarkG2Add(t *testing.B) {
	g2 := NewG2(NewFp2(NewFp()))
	one := g2.fromRawUnchecked(bytes_(48,
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
	))
	randPoint := func() *PointG2 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g2.MulScalar(&PointG2{}, one, k)
	}
	a, b, c := randPoint(), randPoint(), PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g2.Add(&c, a, b)
	}
}

func BenchmarkG2Mul(t *testing.B) {
	g2 := NewG2(NewFp2(NewFp()))
	one := g2.fromRawUnchecked(bytes_(48,
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
	))
	randPoint := func() *PointG2 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g2.MulScalar(&PointG2{}, one, k)
	}
	a, e, c := randPoint(), q, PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g2.MulScalar(&c, a, e)
	}
}
