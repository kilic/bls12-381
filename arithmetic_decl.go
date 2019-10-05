package bls

func add(c, a, b *Fe)

func addn(a, b *Fe) uint64

func sub(c, a, b *Fe)

func subn(a, b *Fe) uint64

func neg(c, a *Fe)

func double(c, a *Fe)

// func mul(c *[12]uint64, a, b *Fe)

// func square(c *[12]uint64, a *Fe)

// func mont(c *Fe, w *[12]uint64)

func montmul(c, a, b *Fe)

func montsquare(c, a *Fe)
