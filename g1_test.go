package bls

import (
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"testing"
)

func (g *G1) one() *PointG1 {
	one, err := g.fromRawUnchecked(fromHex(48,
		"0x17f1d3a73197d7942695638c4fa9ac0fc3688c4f9774b905a14e3a3f171bac586c55e83ff97a1aeffb3af00adb22c6bb",
		"0x08b3f481e3aaa0f1a09e30ed741d8ae4fcf5e095d5d00af600db18cb2c04b3edd03cc744a2888ae40caa232946c5e7e1",
	))
	if err != nil {
		panic(err)
	}
	return one
}

func (g *G1) rand() *PointG1 {
	k, err := rand.Int(rand.Reader, q)
	if err != nil {
		panic(err)
	}
	return g.MulScalar(&PointG1{}, g.one(), k)
}

func (g *G1) randAffine() *PointG1 {
	return g.Affine(g.rand())
}

func TestG1Serialization(t *testing.T) {
	g1 := NewG1()
	for i := 0; i < fuz; i++ {
		a := g1.rand()
		uncompressed := g1.ToUncompressed(a)
		b, err := g1.FromUncompressed(uncompressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g1.Equal(a, b) {
			t.Fatalf("bad serialization 1")
		}
		compressed := g1.ToCompressed(b)
		a, err = g1.FromCompressed(compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g1.Equal(a, b) {
			t.Fatalf("bad serialization 2")
		}
	}
}

func TestG1AdditiveProperties(t *testing.T) {
	g := NewG1()
	t0, t1 := g.New(), g.New()
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

func TestG1MultiplicativeProperties(t *testing.T) {
	g := NewG1()
	t0, t1 := g.New(), g.New()
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

func TestZKCryptoVectorsG1UncompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g1_uncompressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG1()
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

func TestZKCryptoVectorsG1CompressedValid(t *testing.T) {
	data, err := ioutil.ReadFile("tests/g1_compressed_valid_test_vectors.dat")
	if err != nil {
		panic(err)
	}
	g := NewG1()
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

func TestG1MultiExpExpected(t *testing.T) {
	g := NewG1()
	one := g.one()
	var scalars [2]*big.Int
	var bases [2]*PointG1
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

func TestG1MultiExpBatch(t *testing.T) {
	g := NewG1()
	one := g.one()
	n := 1000
	bases := make([]*PointG1, n)
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

func TestG1SWUMap(t *testing.T) {
	// G.9.2.  BLS12381G1_XMD:SHA-256_SSWU_NU_
	// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#appendix-G.9.2
	for i, v := range []struct {
		U []byte
		P []byte
	}{
		{
			U: fromHex(-1, "0x0ccb6bda9b602ab82aae21c0291623e2f639648a6ada1c76d8ffb664130fd18d98a2cc6160624148827a9726678e7cd4"),
			P: fromHex(-1, "0x115281bd55a4103f31c8b12000d98149598b72e5da14e953277def263a24bc2e9fd8fa151df73ea3800f9c8cbb9b245c0796506faf9edbf1957ba8d667a079cab0d3a37e302e5132bd25665b66b26ea8556a0cfb92d6ae2c4890df0029b455ce"),
		},
		{
			U: fromHex(-1, "0x08accd9a1bd4b75bb2e9f014ac354a198cbf607f0061d00a6286f5544cf4f9ecc1439e3194f570cbbc7b96d1a754f231"),
			P: fromHex(-1, "0x04a7a63d24439ade3cd16eaab22583c95b061136bd5013cf109d92983f902c31f49c95cbeb97222577e571e97a68a32e09a8aa8d6e4b409bbe9a6976c016688269024d6e9d378ed25e8b4986194511f479228fa011ec88b8f4c57a621fc12187"),
		},
		{
			U: fromHex(-1, "0x0a359cf072db3a39acf22f086d825fcf49d0daf241d98902342380fc5130b44e55de8f684f300bc11c44dee526413363"),
			P: fromHex(-1, "0x05c59faaf88187f51cd9cc6c20ca47ac66cc38d99af88aef2e82d7f35104168916f200a79562e64bc843f83cdc8a46750b10472100a4aaa665f35f044b14a234b8f74990fa029e3dd06aa60b232fd9c232564ceead8cdb72a8a0320fc1071845"),
		},
		{
			U: fromHex(-1, "0x181d09392c52f7740d5eaae52123c1dfa4808343261d8bdbaf19e7773e5cdfd989165cd9ecc795500e5da2437dde2093"),
			P: fromHex(-1, "0x10147709f8d4f6f2fa6f957f6c6533e3bf9069c01be721f9421d88e0f02d8c617d048c6f8b13b81309d1ef6b56eeddc71048977c38688f1a3acf48ae319216cb1509b6a29bd1e7f3b2e476088a280e8c97d4a4c147f0203c7b3acb3caa566ae8"),
		},
	} {
		g := NewG1()
		p0, err := g.MapToPointSWU(v.U)
		p1, _ := g.fromRawUnchecked(v.P)
		if err != nil {
			t.Fatal("swu mapping fails", i, err)
		}
		if !g.Equal(p0, p1) {
			t.Fatal("bad swu mapping", i, p1)
		}
	}
}

func BenchmarkG1Add(t *testing.B) {
	g1 := NewG1()
	a, b, c := g1.rand(), g1.rand(), PointG1{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g1.Add(&c, a, b)
	}
}

func BenchmarkG1Mul(t *testing.B) {
	g1 := NewG1()
	a, e, c := g1.rand(), q, PointG1{}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		g1.MulScalar(&c, a, e)
	}
}
