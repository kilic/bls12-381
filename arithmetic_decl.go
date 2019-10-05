package bls

//go:noescape
func add(c, a, b *Fe)

//go:noescape
func addn(a, b *Fe) uint64

//go:noescape
func sub(c, a, b *Fe)

//go:noescape
func subn(a, b *Fe) uint64

//go:noescape
func neg(c, a *Fe)

//go:noescape
func double(c, a *Fe)

//go:noescape
func montmul(c, a, b *Fe)
