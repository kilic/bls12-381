package bls

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"testing"
)

func TestFp(t *testing.T) {
	field := NewFp()
	zero := &Fe{0}
	one := &Fe{1}
	t.Run("Encoding & Decoding", func(t *testing.T) {
		t.Run("1", func(t *testing.T) {
			bytes := []byte{0}
			fe := &Fe{}
			fe.FromBytes(bytes)
			if !field.Equal(fe, zero) {
				t.Errorf("bad encoding\n")
			}
		})
		t.Run("2", func(t *testing.T) {
			in := []byte{254, 253}
			fe := &Fe{}
			fe.FromBytes(in)
			if bytes.Equal(in, fe.Bytes()) {
				t.Errorf("bad encoding\n")
			}
		})
		t.Run("3", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			b.FromBytes(a.Bytes())
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
		t.Run("4", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			if _, err := b.SetString(a.String()); err != nil {
				t.Errorf("bad encoding or decoding\n")
			}
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
		t.Run("5", func(t *testing.T) {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			b := &Fe{}
			b.SetBig(a.Big())
			if !field.Equal(a, b) {
				t.Errorf("bad encoding or decoding\n")
			}
		})
	})
	t.Run("Addition", func(t *testing.T) {
		var a, b, c, u, v *Fe
		for i := 0; i < n; i++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Add(u, a, b)
			field.Add(u, u, c)
			field.Add(v, b, c)
			field.Add(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Additive associativity does not hold")
			}
			field.Add(u, a, b)
			field.Add(v, b, a)
			if !field.Equal(u, v) {
				t.Fatalf("Additive commutativity does not hold")
			}
			field.Add(u, a, zero)
			if !field.Equal(u, a) {
				t.Fatalf("Additive identity does not hold")
			}
			field.Neg(u, a)
			field.Add(u, u, a)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Negation\na:%s", a.String())
			}
		}
	})
	t.Run("Doubling", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Double(u, a)
			field.Add(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad doubling\na: %s\nu: %s\nv: %s\n", a, u, v)
			}
		}
	})
	t.Run("Subtraction", func(t *testing.T) {
		var a, b, c, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Sub(u, a, c)
			field.Sub(u, u, b)
			field.Sub(v, a, b)
			field.Sub(v, v, c)
			if !field.Equal(u, v) {
				t.Fatalf("Additive associativity does not hold\na: %s\nb: %s\nc: %s\nu: %s\nv:%s\n", a, b, c, u, v)
			}
			field.Sub(u, a, zero)
			if !field.Equal(u, a) {
				t.Fatalf("Additive identity does not hold\na: %s\nu: %s\n", a, u)
			}
			field.Sub(u, a, b)
			field.Sub(v, b, a)
			field.Add(u, u, v)
			if !field.Equal(u, zero) {
				t.Fatalf("Additive commutativity does not hold\na: %s\nb: %s\nu: %s\nv: %s", a, b, u, v)
			}
			field.Sub(u, a, b)
			field.Sub(v, b, a)
			field.Neg(v, v)
			if !field.Equal(u, u) {
				t.Fatalf("Bad Negation\na:%s", a.String())
			}
		}
	})
	t.Run("Montgomerry", func(t *testing.T) {
		var a, b, c, u, v, w *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			w = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			b, _ = field.RandElement(&Fe{}, rand.Reader)
			c, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Mont(u, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Montgomerry encoding")
			}
			field.Demont(u, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad Montgomerry decoding")
			}
			field.Mont(u, one)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad Montgomerry encoding")
			}
			field.Demont(u, field.One())
			if !field.Equal(u, one) {
				t.Fatalf("Bad Montgomerry decoding")
			}
			field.Mul(u, a, zero)
			if !field.Equal(u, zero) {
				t.Fatalf("Bad zero element")
			}
			field.Mul(u, a, one)
			field.Mul(u, u, r2)
			if !field.Equal(u, a) {
				t.Fatalf("Multiplication identity does not hold")
			}
			field.Mul(u, r2, one)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Multiplication identity does not hold, expected to equal r1")
			}
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.Exp(u, a, p)
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal itself")
			}
			field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u, v *Fe
		for j := 0; j < n; j++ {
			u = &Fe{}
			v = &Fe{}
			a, _ = field.RandElement(&Fe{}, rand.Reader)
			field.InvMontUp(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
			field.Mont(u, a)
			field.InvMontDown(v, u)
			field.Mul(v, v, u)
			if !field.Equal(v, one) {
				t.Fatalf("Bad inversion, expected to equal 1")
			}
			p := new(big.Int).SetBytes(modulus.Bytes())
			field.Exp(u, a, p.Sub(p, big.NewInt(2)))
			field.InvMontUp(v, a)
			if !field.Equal(v, u) {
				t.Fatalf("Bad inversion 1")
			}
			field.InvEEA(u, a)
			field.Mul(u, u, a)
			field.Mul(u, u, r2)
			if !field.Equal(u, one) {
				t.Fatalf("Bad inversion 2")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &Fe{}
		if field.Sqrt(r, nonResidue1) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.RandElement(&Fe{}, rand.Reader)
			aa, rr, r := &Fe{}, &Fe{}, &Fe{}
			field.Square(aa, a)
			if !field.Sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.Square(rr, r)
			if !field.Equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestFp2(t *testing.T) {
	field := NewFp2(nil)
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 96)
		for i := 0; i < 96; i++ {
			in[i] = 1
		}
		fe := &Fe2{}
		if err := field.NewElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.ToBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			v = &Fe2{}
			w = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			b, _ = field.RandElement(&Fe2{}, rand.Reader)
			c, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			v = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			_ = v
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.Exp(u, a, p)
			// if !field.Equal(u, a) {
			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
			// }
			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.Equal(u, field.One()) {
			// 	t.Fatalf("Bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *Fe2
		for j := 0; j < n; j++ {
			u = &Fe2{}
			a, _ = field.RandElement(&Fe2{}, rand.Reader)
			field.Inverse(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("Sqrt", func(t *testing.T) {
		r := &Fe2{}
		if field.Sqrt(r, nonResidue2) {
			t.Fatalf("bad sqrt 1")
		}
		for j := 0; j < n; j++ {
			a, _ := field.RandElement(&Fe2{}, rand.Reader)
			aa, rr, r := &Fe2{}, &Fe2{}, &Fe2{}
			field.Square(aa, a)
			if !field.Sqrt(r, aa) {
				t.Fatalf("bad sqrt 2")
			}
			field.Square(rr, r)
			if !field.Equal(rr, aa) {
				t.Fatalf("bad sqrt 3")
			}
		}
	})
}

func TestFp6(t *testing.T) {
	field := NewFp6(nil)
	// zero := field.Zero()
	t.Run("Encoding & Decoding", func(t *testing.T) {
		in := make([]byte, 288)
		for i := 0; i < 288; i++ {
			in[i] = 1
		}
		fe := &Fe6{}
		if err := field.NewElementFromBytes(fe, in); err != nil {
			panic(err)
		}
		if !bytes.Equal(in, field.ToBytes(fe)) {
			t.Errorf("bad encoding\n")
		}
	})
	t.Run("Multiplication", func(t *testing.T) {
		var a, b, c, u, v, w *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			v = &Fe6{}
			w = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			c, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Mul(u, a, b)
			field.Mul(u, u, c)
			field.Mul(v, b, c)
			field.Mul(v, v, a)
			if !field.Equal(u, v) {
				t.Fatalf("Multiplicative associativity does not hold")
			}
			field.Add(u, a, b)
			field.Mul(u, c, u)
			field.Mul(w, a, c)
			field.Mul(v, b, c)
			field.Add(v, v, w)
			if !field.Equal(u, v) {
				t.Fatalf("Distributivity does not hold")
			}
			field.Square(u, a)
			field.Mul(v, a, a)
			if !field.Equal(u, v) {
				t.Fatalf("Bad squaring")
			}
		}
	})
	t.Run("Exponentiation", func(t *testing.T) {
		var a, u, v *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			v = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Exp(u, a, big.NewInt(0))
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad exponentiation, expected to equal r1")
			}
			_ = v
			field.Exp(u, a, big.NewInt(1))
			if !field.Equal(u, a) {
				t.Fatalf("Bad exponentiation, expected to equal a")
			}
			field.Mul(u, a, a)
			field.Mul(u, u, u)
			field.Mul(u, u, u)
			field.Exp(v, a, big.NewInt(8))
			if !field.Equal(u, v) {
				t.Fatalf("Bad exponentiation")
			}
			// p := new(big.Int).SetBytes(modulus.Bytes())
			// field.Exp(u, a, p)
			// if !field.Equal(u, a) {
			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
			// }
			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
			// if !field.Equal(u, field.One()) {
			// 	t.Fatalf("Bad exponentiation, expected to equal one")
			// }
		}
	})
	t.Run("Inversion", func(t *testing.T) {
		var a, u *Fe6
		for j := 0; j < n; j++ {
			u = &Fe6{}
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			field.Inverse(u, a)
			field.Mul(u, u, a)
			if !field.Equal(u, field.One()) {
				t.Fatalf("Bad inversion, expected to equal r1")
			}
		}
	})
	t.Run("MulBy01", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *Fe6
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			u, _ = field.RandElement(&Fe6{}, rand.Reader)
			fq2.Copy(&b[2], fq2.Zero())
			field.Mul(u, a, b)
			field.mulBy01(a, &b[0], &b[1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 01")
			}
		}
	})
	t.Run("MulBy1", func(t *testing.T) {
		fq2 := field.f
		var a, b, u *Fe6
		for j := 0; j < n; j++ {
			a, _ = field.RandElement(&Fe6{}, rand.Reader)
			b, _ = field.RandElement(&Fe6{}, rand.Reader)
			u, _ = field.RandElement(&Fe6{}, rand.Reader)
			fq2.Copy(&b[2], fq2.Zero())
			fq2.Copy(&b[0], fq2.Zero())
			field.Mul(u, a, b)
			field.mulBy1(a, &b[1])
			if !field.Equal(a, u) {
				t.Fatal("Bad mul by 1")
			}
		}
	})
}

// func TestFp12(t *testing.T) {
// 	field := NewFp12(nil)
// 	t.Run("Encoding & Decoding", func(t *testing.T) {
// 		in := make([]byte, 576)
// 		for i := 0; i < 288; i++ {
// 			in[i] = 1
// 		}
// 		fe := &Fe12{}
// 		if err := field.NewElementFromBytes(fe, in); err != nil {
// 			panic(err)
// 		}
// 		if !bytes.Equal(in, field.ToBytes(fe)) {
// 			t.Errorf("bad encoding\n")
// 		}
// 	})
// 	t.Run("Multiplication", func(t *testing.T) {
// 		var a, b, c, u, v, w *Fe12
// 		for j := 0; j < n; j++ {
// 			u = &Fe12{}
// 			v = &Fe12{}
// 			w = &Fe12{}
// 			a, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			b, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			c, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			field.Mul(u, a, b)
// 			field.Mul(u, u, c)
// 			field.Mul(v, b, c)
// 			field.Mul(v, v, a)
// 			if !field.Equal(u, v) {
// 				t.Fatalf("Multiplicative associativity does not hold")
// 			}
// 			field.Add(u, a, b)
// 			field.Mul(u, c, u)
// 			field.Mul(w, a, c)
// 			field.Mul(v, b, c)
// 			field.Add(v, v, w)
// 			if !field.Equal(u, v) {
// 				t.Fatalf("Distributivity does not hold")
// 			}
// 			field.Square(u, a)
// 			field.Mul(v, a, a)
// 			if !field.Equal(u, v) {
// 				t.Fatalf("Bad squaring")
// 			}
// 		}
// 	})
// 	t.Run("Exponentiation", func(t *testing.T) {
// 		var a, u, v *Fe12
// 		for j := 0; j < n; j++ {
// 			u = &Fe12{}
// 			v = &Fe12{}
// 			a, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			field.Exp(u, a, big.NewInt(0))
// 			if !field.Equal(u, field.One()) {
// 				t.Fatalf("Bad exponentiation, expected to equal r1")
// 			}
// 			_ = v
// 			field.Exp(u, a, big.NewInt(1))
// 			if !field.Equal(u, a) {
// 				t.Fatalf("Bad exponentiation, expected to equal a")
// 			}
// 			field.Mul(u, a, a)
// 			field.Mul(u, u, u)
// 			field.Mul(u, u, u)
// 			field.Exp(v, a, big.NewInt(8))
// 			if !field.Equal(u, v) {
// 				t.Fatalf("Bad exponentiation")
// 			}
// 			// p := new(big.Int).SetBytes(modulus.Bytes())
// 			// field.Exp(u, a, p)
// 			// if !field.Equal(u, a) {
// 			// 	t.Fatalf("Bad exponentiation, expected to equal itself")
// 			// }
// 			// field.Exp(u, a, p.Sub(p, big.NewInt(1)))
// 			// if !field.Equal(u, field.One()) {
// 			// 	t.Fatalf("Bad exponentiation, expected to equal one")
// 			// }
// 		}
// 	})
// 	t.Run("Inversion", func(t *testing.T) {
// 		var a, u *Fe12
// 		for j := 0; j < n; j++ {
// 			u = &Fe12{}
// 			a, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			field.Inverse(u, a)
// 			field.Mul(u, u, a)
// 			if !field.Equal(u, field.One()) {
// 				t.Fatalf("Bad inversion, expected to equal r1")
// 			}
// 		}
// 	})
// 	t.Run("MulBy014", func(t *testing.T) {
// 		fq2 := field.f.f
// 		var a, b, u *Fe12
// 		for j := 0; j < n; j++ {
// 			a, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			b, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			u, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			fq2.Copy(&b[0][2], fq2.Zero())
// 			fq2.Copy(&b[1][0], fq2.Zero())
// 			fq2.Copy(&b[1][2], fq2.Zero())
// 			field.Mul(u, a, b)
// 			field.MulBy014Assign(a, &b[0][0], &b[0][1], &b[1][1])
// 			if !field.Equal(a, u) {
// 				t.Fatal("Bad mul by 014")
// 			}
// 		}
// 	})
// 	t.Run("MulBy034", func(t *testing.T) {
// 		fq2 := field.f.f
// 		var a, b, u *Fe12
// 		for j := 0; j < n; j++ {
// 			a, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			b, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			u, _ = field.RandElement(&Fe12{}, rand.Reader)
// 			fq2.Copy(&b[0][1], fq2.Zero())
// 			fq2.Copy(&b[0][2], fq2.Zero())
// 			fq2.Copy(&b[1][2], fq2.Zero())
// 			field.Mul(u, a, b)
// 			field.MulBy034Assign(a, &b[0][0], &b[1][0], &b[1][1])
// 			if !field.Equal(a, u) {
// 				t.Fatal("Bad mul by 034")
// 			}
// 		}
// 	})
// }

func BenchmarkFp(t *testing.B) {
	var a, b, c Fe
	var x, y, z lfe
	var field = NewFp()
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	field.RandElement(&c, rand.Reader)
	mul(&x, &a, &b)
	mul(&y, &a, &c)
	t.Run("Addition6", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			add6(&c, &a, &b)
		}
	})
	t.Run("LazyAddition6", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			ladd6(&c, &a, &b)
		}
	})
	t.Run("LazyAddition12", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			ladd12(&z, &x, &y)
		}
	})
	t.Run("Subtraction6", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			sub6(&c, &a, &b)
		}
	})
	t.Run("LazySubtraction12", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			lsub12(&z, &x, &y)
		}
	})
	t.Run("mul", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			mul(&z, &a, &b)
		}
	})
	t.Run("mont", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			mont(&c, &z)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.InvMontUp(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Exp(&c, &a, e)
		}
	})
}

func BenchmarkFp2(t *testing.B) {
	var a, b, c Fe2
	var field = NewFp2(nil)
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Addition", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Add(&c, &a, &b)
		}
	})
	t.Run("Subtraction", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Sub(&c, &a, &b)
		}
	})
	t.Run("Doubling", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Double(&c, &a)
		}
	})
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Squaring", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Inversion", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Inverse(&c, &a)
		}
	})
	t.Run("Exponentiation", func(t *testing.B) {
		e := new(big.Int).SetBytes(modulus.Bytes())
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Exp(&c, &a, e)
		}
	})
}

func TestLazyAdd6(t *testing.T) {
	var a, b, c Fe
	field := NewFp()
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		ladd6(&c, &a, &b)
		ab, bb := a.Big(), b.Big()
		cc := new(big.Int).Add(ab, bb)
		expected := make([]byte, 48)
		copy(expected[48-len(cc.Bytes()):], cc.Bytes())
		have := c.Bytes()
		if !bytes.Equal(have, expected) {
			t.Fatalf("")
		}
	}
}

func TestLazyAdd12(t *testing.T) {
	// r := new(big.Int)
	pSq := modulus.Big()
	pSq.Mul(pSq, pSq)
	// r.SetBit(r, 384, 1)
	// c := new(big.Int).Mul(pBig, r)
	// field := NewFp()
	for i := 0; i < n; i++ {
		aBig, _ := rand.Int(rand.Reader, pSq)
		bBig, _ := rand.Int(rand.Reader, pSq)
		cBig := new(big.Int).Add(aBig, bBig)
		a, b, c := &lfe{}, &lfe{}, &lfe{}
		a.FromBytes(aBig.Bytes())
		b.FromBytes(bBig.Bytes())
		ladd12(c, a, b)
		expected := make([]byte, 96)
		copy(expected[96-len(cBig.Bytes()):], cBig.Bytes())
		if !bytes.Equal(c.Bytes(), expected) {
			t.Fatalf("")
		}
	}
}

func TestLazySub12(t *testing.T) {
	p := modulus.Big()
	pSq := new(big.Int)
	pSq.Mul(p, p)
	r := new(big.Int)
	r2 := new(big.Int)
	r.SetBit(r, 384, 1)
	r2.SetBit(r, 384*2, 1)
	bound := new(big.Int).Mul(p, r)
	// fmt.Printf("%x\n", bound)
	for i := 0; i < n; i++ {
		// aBig := big.NewInt(1)
		aBig, _ := rand.Int(rand.Reader, pSq)
		// bBig := big.NewInt(2
		bBig, _ := rand.Int(rand.Reader, pSq)
		cBig := new(big.Int)
		cBig = cBig.Sub(aBig, bBig)
		if cBig.Sign() == -1 {
			cBig.Sub(r2, cBig)
		}
		if cBig.Cmp(bound) == 1 {
			cBig.Sub(cBig, bound)
		}
		a, b, c := &lfe{}, &lfe{}, &lfe{}
		a.FromBytes(aBig.Bytes())
		b.FromBytes(bBig.Bytes())
		lsub12(c, a, b)
		expected := make([]byte, 96)
		copy(expected[96-len(cBig.Bytes()):], cBig.Bytes())
		if !bytes.Equal(c.Bytes(), expected) {
			fmt.Printf("%x\n", c.Bytes())
			fmt.Printf("%x\n", expected)
			t.Fatalf("")
		}
	}
}

func TestLazyMul(t *testing.T) {
	var a, b Fe
	c := &lfe{}
	field := NewFp()
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		mul(c, &a, &b)
		ab, bb := a.Big(), b.Big()
		cc := new(big.Int).Mul(ab, bb)
		expected := make([]byte, 96)
		copy(expected[96-len(cc.Bytes()):], cc.Bytes())
		if !bytes.Equal(c.Bytes(), expected) {
			t.Fatalf("")
		}
	}
}

func TestMontRed(t *testing.T) {
	field := NewFp()
	for i := 0; i < n; i++ {
		var a, c Fe
		var lc lfe
		field.RandElement(&a, rand.Reader)
		cBig := new(big.Int).Mul(a.Big(), r1.Big())
		lc.FromBytes(cBig.Bytes())
		mont(&c, &lc)
		if !c.Equals(&a) {
			t.Fatalf("")
		}
	}
}

func TestFp2LazyMul(t *testing.T) {
	var a, b, c, c2, c3 Fe2
	var lc lfe2
	var field = NewFp2(nil)
	_, _ = c2, c3
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.mul(&c, &a, &b)
		field.lmul(&lc, &a, &b)
		mont(&c2[0], &lc[0])
		mont(&c2[1], &lc[1])
		field.Mul(&c3, &a, &b)
		if !field.Equal(&c, &c2) {
			fmt.Println(c)
			fmt.Println(c2)
			fmt.Println(i)
			t.Fatalf("c2")
		}
		if !field.Equal(&c, &c3) {
			fmt.Println(c)
			fmt.Println(c3)
			fmt.Println(i)
			t.Fatalf("c3")
		}
	}
}

func TestFp2LazySq(t *testing.T) {
	var a, b, c, c2 Fe2
	var lc lfe2
	var field = NewFp2(nil)
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.square(&c, &a)
		field.lsquare(&lc, &a)
		mont(&c2[0], &lc[0])
		mont(&c2[1], &lc[1])
		if !field.Equal(&c, &c2) {
			fmt.Println(c)
			fmt.Println(c2)
			fmt.Println(i)
			t.Fatalf("c2")
		}
	}
}

// func TestFp6Mul01(t *testing.T) {
// 	var a, b, c, c2 Fe6
// 	var field = NewFp6(nil)
// 	for i := 0; i < n; i++ {
// 		field.RandElement(&a, rand.Reader)
// 		field.RandElement(&b, rand.Reader)
// 		field.f.Copy(&b[2], field.f.Zero())
// 		field.mul(&c, &a, &b)
// 		field.mulBy01(&c2, &a, &b[0], &b[1])
// 		if !field.Equal(&c, &c2) {
// 			fmt.Println(c)
// 			fmt.Println()
// 			fmt.Println(c2)
// 			fmt.Println(i)
// 		}
// 	}
// }

// func TestFp6Mul1(t *testing.T) {
// 	var a, b, c, c2 Fe6
// 	var field = NewFp6(nil)
// 	for i := 0; i < n; i++ {
// 		field.RandElement(&a, rand.Reader)
// 		field.RandElement(&b, rand.Reader)
// 		field.f.Copy(&b[2], field.f.Zero())
// 		field.f.Copy(&b[0], field.f.Zero())
// 		field.mul(&c, &a, &b)
// 		field.mulBy1(&c2, &a, &b[1])
// 		if !field.Equal(&c, &c2) {
// 			fmt.Println(c)
// 			fmt.Println()
// 			fmt.Println(c2)
// 			fmt.Println(i)
// 		}
// 	}
// }

// func TestInverse6(t *testing.T) {
// 	field := NewFp6(nil)
// 	var a, u, v Fe6
// 	field.RandElement(&a, rand.Reader)
// 	field.Inverse(&u, &a)
// 	field.inverse(&v, &a)
// 	if !field.Equal(&u, &v) {
// 		fmt.Println(u)
// 		fmt.Println()
// 		fmt.Println(v)
// 	}
// }

// t.Run("MulBy01", func(t *testing.T) {
// 	fq2 := field.f
// 	var a, b, u *Fe6
// 	for j := 0; j < n; j++ {
// 		a, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		b, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		u, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		fq2.Copy(&b[2], fq2.Zero())
// 		field.Mul(u, a, b)
// 		field.MulBy01(a, &b[0], &b[1])
// 		if !field.Equal(a, u) {
// 			t.Fatal("Bad mul by 01")
// 		}
// 	}
// })
// t.Run("MulBy1", func(t *testing.T) {
// 	fq2 := field.f
// 	var a, b, u *Fe6
// 	for j := 0; j < n; j++ {
// 		a, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		b, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		u, _ = field.RandElement(&Fe6{}, rand.Reader)
// 		fq2.Copy(&b[2], fq2.Zero())
// 		fq2.Copy(&b[0], fq2.Zero())
// 		field.Mul(u, a, b)
// 		field.MulBy1(a, &b[1])
// 		if !field.Equal(a, u) {
// 			t.Fatal("Bad mul by 1")
// 		}
// 	}

// func TestFp12LazyMul(t *testing.T) {
// 	var a, b, c, c2 Fe12
// 	var field = NewFp12(nil)
// 	for i := 0; i < n; i++ {
// 		field.RandElement(&a, rand.Reader)
// 		field.RandElement(&b, rand.Reader)
// 		field.Mul(&c, &a, &b)
// 		field.mul(&c2, &a, &b)
// 		if !field.Equal(&c, &c2) {
// 			fmt.Println(c)
// 			fmt.Println()
// 			fmt.Println(c2)
// 			fmt.Println(i)
// 			t.Fatalf(":(")
// 		}
// 	}
// }

func TestFp6LazyMul(t *testing.T) {
	var a, b, c, c2, c3 Fe6
	var lc lfe6
	var field = NewFp6(nil)
	_, _ = c2, c3
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.Mul(&c, &a, &b)
		field.lmul(&lc, &a, &b)
		field.mont(&c2, &lc)
		if !field.Equal(&c, &c2) {
			// if !c[0][0].Equals(&c2[0][0]) {
			fmt.Println(c)
			fmt.Println()
			fmt.Println(c2)
			fmt.Println(i)
			t.Fatalf(":(")
		}
	}
}

func TestFp6LazySquare(t *testing.T) {
	var a, c, c2 Fe6
	var field = NewFp6(nil)
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.Square(&c, &a)
		field.square(&c2, &a)
		if !field.Equal(&c, &c2) {
			// if !c[0][0].Equals(&c2[0][0]) {
			fmt.Println(c)
			fmt.Println()
			fmt.Println(c2)
			fmt.Println(i)
			t.Fatalf(":(")
		}
	}
}

func BenchmarkFp6X(t *testing.B) {
	var a, b, c Fe6
	var field = NewFp6(nil)
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Faster Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
}

func BenchmarkFp2X(t *testing.B) {
	var a, b, c Fe2
	var field = NewFp2(nil)
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Mul(&c, &a, &b)
		}
	})
	t.Run("Faster Multiplication", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.mul(&c, &a, &b)
		}
	})
	t.Run("Sq1", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Sq2", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.square(&c, &a)
		}
	})
}

func BenchmarkFp12X(t *testing.B) {
	var a, b, c Fe12
	var field = NewFp12(nil)
	field.RandElement(&a, rand.Reader)
	field.RandElement(&b, rand.Reader)
	t.Run("Sq", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.Square(&c, &a)
		}
	})
	t.Run("Faster", func(t *testing.B) {
		t.ResetTimer()
		for i := 0; i < t.N; i++ {
			field.cyclotomicSquare(&c, &a)
		}
	})
}

func TestMixed(t *testing.T) {

	var a, b, c, z, z2 Fe
	var lz, lc lfe
	var field = NewFp()
	for i := 0; i < n; i++ {
		field.RandElement(&a, rand.Reader)
		field.RandElement(&b, rand.Reader)
		field.RandElement(&c, rand.Reader)
		field.mul(&z, &a, &b)
		field.add(&z, &z, &c)

		field.lmul(&lz, &a, &b)
		field.copyMixed(&lc, &b)
		field.add12(&lz, &lz, &lc)
		field.mont(&z2, &lz)
		if field.Equal(&z, &z2) {
			fmt.Println(z)
			fmt.Println()
			fmt.Println(z2)
			t.Fatal(i)
		}
	}
}

func TestFp12X(t *testing.T) {
	var a, b, c Fe12
	var field = NewFp12(nil)
	field.RandElement(&a, rand.Reader)
	// field.RandElement(&b, rand.Reader)
	field.Copy(&b, &a)
	//field.Square(&c, &a)
	field.mul(&c, &a, &a)
	fmt.Println(c)
	fmt.Println()
	field.Square(&c, &a)
	fmt.Println(c)
	// if !field.Equal(&a, &b) {
	// 	fmt.Println(a)
	// 	fmt.Println()
	// 	fmt.Println(b)
	// }
}

// func BenchmarkArithmetic(t *testing.B) {
// 	var a, b, c Fe
// 	var field = NewFp()
// 	field.RandElement(&a, rand.Reader)
// 	field.RandElement(&b, rand.Reader)
// 	t.Run("1", func(t *testing.B) {
// 		t.ResetTimer()
// 		for i := 0; i < t.N; i++ {
// 			add6(&c, &a, &b)
// 		}
// 	})
// 	t.Run("2", func(t *testing.B) {
// 		t.ResetTimer()
// 		for i := 0; i < t.N; i++ {
// 			add6alt(&c, &a, &b)
// 		}
// 	})
// 	t.Run("3", func(t *testing.B) {
// 		t.ResetTimer()
// 		for i := 0; i < t.N; i++ {
// 			sub6(&c, &a, &b)
// 		}
// 	})
// 	t.Run("4", func(t *testing.B) {
// 		t.ResetTimer()
// 		for i := 0; i < t.N; i++ {
// 			sub6alt(&c, &a, &b)
// 		}
// 	})
// }

// func TestFp12CS(t *testing.T) {
// 	var a, c, c2 Fe12
// 	var field = NewFp12(nil)
// 	for i := 0; i < n; i++ {
// 		field.RandElement(&a, rand.Reader)
// 		field.CyclotomicSquare(&c, &a)
// 		field.cyclotomicSquare(&c2, &a)
// 		if !field.Equal(&c, &c2) {
// 			// if !c[0][0].Equals(&c2[0][0]) {
// 			fmt.Println(c)
// 			fmt.Println()
// 			fmt.Println(c2)
// 			fmt.Println()
// 			// fmt.Println(c3)
// 			fmt.Println(i)
// 			t.Fatalf(":(")
// 		}
// 	}
// }

func BenchmarkCS(t *testing.B) {
	var a, c Fe12
	var field = NewFp12(nil)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		field.cyclotomicSquare(&c, &a)
	}
}

// func BenchmarkCS2(t *testing.B) {
// 	var a, c Fe12
// 	var field = NewFp12(nil)
// 	t.ResetTimer()
// 	for i := 0; i < t.N; i++ {
// 		field.cyclotomicSquare(&c, &a)
// 	}
// }
