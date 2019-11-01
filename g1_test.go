package bls

import (
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"testing"
)

func TestG1(t *testing.T) {
	g1 := NewG1(newFp())
	one := g1.fromRawUnchecked(bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	))
	randPoint := func() *PointG1 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g1.MulScalar(&PointG1{}, one, k)
	}
	zero := g1.Zero()
	var t0, t1 PointG1
	t.Run("Generator", func(t *testing.T) {
		if !g1.IsOnCurve(one) {
			t.Fatalf("generator is not on curve")
		}
	})
	t.Run("Serialization", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a := randPoint()
			uncompressed := g1.ToUncompressed(a)
			b, err := g1.FromUncompressed(uncompressed)
			if err != nil {
				t.Fatal(err)
			}
			if !g1.Equal(a, b) {
				t.Fatalf("bad encoding & decoding 1")
			}
			compressed := g1.ToCompressed(b)
			a, err = g1.FromCompressed(compressed)
			if err != nil {
				t.Fatal(err)
			}
			if !g1.Equal(a, b) {
				t.Fatalf("bad encoding & decoding 2")
			}

		}
	})
	t.Run("Addition", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a, b := randPoint(), randPoint()
			g1.Add(&t0, a, one)
			g1.Add(&t0, &t0, b)
			g1.Add(&t1, one, b)
			g1.Add(&t1, &t1, a)
			if !g1.Equal(&t0, &t1) || !g1.IsOnCurve(&t1) || !g1.IsOnCurve(&t0) {
				t.Fatalf("")
			}
			g1.Add(b, a, zero)
			if !g1.Equal(a, b) || !g1.IsOnCurve(b) {
				t.Fatalf("")
			}
		}
	})
	t.Run("Doubling", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a := randPoint()
			g1.Double(&t0, a)
			g1.Sub(&t0, &t0, a)
			if !g1.Equal(&t0, a) || !g1.IsOnCurve(&t0) {
				t.Fatalf("")
			}
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		for i := 0; i < n; i++ {
			s1, s2, s3 := randScalar(q), randScalar(q), randScalar(q)
			g1.MulScalar(&t0, one, s1)
			g1.MulScalar(&t0, &t0, s2)
			s3.Mul(s1, s2)
			g1.MulScalar(&t1, one, s3)
			if !g1.Equal(&t0, &t1) || !g1.IsOnCurve(&t1) || !g1.IsOnCurve(&t0) {
				t.Errorf("")
			}
			a := randPoint()
			g1.MulScalar(a, a, big.NewInt(0))
			if !g1.Equal(a, zero) || !g1.IsOnCurve(a) {
				t.Errorf("")
			}
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		for i := 0; i < n; i++ {
			a, s := randPoint(), randScalar(q)
			var p PointG1
			g1.MulScalar(&p, a, new(big.Int).ModInverse(s, q))
			g1.MulScalar(&p, &p, s)
			if !g1.Equal(&p, a) || !g1.IsOnCurve(&p) || !g1.IsOnCurve(a) {
				t.Errorf("")
				return
			}
		}
	})
}

func TestZKCryptoVectors_G1UncompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g1_uncompressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG1(nil)
	p1 := g.Zero()
	for i := 0; i < 1000; i++ {
		vector := data[i*96 : (i+1)*96]
		p2, err := g.FromUncompressed(vector)
		if err != nil {
			t.Fatal("decoing fails", err, i)
		}
		if !g.Equal(p1, p2) {
			t.Fatalf("\nwant: %s\nhave: %s\n", p1, p2)
		}
		g.Add(p1, p1, &g1One)
	}
}

func TestZKCryptoVectors_G1CompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g1_compressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG1(nil)
	p1 := g.Zero()
	for i := 0; i < 1000; i++ {
		vector := data[i*48 : (i+1)*48]
		p2, err := g.FromCompressed(vector)
		if err != nil {
			t.Fatal("decoing fails", err, i)
		}
		if !g.Equal(p1, p2) {
			t.Fatalf("\nwant: %s\nhave: %s\n", p1, p2)
		}
		g.Add(p1, p1, &g1One)
	}
}

func BenchmarkG1Add(t *testing.B) {
	g1 := NewG1(newFp())
	one := g1.fromRawUnchecked(bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	))
	randPoint := func() *PointG1 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g1.MulScalar(&PointG1{}, one, k)
	}
	a, b, c := randPoint(), randPoint(), PointG1{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g1.Add(&c, a, b)
	}
}

func BenchmarkG1Mul(t *testing.B) {
	g1 := NewG1(newFp())
	one := g1.fromRawUnchecked(bytes_(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	))
	randPoint := func() *PointG1 {
		k, err := rand.Int(rand.Reader, q)
		if err != nil {
			panic(err)
		}
		return g1.MulScalar(&PointG1{}, one, k)
	}
	a, e, c := randPoint(), q, PointG1{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g1.MulScalar(&c, a, e)
	}
}
