package bls12381

import (
	"crypto/rand"
	"errors"
	"flag"
	"math/big"
	"testing"
)

var fuz int

func TestMain(m *testing.M) {
	_fuz := flag.Int("fuzz", 10, "# of iterations")
	adx := flag.Bool("noadx", false, "to enfoce non adx arch")
	flag.Parse()
	forceNonADXArch = *adx
	fuz = *_fuz
	cfgArch()
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

// func TestSwuMapperWithIETFVectors(t *testing.T) {
// 	mapper := newSwuMapper()
// 	for i, vector := range []struct {
// 		group string
// 		U     string
// 		P     string
// 	}{
// 		// G.9.2.  BLS12381G1_XMD:SHA-256_SSWU_NU_
// 		// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#appendix-G.9.2
// 		{
// 			U: "0ccb6bda9b602ab82aae21c0291623e2f639648a6ada1c76d8ffb664130fd18d98a2cc6160624148827a9726678e7cd4",
// 			P: "115281bd55a4103f31c8b12000d98149598b72e5da14e953277def263a24bc2e9fd8fa151df73ea3800f9c8cbb9b245c0796506faf9edbf1957ba8d667a079cab0d3a37e302e5132bd25665b66b26ea8556a0cfb92d6ae2c4890df0029b455ce",
// 		},
// 		{
// 			U: "08accd9a1bd4b75bb2e9f014ac354a198cbf607f0061d00a6286f5544cf4f9ecc1439e3194f570cbbc7b96d1a754f231",
// 			P: "04a7a63d24439ade3cd16eaab22583c95b061136bd5013cf109d92983f902c31f49c95cbeb97222577e571e97a68a32e09a8aa8d6e4b409bbe9a6976c016688269024d6e9d378ed25e8b4986194511f479228fa011ec88b8f4c57a621fc12187",
// 		},
// 		{
// 			U: "0a359cf072db3a39acf22f086d825fcf49d0daf241d98902342380fc5130b44e55de8f684f300bc11c44dee526413363",
// 			P: "05c59faaf88187f51cd9cc6c20ca47ac66cc38d99af88aef2e82d7f35104168916f200a79562e64bc843f83cdc8a46750b10472100a4aaa665f35f044b14a234b8f74990fa029e3dd06aa60b232fd9c232564ceead8cdb72a8a0320fc1071845",
// 		},
// 		{
// 			U: "181d09392c52f7740d5eaae52123c1dfa4808343261d8bdbaf19e7773e5cdfd989165cd9ecc795500e5da2437dde2093",
// 			P: "10147709f8d4f6f2fa6f957f6c6533e3bf9069c01be721f9421d88e0f02d8c617d048c6f8b13b81309d1ef6b56eeddc71048977c38688f1a3acf48ae319216cb1509b6a29bd1e7f3b2e476088a280e8c97d4a4c147f0203c7b3acb3caa566ae8",
// 		},
// 		// G.10.2.  BLS12381G2_XMD:SHA-256_SSWU_NU_
// 		// https://tools.ietf.org/html/draft-irtf-cfrg-hash-to-curve-06#appendix-G.10.2
// 		{
// 			U: "094376a68cdc8f64bd981d59bf762f9b2960df6b135f6e09ceada2fe8d0000bbf04023492796c09f8ef04016a2e8365f09367e3b485dda3925e82cc458e5009051281d3e442e94f9ef9feec44ee26375d6dc904dc1aa1f831f2aebd7b437ad12",
// 			P: "04264ddf941f7c9ea5ad62027c72b194c6c3f62a92fcdb56ddc9de7990489af1f81c576e7f451c2cd416102253e040f0170919c7845a9e623cef297e17484606a3eb2ae21ed8a21ff2b258861daefa3ac36955c0b374c6f4925868920d9c5f0b02d03d852629f70563e3a653ccc2e114439f551a2fd87c8136eb205b84e22c3f40507beccdcdc52c921b69a57968ec7c0ce03abe6c55ff0640b2b303440d88bd1a2b0cbfe3274b2802c1f58b1085e4dd8795c9c4d9c166d2f033e3c438e7f8a9",
// 		},
// 		{
// 			U: "0f105595e14847cc9a41fd70deb3240337678b266304100ec261add2585b991c7268bb1a325d2f871b327e8d04fd579b17ecd5d41a860b8886cb1210874b254f59945b089f774dcc14bc1aca7d4e3c975bce0d28510c442e9a932be5880ee5b1",
// 			P: "019a3b47aa956b2b548cc04d9e109dec06642d6e28814f7e35f807e1ce609e2eae3a155af406c842529776d8192f562e16d830a4e12fddfbdaf9a667f94f21e490879fd3ccc5ee6f039cd7c2174fb47ea8027af78779a978d2a921612844587f15adde069459ab2012b44c7703119185b96b7f04ad59b39f4f6aea35fdbb9c5c7d876b5f89afb55b67e7da96ad489dc315930174c11aa9b51a5cc3ebfa1ab6377e2318c4ea2df387bdb84b28687a02c86e6401b195bbcabb6e95d6ae43669e12",
// 		},
// 		{
// 			U: "1107a6f450c6c9580c720190b577f52c633cf5f3defb528ae873d3723bccc8fa433014e9120a1da31abc27c674f37ae4032ae17a23a76c94745a5460cd9f1191c0ebeec7adfc4df28b0833e536b7dbabf498dc076ff16cc11c6a6ef5105df693",
// 			P: "0910b2d55e210122fab2d2dae81e6a440fd22e925e422aaf16a8fd28477bacb12aa888de0faeea203e372a1c1cd9578c1498937f0ed18c49ebbcdee579b58ce235f3ab03be5dc809e1df25e2e0b4eb4c672f4eaf26df91f3755d6367df55d5be102631eb4e684d759312d7eab78598f487c2c10ad3d3552cb43ce6f09a11eb46e551864863077906d3ecfd921f1fe541033b1948575e70fed67fb4f7bd86b5452dfc0afeb74ecf5cab4a6872e33f0eade9564d3d5b9fcb9d4c498afda0bc037d",
// 		},
// 		{
// 			U: "0306162d24592a18fa8de2007d7b69d04bb7a71a5a7965d15bdcbaa4ddf9b599079fbdae9f67d55ab6dba044f9daf1790cda6b874f8c41862c078099aa76d607be51d913a2e3f997539a0993bda31892292818c74aa9be035f234df2576fe49a",
// 			P: "021f7faa0550e5a5d08338b4c0a5d30240dec7989fc7c77b6ffba9bfd5d64ce45af5aad8da8482bf0da91af4f29d371f18af6eedb7ed3be66c5a1d998ad4d9640f557b189558baec41f6e712ff2a39f795a35494b4b12343b7a1a2b17686d793166c1abec65af593d291dbd05e5d7d28f1a9ffb73751d65f49d76084493f3da707ee2bbf54cf6de5bbaac2ffa0028c310cc46cea229960bfbe25831162c27f96cf8bb14c017938e35b636987a306521915456fbd40633c6d5a30f61bce52a3f5",
// 		},
// 	} {
// 		if len(vector.U) <= 96 {
// 			u, _ := fromString(vector.U)
// 			expected, _ := mapper.g1.fromRawUnchecked(fromHex(-1, vector.P))
// 			actual, ok := mapper.toG1(u)
// 			if !ok {
// 				t.Errorf("g1 vector#%d test failed", i)
// 			} else {
// 				if !mapper.g1.Equal(actual, expected) {
// 					t.Errorf("g1 vector#%d point is not equal to expected", i)
// 				}
// 			}
// 		} else {
// 			u, _ := mapper.g2.f.fromBytes(fromHex(-1, vector.U))
// 			expected, err := mapper.g2.fromRawUnchecked(fromHex(-1, vector.P))
// 			if err != nil {
// 				t.Fatal(err)
// 			}
// 			actual, ok := mapper.toG2(u)
// 			if !ok {
// 				t.Errorf("g2 vector#%d test failed", i)
// 			} else {
// 				if !mapper.g2.Equal(actual, expected) {
// 					t.Errorf("g2 vector#%d point is not equal to expected", i)
// 				}
// 			}
// 		}

// 	}
// }
