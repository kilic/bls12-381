package bls

import (
	"crypto/rand"
	"io/ioutil"
	"math/big"
	"testing"
)

func (g *G2) one() *PointG2 {
	one, err := g.fromBytesUnchecked(fromHex(48,
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

func TestG2Serialization(t *testing.T) {
	var err error
	g2 := NewG2(nil)
	zero := g2.Zero()
	b0 := g2.ToUncompressed(zero)
	p0, err := g2.FromUncompressed(b0)
	if err != nil {
		t.Fatal(err)
	}
	if !g2.IsZero(p0) {
		t.Fatalf("bad infinity serialization 1")
	}
	b0 = g2.ToCompressed(zero)
	p0, err = g2.FromCompressed(b0)
	if err != nil {
		t.Fatal(err)
	}
	if !g2.IsZero(p0) {
		t.Fatalf("bad infinity serialization 2")
	}
	b0 = g2.ToBytes(zero)
	p0, err = g2.FromBytes(b0)
	if err != nil {
		t.Fatal(err)
	}
	if !g2.IsZero(p0) {
		t.Fatalf("bad infinity serialization 3")
	}
	for i := 0; i < fuz; i++ {
		a := g2.rand()
		uncompressed := g2.ToUncompressed(a)
		b, err := g2.FromUncompressed(uncompressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g2.Equal(a, b) {
			t.Fatalf("bad serialization 1")
		}
		compressed := g2.ToCompressed(b)
		a, err = g2.FromCompressed(compressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g2.Equal(a, b) {
			t.Fatalf("bad serialization 2")
		}
	}
	for i := 0; i < fuz; i++ {
		a := g2.rand()
		uncompressed := g2.ToBytes(a)
		b, err := g2.FromBytes(uncompressed)
		if err != nil {
			t.Fatal(err)
		}
		if !g2.Equal(a, b) {
			t.Fatalf("bad serialization 3")
		}
	}
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

func TestG2SWUMap(t *testing.T) {
	// G.10.2.  BLS12381G2_XMD:SHA-256_SSWU_NU_
	// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#appendix-G.10.2
	for i, v := range []struct {
		U []byte
		P []byte
	}{
		{
			U: fromHex(-1, "094376a68cdc8f64bd981d59bf762f9b2960df6b135f6e09ceada2fe8d0000bbf04023492796c09f8ef04016a2e8365f09367e3b485dda3925e82cc458e5009051281d3e442e94f9ef9feec44ee26375d6dc904dc1aa1f831f2aebd7b437ad12"),
			P: fromHex(-1, "04264ddf941f7c9ea5ad62027c72b194c6c3f62a92fcdb56ddc9de7990489af1f81c576e7f451c2cd416102253e040f0170919c7845a9e623cef297e17484606a3eb2ae21ed8a21ff2b258861daefa3ac36955c0b374c6f4925868920d9c5f0b02d03d852629f70563e3a653ccc2e114439f551a2fd87c8136eb205b84e22c3f40507beccdcdc52c921b69a57968ec7c0ce03abe6c55ff0640b2b303440d88bd1a2b0cbfe3274b2802c1f58b1085e4dd8795c9c4d9c166d2f033e3c438e7f8a9"),
		},
		{
			U: fromHex(-1, "0f105595e14847cc9a41fd70deb3240337678b266304100ec261add2585b991c7268bb1a325d2f871b327e8d04fd579b17ecd5d41a860b8886cb1210874b254f59945b089f774dcc14bc1aca7d4e3c975bce0d28510c442e9a932be5880ee5b1"),
			P: fromHex(-1, "019a3b47aa956b2b548cc04d9e109dec06642d6e28814f7e35f807e1ce609e2eae3a155af406c842529776d8192f562e16d830a4e12fddfbdaf9a667f94f21e490879fd3ccc5ee6f039cd7c2174fb47ea8027af78779a978d2a921612844587f15adde069459ab2012b44c7703119185b96b7f04ad59b39f4f6aea35fdbb9c5c7d876b5f89afb55b67e7da96ad489dc315930174c11aa9b51a5cc3ebfa1ab6377e2318c4ea2df387bdb84b28687a02c86e6401b195bbcabb6e95d6ae43669e12"),
		},
		{
			U: fromHex(-1, "1107a6f450c6c9580c720190b577f52c633cf5f3defb528ae873d3723bccc8fa433014e9120a1da31abc27c674f37ae4032ae17a23a76c94745a5460cd9f1191c0ebeec7adfc4df28b0833e536b7dbabf498dc076ff16cc11c6a6ef5105df693"),
			P: fromHex(-1, "0910b2d55e210122fab2d2dae81e6a440fd22e925e422aaf16a8fd28477bacb12aa888de0faeea203e372a1c1cd9578c1498937f0ed18c49ebbcdee579b58ce235f3ab03be5dc809e1df25e2e0b4eb4c672f4eaf26df91f3755d6367df55d5be102631eb4e684d759312d7eab78598f487c2c10ad3d3552cb43ce6f09a11eb46e551864863077906d3ecfd921f1fe541033b1948575e70fed67fb4f7bd86b5452dfc0afeb74ecf5cab4a6872e33f0eade9564d3d5b9fcb9d4c498afda0bc037d"),
		},
		{
			U: fromHex(-1, "0306162d24592a18fa8de2007d7b69d04bb7a71a5a7965d15bdcbaa4ddf9b599079fbdae9f67d55ab6dba044f9daf1790cda6b874f8c41862c078099aa76d607be51d913a2e3f997539a0993bda31892292818c74aa9be035f234df2576fe49a"),
			P: fromHex(-1, "021f7faa0550e5a5d08338b4c0a5d30240dec7989fc7c77b6ffba9bfd5d64ce45af5aad8da8482bf0da91af4f29d371f18af6eedb7ed3be66c5a1d998ad4d9640f557b189558baec41f6e712ff2a39f795a35494b4b12343b7a1a2b17686d793166c1abec65af593d291dbd05e5d7d28f1a9ffb73751d65f49d76084493f3da707ee2bbf54cf6de5bbaac2ffa0028c310cc46cea229960bfbe25831162c27f96cf8bb14c017938e35b636987a306521915456fbd40633c6d5a30f61bce52a3f5"),
		},
	} {
		g := NewG2(nil)
		p0, err := g.MapToPointSWU(v.U)
		p1, _ := g.fromBytesUnchecked(v.P)
		if err != nil {
			t.Fatal("swu mapping fails", i, err)
		}
		if !g.Equal(p0, p1) {
			t.Fatal("bad swu mapping", i, p1)
		}
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

func BenchmarkG2SWUMap(t *testing.B) {
	a := fromHex(96, "0x1234")
	g2 := NewG2(nil)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		_, err := g2.MapToPointSWU(a)
		if err != nil {
			t.Fatal(err)
		}
	}
}
