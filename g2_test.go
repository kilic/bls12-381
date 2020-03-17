package bls

import (
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"testing"
)

func (g *G2) one() *PointG2 {
	one, err := g.fromRawUnchecked(fromHex(48,
		"0x13e02b6052719f607dacd3a088274f65596bd0d09920b61ab5da61bbdc7f5049334cf11213945d57e5ac7d055d042b7e",
		"0x024aa2b2f08f0a91260805272dc51051c6e47ad4fa403b02b4510b647ae3d1770bac0326a805bbefd48056c8c121bdb8",
		"0x0606c4a02ea734cc32acd2b02bc28b99cb3e287e85a763af267492ab572e99ab3f370d275cec1da1aaa9075ff05f79be",
		"0x0ce5d527727d6e118cc9cdc6da2e351aadfd9baa8cbdd3a76d429a695160d12c923ac9cc3baca289e193548608b82801",
	))
	if err != nil {
		panic(err)
	}
	return one
}

func (g *G2) rand() *PointG2 {
	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		panic(err)
	}
	return g.MulScalar(&PointG2{}, g.one(), k)
}

func (g *G2) randAffine() *PointG2 {
	return g.Affine(g.rand())
}

func (g *G2) new() *PointG2 {
	return g.Zero()
}

func TestG2AdditiveProperties(t *testing.T) {
	g := NewG2(newFp2())
	t0, t1 := g.new(), g.new()
	zero := g.Zero()
	for i := 0; i < fuz; i++ {
		a, b := g.rand(), g.rand()
		g.Add(t0, a, zero)
		if !g.Equal(t0, a) {
			t.Fatalf("a + 0 == a")
		}
		g.Add(t0, zero, zero)
		if !g.Equal(t0, zero) {
			t.Fatalf("0 + 0 == 0")
		}
		g.Sub(t0, a, zero)
		if !g.Equal(t0, a) {
			t.Fatalf("a - 0 == a")
		}
		g.Sub(t0, zero, zero)
		if !g.Equal(t0, zero) {
			t.Fatalf("0 - 0 == 0")
		}
		g.Neg(t0, zero)
		if !g.Equal(t0, zero) {
			t.Fatalf("- 0 == 0")
		}
		g.Sub(t0, zero, a)
		g.Neg(t0, t0)
		if !g.Equal(t0, a) {
			t.Fatalf(" - (0 - a) == a")
		}
		g.Double(t0, zero)
		if !g.Equal(t0, zero) {
			t.Fatalf("2 * 0 == 0")
		}
		g.Double(t0, a)
		g.Sub(t0, t0, a)
		if !g.Equal(t0, a) || !g.IsOnCurve(t0) {
			t.Fatalf(" (2 * a) - a == a")
		}
		g.Add(t0, a, b)
		g.Add(t1, b, a)
		if !g.Equal(t0, t1) {
			t.Fatalf("a + b == b + a")
		}
		g.Sub(t0, a, b)
		g.Sub(t1, b, a)
		g.Neg(t1, t1)
		if !g.Equal(t0, t1) {
			t.Fatalf("a - b == - ( b - a )")
		}
		c := g.rand()
		g.Add(t0, a, b)
		g.Add(t0, t0, c)
		g.Add(t1, a, c)
		g.Add(t1, t1, b)
		if !g.Equal(t0, t1) {
			t.Fatalf("(a + b) + c == (a + c ) + b")
		}
		g.Sub(t0, a, b)
		g.Sub(t0, t0, c)
		g.Sub(t1, a, c)
		g.Sub(t1, t1, b)
		if !g.Equal(t0, t1) {
			t.Fatalf("(a - b) - c == (a - c) -b")
		}
	}
}

func TestG2MultiplicativeProperties(t *testing.T) {
	g := NewG2(newFp2())
	t0, t1 := g.new(), g.new()
	zero := g.Zero()
	for i := 0; i < fuz; i++ {
		a := g.rand()
		s1, s2, s3 := randScalar(q), randScalar(q), randScalar(q)
		sone := big.NewInt(1)
		g.MulScalar(t0, zero, s1)
		if !g.Equal(t0, zero) {
			t.Fatalf(" 0 ^ s == 0")
		}
		g.MulScalar(t0, a, sone)
		if !g.Equal(t0, a) {
			t.Fatalf(" a ^ 1 == a")
		}
		g.MulScalar(t0, zero, s1)
		if !g.Equal(t0, zero) {
			t.Fatalf(" 0 ^ s == a")
		}
		g.MulScalar(t0, a, s1)
		g.MulScalar(t0, t0, s2)
		s3.Mul(s1, s2)
		g.MulScalar(t1, a, s3)
		if !g.Equal(t0, t1) {
			t.Errorf(" (a ^ s1) ^ s2 == a ^ (s1 * s2)")
		}
		g.MulScalar(t0, a, s1)
		g.MulScalar(t1, a, s2)
		g.Add(t0, t0, t1)
		s3.Add(s1, s2)
		g.MulScalar(t1, a, s3)
		if !g.Equal(t0, t1) {
			t.Errorf(" (a ^ s1) + (a ^ s2) == a ^ (s1 + s2)")
		}
	}
}

func TestZKCryptoVectorsG2UncompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g2_uncompressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG2(nil)
	p1 := g.Zero()
	for i := 0; i < 1000; i++ {
		vector := data[i*192 : (i+1)*192]
		p2, err := g.FromUncompressed(vector)
		if err != nil {
			t.Fatal("decoing fails", err, i)
		}
		if !g.Equal(p1, p2) {
			t.Fatalf("\nwant: %s\nhave: %s\n", p1, p2)
		}
		g.Add(p1, p1, &g2One)
	}
}

func TestZKCryptoVectorsG2CompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g2_compressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG2(nil)
	p1 := g.Zero()
	for i := 0; i < 1000; i++ {
		vector := data[i*96 : (i+1)*96]
		p2, err := g.FromCompressed(vector)
		if err != nil {
			t.Fatal("decoing fails", err, i)
		}
		if !g.Equal(p1, p2) {
			t.Fatalf("\nwant: %s\nhave: %s\n", p1, p2)
		}
		g.Add(p1, p1, &g2One)
	}
}

func TestG2MultiExpExpected(t *testing.T) {
	g := NewG2(nil)
	one := g.one()
	var scalars [2]*big.Int
	var bases [2]*PointG2
	scalars[0] = big.NewInt(2)
	scalars[1] = big.NewInt(3)
	bases[0], bases[1] = g.New(), g.New()
	g.Copy(bases[0], one)
	g.Copy(bases[1], one)
	expected, result := g.New(), g.New()
	g.MulScalar(expected, one, big.NewInt(5))
	g.MultiExp(result, bases[:], scalars[:])
	if !g.Equal(expected, result) {
		t.Fatalf("bad multi-exponentiation")
	}
}

func TestG2MultiExpBatch(t *testing.T) {
	g := NewG2(nil)
	one := g.one()
	n := 1000
	bases := make([]*PointG2, n)
	scalars := make([]*big.Int, n)
	// scalars: [s0,s1 ... s(n-1)]
	// bases: [P0,P1,..P(n-1)] = [s(n-1)*G, s(n-2)*G ... s0*G]
	for i, j := 0, n-1; i < n; i, j = i+1, j-1 {
		scalars[j], _ = rand.Int(rand.Reader, big.NewInt(100000))
		bases[i] = g.New()
		g.MulScalar(bases[i], one, scalars[j])
	}
	// expected: s(n-1)*P0 + s(n-2)*P1 + s0*P(n-1)
	expected, tmp := g.New(), g.New()
	for i := 0; i < n; i++ {
		g.MulScalar(tmp, bases[i], scalars[i])
		g.Add(expected, expected, tmp)
	}
	result := g.New()
	g.MultiExp(result, bases, scalars)
	if !g.Equal(expected, result) {
		t.Fatalf("bad multi-exponentiation")
	}
}

func BenchmarkG2Add(t *testing.B) {
	g2 := NewG2(newFp2())
	a, b, c := g2.rand(), g2.rand(), PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g2.Add(&c, a, b)
	}
}

func BenchmarkG2Mul(t *testing.B) {
	g2 := NewG2(newFp2())
	a, e, c := g2.rand(), q, PointG2{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g2.MulScalar(&c, a, e)
	}
}
