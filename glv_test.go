package bls12381

import (
	"crypto/rand"
	"math/big"
	"testing"
)

func TestGLVConstruction(t *testing.T) {

	t.Run("Parameters", func(t *testing.T) {

		t0, t1 := new(Fr), new(Fr)
		one := new(Fr).setUInt64(1)
		t0.Square(glvLambda1)
		t0.Add(t0, glvLambda1)
		t1.Sub(q, one)
		if !t0.Equal(t1) {
			t.Fatal("lambda**2 + lambda + 1 = 0")
		}
		// t0.Square(glvLambda1Alt)
		// t0.Add(t0, glvLambda1Alt)
		// t1.Sub(q, one)
		// if !t0.Equal(t1) {
		// 	t.Fatal("lambda**2 + lambda + 1 = 0")
		// }
		c0 := new(fe)
		square(c0, glvPhi)
		mul(c0, c0, glvPhi)
		if !c0.isOne() {
			t.Fatal("phi^3 = 1")
		}
		// square(c0, glvPhiAlt)
		// mul(c0, c0, glvPhiAlt)
		// if !c0.isOne() {
		// 	t.Fatal("phi^3 = 1")
		// }
	})
	t.Run("Endomorphism G1", func(t *testing.T) {

		g := NewG1()
		{
			p0, p1 := g.randAffine(), g.New()
			g.MulScalar(p1, p0, glvLambda1)
			g.Affine(p1)
			x := new(fe)
			mul(x, &p0[0], glvPhi)
			if !x.equal(&p1[0]) || !p0[1].equal(&p1[1]) {
				t.Fatal("f(x, y) = (phi * x, y)")
			}
		}
		// {
		// 	p0, p1 := g.randAffine(), g.New()
		// 	g.MulScalar(p1, p0, glvLambda1Alt)
		// 	g.Affine(p1)
		// 	x := new(fe)
		// 	mul(x, &p0[0], glvPhiAlt)
		// 	if !x.equal(&p1[0]) || !p0[1].equal(&p1[1]) {
		// 		t.Fatal("f(x, y) = (phi * x, y)")
		// 	}
		// }
	})
	t.Run("Scalar Decomposition", func(t *testing.T) {
		for i := 0; i < fuz; i++ {
			m, err := new(Fr).Rand(rand.Reader)
			if err != nil {
				t.Fatal(err)
			}
			mBig := m.ToBig()

			var v *glvVector
			{
				v = decompose(m)

				if v.m1.Cmp(r128) >= 0 {
					t.Fatal("bad scalar component, m1")
				}
				if v.m2.Cmp(r128) >= 0 {
					t.Fatal("bad scalar component, m2")
				}

				k := new(Fr)
				if v.neg1 && v.neg2 {
					k.Mul(glvLambda1, v.m2)
					k.Sub(k, v.m1)
				} else if v.neg1 {
					k.Mul(glvLambda1, v.m2)
					k.Add(k, v.m1)
					k.Neg(k)
				} else if v.neg2 {
					k.Mul(glvLambda1, v.m2)
					k.Add(v.m1, k)
				} else {
					k.Mul(glvLambda1, v.m2)
					k.Sub(v.m1, k)
				}

				if !k.Equal(m) {
					t.Fatal("scalar decomposing failed")
				}
			}

			k1Abs, k2Abs, k1Big, k2Big := new(big.Int), new(big.Int), new(big.Int), new(big.Int)
			r128Big := r128.ToBig()
			{
				k1Big, k2Big = decomposeBig(mBig)
				k1Abs.Abs(k1Big)
				k2Abs.Abs(k2Big)
				if k1Abs.Cmp(r128Big) >= 0 {
					t.Fatal("bad scalar component, big m1")
				}
				if k2Abs.Cmp(r128Big) >= 0 {
					t.Fatal("bad scalar component, big m2")
				}

				k := new(big.Int)
				k.Mul(glvLambda1Big, k2Big)
				k.Sub(k1Big, k).Mod(k, qBig)
				if k.Cmp(mBig) != 0 {
					t.Fatal("scalar decomposing with big.Int failed", i)
				}

			}

			zeroBig := new(big.Int)

			if v.neg1 != (k1Big.Cmp(zeroBig) == -1) {
				t.Fatal("cross: scalar decomposing with failed neg1")
			}
			if v.neg2 != (k2Big.Cmp(zeroBig) == -1) {
				t.Fatal("cross: scalar decomposing with failed neg2")
			}
			if k1Abs.Cmp(v.m1.ToBig()) != 0 {
				t.Fatal("cross: scalar decomposing with failed m1", i)
			}
			if k2Abs.Cmp(v.m2.ToBig()) != 0 {
				t.Fatal("cross: scalar decomposing with failed m2", i)
			}
		}
	})
}
