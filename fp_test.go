package bls

import (
	"bytes"
	"crypto/rand"
	"math/big"
	"testing"
)

func TestFpOne(t *testing.T) {
	t.Run("Encoding & Decoding", func(t *testing.T) {
		field := newFp()
		zero := &fe{0}
		t.Run("1", func(t *testing.T) {
			in := make([]byte, 48)
			fe := &fe{}
			if err := field.newElementFromBytes(fe, in); err != nil {
				t.Fatal(err)
			}
			if !field.equal(fe, zero) {
				t.Fatalf("bad encoding\n")
			}
			if !bytes.Equal(in, field.toBytes(fe)) {
				t.Fatalf("bad encoding\n")
			}
		})
		t.Run("2", func(t *testing.T) {
			in := make([]byte, 48)
			copy(in, []byte{0x11, 0x12})
			fe := &fe{}
			if err := field.newElementFromBytes(fe, in); err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(in, field.toBytes(fe)) {
				t.Fatalf("bad encoding\n")
			}
		})
		t.Run("3", func(t *testing.T) {
			for i := 0; i < n; i++ {
				a, _ := field.randElement(&fe{}, rand.Reader)
				b := &fe{}
				if err := field.newElementFromBytes(b, field.toBytes(a)); err != nil {
					t.Fatal(err)
				}
				if !field.equal(a, b) {
					t.Fatalf("bad encoding or decoding\n")
				}
			}
		})
		t.Run("4", func(t *testing.T) {
			for i := 0; i < n; i++ {
				a, _ := field.randElement(&fe{}, rand.Reader)
				b, err := field.newElementFromString(field.toString(a))
				if err != nil {
					t.Fatal(err)
				}
				if !field.equal(a, b) {
					t.Fatalf("bad encoding or decoding\n")
				}
			}
		})
		t.Run("5", func(t *testing.T) {
			for i := 0; i < n; i++ {
				a, _ := field.randElement(&fe{}, rand.Reader)
				b, err := field.newElementFromBig(field.toBig(a))
				if err != nil {
					t.Fatal(err)
				}
				if !field.equal(a, b) {
					t.Fatalf("bad encoding or decoding\n")
				}
			}
		})
	})
	t.Run("Addition", func(t *testing.T) {
		field := newFp()
		zero := &fe{}
		var a, b, c, u, v *fe
		for i := 0; i < n; i++ {
			u = &fe{}
			v = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			b, _ = field.randElement(&fe{}, rand.Reader)
			c, _ = field.randElement(&fe{}, rand.Reader)
			field.add(u, a, b)
			field.add(u, u, c)
			field.add(v, b, c)
			field.add(v, v, a)
			if !field.equal(u, v) {
				t.Fatalf("additive associativity does not hold")
			}
			field.add(u, a, b)
			field.add(v, b, a)
			if !field.equal(u, v) {
				t.Fatalf("additive commutativity does not hold")
			}
			field.add(u, a, zero)
			if !field.equal(u, a) {
				t.Fatalf("additive identity does not hold")
			}
			field.add(u, zero, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad zero addition")
			}
			field.neg(u, a)
			field.add(u, u, a)
			if !field.equal(u, zero) {
				t.Fatalf("bad Negation")
			}
		}
	})
	t.Run("Doubling", func(t *testing.T) {
		field := newFp()
		zero := &fe{}
		var a, u, v *fe
		for j := 0; j < n; j++ {
			u = &fe{}
			v = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			field.double(u, a)
			field.add(v, a, a)
			if !field.equal(u, v) {
				t.Fatalf("bad doubling\na: %s\nu: %s\nv: %s\n", a, u, v)
			}
			field.double(u, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad zero addition")
			}
		}
	})
	t.Run("Subtraction", func(t *testing.T) {
		field := newFp()
		zero := &fe{}
		var a, b, c, u, v *fe
		for j := 0; j < n; j++ {
			u = &fe{}
			v = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			b, _ = field.randElement(&fe{}, rand.Reader)
			c, _ = field.randElement(&fe{}, rand.Reader)
			field.sub(u, a, c)
			field.sub(u, u, b)
			field.sub(v, a, b)
			field.sub(v, v, c)
			if !field.equal(u, v) {
				t.Fatalf("additive associativity does not hold\na: %s\nb: %s\nc: %s\nu: %s\nv:%s\n", a, b, c, u, v)
			}
			field.sub(u, a, zero)
			if !field.equal(u, a) {
				t.Fatalf("additive identity does not hold\na: %s\nu: %s\n", a, u)
			}
			field.sub(u, a, b)
			field.sub(v, b, a)
			field.add(u, u, v)
			if !field.equal(u, zero) {
				t.Fatalf("additive commutativity does not hold\na: %s\nb: %s\nu: %s\nv: %s", a, b, u, v)
			}
			field.sub(u, zero, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad zero subtraction")
			}
			field.sub(u, a, a)
			if !field.equal(u, zero) {
				t.Fatalf("bad subtraction")
			}
			field.sub(u, a, b)
			field.sub(v, b, a)
			field.neg(v, v)
			if !field.equal(u, u) {
				t.Fatalf("bad negation\na:%s", a.String())
			}
		}
	})
	t.Run("Montgomerry", func(t *testing.T) {
		field := newFp()
		zero := &fe{}
		one := &fe{1}
		var a, b, c, u, v, w *fe
		for j := 0; j < n; j++ {
			u = &fe{}
			v = &fe{}
			w = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			b, _ = field.randElement(&fe{}, rand.Reader)
			c, _ = field.randElement(&fe{}, rand.Reader)
			field.mont(u, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad montgomerry encoding")
			}
			field.demont(u, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad montgomerry decoding")
			}
			field.mont(u, one)
			if !field.equal(u, field.one()) {
				t.Fatalf("bad montgomerry encoding")
			}
			field.demont(u, field.one())
			if !field.equal(u, one) {
				t.Fatalf("bad montgomerry decoding")
			}
			field.mul(u, a, zero)
			if !field.equal(u, zero) {
				t.Fatalf("bad zero element")
			}
			field.mul(u, a, one)
			field.mul(u, u, r2)
			if !field.equal(u, a) {
				t.Fatalf("multiplication identity does not hold")
			}
			field.mul(u, r2, one)
			if !field.equal(u, field.one()) {
				t.Fatalf("multiplication identity does not hold, expected to equal r1")
			}
			field.mul(u, a, b)
			field.mul(u, u, c)
			field.mul(v, b, c)
			field.mul(v, v, a)
			if !field.equal(u, v) {
				t.Fatalf("multiplicative associativity does not hold")
			}
			field.add(u, a, b)
			field.mul(u, c, u)
			field.mul(w, a, c)
			field.mul(v, b, c)
			field.add(v, v, w)
			if !field.equal(u, v) {
				t.Fatalf("distributivity does not hold")
			}
			field.square(u, a)
			field.mul(v, a, a)
			if !field.equal(u, v) {
				t.Fatalf("bad squaring")
			}
			field.mul(u, a, b)
			field.copy(v, a)
			field.mulAssign(v, b)
			if !field.equal(u, v) {
				t.Fatalf("bad mul assign")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *fe
		field := newFp()
		for j := 0; j < n; j++ {
			u = &fe{}
			v = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			field.exp(u, a, big.NewInt(0))
			if !field.equal(u, field.one()) {
				t.Fatalf("bad exponentiation, expected to equal r1")
			}
			field.exp(u, a, big.NewInt(1))
			if !field.equal(u, a) {
				t.Fatalf("bad exponentiation, expected to equal a")
			}
			field.mul(u, a, a)
			field.mul(u, u, u)
			field.mul(u, u, u)
			field.exp(v, a, big.NewInt(8))
			if !field.equal(u, v) {
				t.Fatalf("bad exponentiation")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.exp(u, a, p)
			if !field.equal(u, a) {
				t.Fatalf("bad exponentiation, expected to equal itself")
			}
			field.exp(u, a, p.Sub(p, big.NewInt(1)))
			if !field.equal(u, field.one()) {
				t.Fatalf("bad exponentiation, expected to equal r1")
			}
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u, v *fe
		field := newFp()
		one := &fe{1}
		for j := 0; j < n; j++ {
			u = &fe{}
			v = &fe{}
			a, _ = field.randElement(&fe{}, rand.Reader)
			field.invMontUp(u, a)
			field.mul(u, u, a)
			if !field.equal(u, field.one()) {
				t.Fatalf("bad inversion, expected to equal r1")
			}
			field.mont(u, a)
			field.invMontDown(v, u)
			field.mul(v, v, u)
			if !field.equal(v, one) {
				t.Fatalf("bad inversion, expected to equal 1")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.exp(u, a, p.Sub(p, big.NewInt(2)))
			field.invMontUp(v, a)
			if !field.equal(v, u) {
				t.Fatalf("bad inversion 1")
			}
			field.invEEA(u, a)
			field.mul(u, u, a)
			field.mul(u, u, r2)
			if !field.equal(u, one) {
				t.Fatalf("bad inversion 2")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &fe{}
		field := newFp()
		if field.sqrt(r, nonResidue1) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.randElement(&fe{}, rand.Reader)
			aa, rr, r := &fe{}, &fe{}, &fe{}
			field.square(aa, a)
			if !field.sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.square(rr, r)
			if !field.equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestInv(t *testing.T) {
	var a, u, v *fe
	field := newFp()
	one := &fe{1}
	for j := 0; j < n; j++ {
		u = &fe{}
		v = &fe{}
		a, _ = field.randElement(&fe{}, rand.Reader)
		field.invMontUp(u, a)
		field.mul(u, u, a)
		if !field.equal(u, field.one()) {
			t.Fatalf("bad inversion, expected to equal r1")
		}
		field.mont(u, a)
		field.invMontDown(v, u)
		field.mul(v, v, u)
		if !field.equal(v, one) {
			t.Fatalf("bad inversion, expected to equal 1")
		}
		p := new(big.Int).SetBytes(modulus.Bytes())
		field.exp(u, a, p.Sub(p, big.NewInt(2)))
		field.invMontUp(v, a)
		if !field.equal(v, u) {
			t.Fatalf("bad inversion 1")
		}
		field.invEEA(u, a)
		field.mul(u, u, a)
		field.mul(u, u, r2)
		if !field.equal(u, one) {
			t.Fatalf("bad inversion 2")
		}
	}
}

func TestFpTwo(t *testing.T) {
	field := newFp2(nil)
	t.Run("encoding & decoding", func(t *testing.T) {
		in := make([]byte, 96)
		for i := 0; i < 96; i++ {
			in[i] = 1
		}
		fe := &fe2{}
		if err := field.newElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.toBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *fe2
		for j := 0; j < n; j++ {
			u = &fe2{}
			v = &fe2{}
			w = &fe2{}
			a, _ = field.randElement(&fe2{}, rand.Reader)
			b, _ = field.randElement(&fe2{}, rand.Reader)
			c, _ = field.randElement(&fe2{}, rand.Reader)
			field.mul(u, a, b)
			field.mul(u, u, c)
			field.mul(v, b, c)
			field.mul(v, v, a)
			if !field.equal(u, v) {
				t.Fatalf("multiplicative associativity does not hold")
			}
			field.add(u, a, b)
			field.mul(u, c, u)
			field.mul(w, a, c)
			field.mul(v, b, c)
			field.add(v, v, w)
			if !field.equal(u, v) {
				t.Fatalf("distributivity does not hold")
			}
			field.square(u, a)
			field.mul(v, a, a)
			if !field.equal(u, v) {
				t.Fatalf("bad squaring")
			}
			field.mul(u, a, b)
			field.copy(v, a)
			field.mulAssign(v, b)
			if !field.equal(u, v) {
				t.Fatalf("bad mul assign")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *fe2
		for j := 0; j < n; j++ {
			u = &fe2{}
			v = &fe2{}
			a, _ = field.randElement(&fe2{}, rand.Reader)
			field.exp(u, a, big.NewInt(0))
			if !field.equal(u, field.one()) {
				t.Fatalf("bad exponentiation, expected to equal r1")
			}
			_ = v
			field.exp(u, a, big.NewInt(1))
			if !field.equal(u, a) {
				t.Fatalf("bad exponentiation, expected to equal a")
			}
			field.mul(u, a, a)
			field.mul(u, u, u)
			field.mul(u, u, u)
			field.exp(v, a, big.NewInt(8))
			if !field.equal(u, v) {
				t.Fatalf("bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.exp(u, a, p)
			// if !field.equal(u, a) {
			// 	t.Fatalf("bad exponentiation, expected to equal itself")
			// }
			// field.exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.equal(u, field.one()) {
			// 	t.Fatalf("bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *fe2
		for j := 0; j < n; j++ {
			u = &fe2{}
			a, _ = field.randElement(&fe2{}, rand.Reader)
			field.inverse(u, a)
			field.mul(u, u, a)
			if !field.equal(u, field.one()) {
				t.Fatalf("bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &fe2{}
		if field.sqrt(r, nonResidue2) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.randElement(&fe2{}, rand.Reader)
			aa, rr, r := &fe2{}, &fe2{}, &fe2{}
			field.square(aa, a)
			if !field.sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.square(rr, r)
			if !field.equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestFpSix(t *testing.T) {
	field := newFp6(nil)
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 288)
		for i := 0; i < 288; i++ {
			in[i] = 1
		}
		fe := &fe6{}
		if err := field.newElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.toBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *fe6
		for j := 0; j < n; j++ {
			u = &fe6{}
			v = &fe6{}
			w = &fe6{}
			a, _ = field.randElement(&fe6{}, rand.Reader)
			b, _ = field.randElement(&fe6{}, rand.Reader)
			c, _ = field.randElement(&fe6{}, rand.Reader)
			field.mul(u, a, b)
			field.mul(u, u, c)
			field.mul(v, b, c)
			field.mul(v, v, a)
			if !field.equal(u, v) {
				t.Fatalf("multiplicative associativity does not hold")
			}
			field.add(u, a, b)
			field.mul(u, c, u)
			field.mul(w, a, c)
			field.mul(v, b, c)
			field.add(v, v, w)
			if !field.equal(u, v) {
				t.Fatalf("distributivity does not hold")
			}
			field.square(u, a)
			field.mul(v, a, a)
			if !field.equal(u, v) {
				t.Fatalf("bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *fe6
		for j := 0; j < n; j++ {
			u = &fe6{}
			v = &fe6{}
			a, _ = field.randElement(&fe6{}, rand.Reader)
			field.exp(u, a, big.NewInt(0))
			if !field.equal(u, field.one()) {
				t.Fatalf("bad exponentiation, expected to equal r1")
			}
			_ = v
			field.exp(u, a, big.NewInt(1))
			if !field.equal(u, a) {
				t.Fatalf("bad exponentiation, expected to equal a")
			}
			field.mul(u, a, a)
			field.mul(u, u, u)
			field.mul(u, u, u)
			field.exp(v, a, big.NewInt(8))
			if !field.equal(u, v) {
				t.Fatalf("bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.exp(u, a, p)
			// if !field.equal(u, a) {
			// 	t.Fatalf("bad exponentiation, expected to equal itself")
			// }
			// field.exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.equal(u, field.one()) {
			// 	t.Fatalf("bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *fe6
		for j := 0; j < n; j++ {
			u = &fe6{}
			a, _ = field.randElement(&fe6{}, rand.Reader)
			field.inverse(u, a)
			field.mul(u, u, a)
			if !field.equal(u, field.one()) {
				t.Fatalf("bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("MulBy01", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *fe6
		c := &fe6{}
		for j := 0; j < n; j++ {
			a, _ = field.randElement(&fe6{}, rand.Reader)
			b, _ = field.randElement(&fe6{}, rand.Reader)
			u, _ = field.randElement(&fe6{}, rand.Reader)
			fq2.copy(&b[2], fq2.zero())
			field.mul(u, a, b)
			field.mulBy01(a, a, &b[0], &b[1])
			if !field.equal(a, u) {
				t.Fatal("bad mul by 01")
			}
		}
		_ = c
	})
	t.Run("MulBy1", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *fe6
		for j := 0; j < n; j++ {
			a, _ = field.randElement(&fe6{}, rand.Reader)
			b, _ = field.randElement(&fe6{}, rand.Reader)
			u, _ = field.randElement(&fe6{}, rand.Reader)
			fq2.copy(&b[2], fq2.zero())
			fq2.copy(&b[0], fq2.zero())
			field.mul(u, a, b)
			field.mulBy1(a, a, &b[1])
			if !field.equal(a, u) {
				t.Fatal("bad mul by 1")
			}
		}
	})
}

func TestFpTwelve(t *testing.T) {
	t.Run("Encoding & Decoding", func(t *testing.T) {
		field := newFp12(nil)
		in := make([]byte, 576)
		for i := 0; i < 288; i++ {
			in[i] = 1
		}
		fe := &fe12{}
		if err := field.newElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.toBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *fe12
		field := newFp12(nil)
		for j := 0; j < n; j++ {
			u = &fe12{}
			v = &fe12{}
			w = &fe12{}
			a, _ = field.randElement(&fe12{}, rand.Reader)
			b, _ = field.randElement(&fe12{}, rand.Reader)
			c, _ = field.randElement(&fe12{}, rand.Reader)
			field.mul(u, a, b)
			field.mul(u, u, c)
			field.mul(v, b, c)
			field.mul(v, v, a)
			if !field.equal(u, v) {
				t.Fatalf("multiplicative associativity does not hold")
			}
			field.add(u, a, b)
			field.mul(u, c, u)
			field.mul(w, a, c)
			field.mul(v, b, c)
			field.add(v, v, w)
			if !field.equal(u, v) {
				t.Fatalf("distributivity does not hold")
			}
			field.square(u, a)
			field.mul(v, a, a)
			if !field.equal(u, v) {
				t.Fatalf("bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *fe12
		field := newFp12(nil)
		for j := 0; j < n; j++ {
			u = &fe12{}
			v = &fe12{}
			a, _ = field.randElement(&fe12{}, rand.Reader)
			field.exp(u, a, big.NewInt(0))
			if !field.equal(u, field.one()) {
				t.Fatalf("bad exponentiation, expected to equal r1")
			}
			_ = v
			field.exp(u, a, big.NewInt(1))
			if !field.equal(u, a) {
				t.Fatalf("bad exponentiation, expected to equal a")
			}
			field.mul(u, a, a)
			field.mul(u, u, u)
			field.mul(u, u, u)
			field.exp(v, a, big.NewInt(8))
			if !field.equal(u, v) {
				t.Fatalf("bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.exp(u, a, p)
			// if !field.equal(u, a) {
			// 	t.Fatalf("bad exponentiation, expected to equal itself")
			// }
			// field.exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.equal(u, field.one()) {
			// 	t.Fatalf("bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		field := newFp12(nil)
		var a, u *fe12
		for j := 0; j < n; j++ {
			u = &fe12{}
			a, _ = field.randElement(&fe12{}, rand.Reader)
			field.inverse(u, a)
			field.mul(u, u, a)
			if !field.equal(u, field.one()) {
				t.Fatalf("bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("MulBy014", func(t *testing.T) {
		field := newFp12(nil)
		fq2 := field.f.f
		var a, b, u *fe12
		for j := 0; j < n; j++ {
			a, _ = field.randElement(&fe12{}, rand.Reader)
			b, _ = field.randElement(&fe12{}, rand.Reader)
			u, _ = field.randElement(&fe12{}, rand.Reader)
			fq2.copy(&b[0][2], fq2.zero())
			fq2.copy(&b[1][0], fq2.zero())
			fq2.copy(&b[1][2], fq2.zero())
			field.mul(u, a, b)
			field.mulBy014Assign(a, &b[0][0], &b[0][1], &b[1][1])
			if !field.equal(a, u) {
				t.Fatal("bad mul by 014")
			}
		}
	})
}

func BenchmarkFp1(t *testing.B) {
	var a, b, c fe
	var field = newFp()
	field.randElement(&a, rand.Reader)
	field.randElement(&b, rand.Reader)
	field.randElement(&c, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.exp(&c, &a, e)
		}
	})
	t.Run("Copy", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.copy(&c, &a)
		}
	})
}

func BenchmarkFp2(t *testing.B) {
	var a, b, c fe2
	var field = newFp2(nil)
	field.randElement(&a, rand.Reader)
	field.randElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.exp(&c, &a, e)
		}
	})
	t.Run("Copy", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.copy(&c, &a)
		}
	})
}

func BenchmarkFp6(t *testing.B) {
	var a, b, c fe6
	var field = newFp6(nil)
	field.randElement(&a, rand.Reader)
	field.randElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.exp(&c, &a, e)
		}
	})
	t.Run("Copy", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.copy(&c, &a)
		}
	})
}

func BenchmarkFp12(t *testing.B) {
	var a, b, c fe12
	var field = newFp12(nil)
	field.randElement(&a, rand.Reader)
	field.randElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.square(&c, &a)
		}
	})
	t.Run("Cyclotomic Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.cyclotomicSquare(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.exp(&c, &a, e)
		}
	})
	t.Run("Cyclotomic Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.cyclotomicExp(&c, &a, e)
		}
	})
	t.Run("Copy", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.copy(&c, &a)
		}
	})
}

func BenchmarkXXX(t *testing.B) {
	var a, b, c fe12
	var field = newFp12(nil)
	field.randElement(&a, rand.Reader)
	field.randElement(&b, rand.Reader)
	t.Run("1", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("2", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mulAssign(&a, &b)
		}
	})
	t.Run("1", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("1", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mulAssign(&a, &b)
		}
	})
}
