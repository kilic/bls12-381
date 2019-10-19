package bls

//go:noescape
func add6(c, a, b *Fe)

//go:noescape
func add6alt(c, a, b *Fe)

//go:noescape
func ladd6(c, a, b *Fe)

//go:noescape
func ladd12(c, a, b *lfe)

//go:noescape
func ladd12_opt2(c, a, b *lfe)

//go:noescape
func addn(a, b *Fe) uint64

//go:noescape
func double6(c, a *Fe)

//go:noescape
func ldouble6(c, a *Fe)

//go:noescape
func ldouble12(c, a *lfe)

//go:noescape
func ldouble12_opt2(c, a *lfe)

//go:noescape
func sub6(c, a, b *Fe)

//go:noescape
func lsub6(c, a, b *Fe)

//go:noescape
func sub6alt(c, a, b *Fe)

//go:noescape
func lsub12(c, a, b *lfe)

//go:noescape
func lsub12_opt2(c, a, b *lfe)

//go:noescape
func lsub12_opt1_h2(c, a, b *lfe)

//go:noescape
func lsub12_opt1_h1(c, a, b *lfe)

//go:noescape
func subn(a, b *Fe) uint64

//go:noescape
func neg(c, a *Fe)

//go:noescape
func mul(c *lfe, a, b *Fe)

//go:noescape
func mont(c *Fe, a *lfe)

//go:noescape
func montmul(c, a, b *Fe)
