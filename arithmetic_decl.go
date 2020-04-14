// +build amd64,!generic

package bls

import (
	"golang.org/x/sys/cpu"
)

var mul func(c, a, b *fe) = mulADX
var mulAssign func(a, b *fe) = mulAssignADX

func cfgArch() {
	if !x86ArchitectureSet {
		if !(cpu.X86.HasADX && cpu.X86.HasBMI2) || forceNonADXArch {
			mul = mulNoADX
			mulAssign = mulAssignNoADX
		}
		x86ArchitectureSet = true
	}
}

func square(c, a *fe) {
	mul(c, a, a)
}

func neg(c, a *fe) {
	if a.IsZero() {
		c.Set(a)
	} else {
		_neg(c, a)
	}
}

//go:noescape
func add(c, a, b *fe)

//go:noescape
func addAssign(a, b *fe)

//go:noescape
func ladd(c, a, b *fe)

//go:noescape
func laddAssign(a, b *fe)

//go:noescape
func double(c, a *fe)

//go:noescape
func doubleAssign(a *fe)

//go:noescape
func ldouble(c, a *fe)

//go:noescape
func sub(c, a, b *fe)

//go:noescape
func subAssign(a, b *fe)

//go:noescape
func lsubAssign(a, b *fe)

//go:noescape
func _neg(c, a *fe)

//go:noescape
func mulNoADX(c, a, b *fe)

//go:noescape
func mulAssignNoADX(a, b *fe)

//go:noescape
func mulADX(c, a, b *fe)

//go:noescape
func mulAssignADX(a, b *fe)
