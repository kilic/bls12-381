// +build !nolazy

package bls

type fp12 struct {
	f   *fp6
	t   [5]*fe6
	t2  [9]*fe2
	lt  [4]*lfe6
	t12 *fe12
}

func newFp12(f *fp6) *fp12 {
	t := [5]*fe6{}
	t2 := [9]*fe2{}
	lt := [4]*lfe6{}
	for i := 0; i < len(t); i++ {
		t[i] = &fe6{}
	}
	for i := 0; i < len(t2); i++ {
		t2[i] = &fe2{}
	}
	for i := 0; i < len(lt); i++ {
		lt[i] = &lfe6{}
	}
	if f == nil {
		return &fp12{newFp6(nil), t, t2, lt, &fe12{}}
	}
	return &fp12{f, t, t2, lt, &fe12{}}
}

func (fp *fp12) mul(c, a, b *fe12) {
	t := fp.t
	lt := fp.lt
	fp.f.lmul(lt[0], &a[0], &b[0])
	fp.f.lmul(lt[1], &a[1], &b[1])
	fp.f.add(t[0], &a[0], &a[1])
	fp.f.add(t[1], &b[0], &b[1])
	fp.f.lmul(lt[2], t[0], t[1])
	fp.f.add12(lt[3], lt[0], lt[1])
	fp.f.subAssign12(lt[2], lt[3])
	fp.f.reduce(&c[1], lt[2])
	fp.f.mulByNonResidue12(lt[2], lt[1])
	fp.f.addAssign12(lt[2], lt[0])
	fp.f.reduce(&c[0], lt[2])
}

func (fp *fp12) mulAssign(a, b *fe12) {
	t := fp.t
	lt := fp.lt
	fp.f.lmul(lt[0], &a[0], &b[0])
	fp.f.lmul(lt[1], &a[1], &b[1])
	fp.f.add(t[0], &a[0], &a[1])
	fp.f.add(t[1], &b[0], &b[1])
	fp.f.lmul(lt[2], t[0], t[1])
	fp.f.add12(lt[3], lt[0], lt[1])
	fp.f.subAssign12(lt[2], lt[3])
	fp.f.reduce(&a[1], lt[2])
	fp.f.mulByNonResidue12(lt[2], lt[1])
	fp.f.addAssign12(lt[2], lt[0])
	fp.f.reduce(&a[0], lt[2])
}
