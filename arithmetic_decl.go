package bls

//go:noescape
func add6(c, a, b *fe)

//go:noescape
func add_assign_6(a, b *fe)

//go:noescape
func ladd6(c, a, b *fe)

//go:noescape
func ladd_assign_6(a, b *fe)

//go:noescape
func double6(c, a *fe)

//go:noescape
func double_assign_6(a *fe)

//go:noescape
func ldouble6(c, a *fe)

//go:noescape
func sub6(c, a, b *fe)

//go:noescape
func sub_assign_6(a, b *fe)

//go:noescape
func lsub6(c, a, b *fe)

//go:noescape
func lsub_assign_nc_6(a, b *fe)

//go:noescape
func neg(c, a *fe)

//go:noescape
func montmul_nobmi2(c, a, b *fe)

//go:noescape
func montmul_assign_nobmi2(a, b *fe)

//go:noescape
func montmul_bmi2(c, a, b *fe)

//go:noescape
func montmul_assign_bmi2(a, b *fe)
