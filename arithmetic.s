#include "textflag.h"

// func add(c *Fe, a *Fe, b *Fe)
TEXT ·add(SB), NOSPLIT, $0-24
	// |
	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI
	XORQ AX, AX

	// |
	MOVQ (DI), R8
	ADDQ (SI), R8
	MOVQ 8(DI), R9
	ADCQ 8(SI), R9
	MOVQ 16(DI), R10
	ADCQ 16(SI), R10
	MOVQ 24(DI), R11
	ADCQ 24(SI), R11
	MOVQ 32(DI), R12
	ADCQ 32(SI), R12
	MOVQ 40(DI), R13
	ADCQ 40(SI), R13

	// |
	MOVQ R8, R14
	SUBQ ·modulus+0(SB), R14
	MOVQ R9, R15
	SBBQ ·modulus+8(SB), R15
	MOVQ R10, CX
	SBBQ ·modulus+16(SB), CX
	MOVQ R11, DX
	SBBQ ·modulus+24(SB), DX
	MOVQ R12, SI
	SBBQ ·modulus+32(SB), SI
	MOVQ R13, BX
	SBBQ ·modulus+40(SB), BX

	// |
	MOVQ    c+0(FP), DI
	CMOVQCC R14, R8
	MOVQ    R8, (DI)
	CMOVQCC R15, R9
	MOVQ    R9, 8(DI)
	CMOVQCC CX, R10
	MOVQ    R10, 16(DI)
	CMOVQCC DX, R11
	MOVQ    R11, 24(DI)
	CMOVQCC SI, R12
	MOVQ    R12, 32(DI)
	CMOVQCC BX, R13
	MOVQ    R13, 40(DI)
	RET

// func addn(a *Fe, b *Fe) uint64
TEXT ·addn(SB), NOSPLIT, $0-24
	// |
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	XORQ AX, AX

	// |
	MOVQ (DI), R8
	ADDQ (SI), R8
	MOVQ 8(DI), R9
	ADCQ 8(SI), R9
	MOVQ 16(DI), R10
	ADCQ 16(SI), R10
	MOVQ 24(DI), R11
	ADCQ 24(SI), R11
	MOVQ 32(DI), R12
	ADCQ 32(SI), R12
	MOVQ 40(DI), R13
	ADCQ 40(SI), R13
	ADCQ $0, AX

	// |
	MOVQ R8, (DI)
	MOVQ R9, 8(DI)
	MOVQ R10, 16(DI)
	MOVQ R11, 24(DI)
	MOVQ R12, 32(DI)
	MOVQ R13, 40(DI)
	MOVQ AX, ret+16(FP)
	RET

// func sub(c *Fe, a *Fe, b *Fe)
TEXT ·sub(SB), NOSPLIT, $0-24
	// |
	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI
	XORQ AX, AX
	MOVQ (DI), R8
	SUBQ (SI), R8
	MOVQ 8(DI), R9
	SBBQ 8(SI), R9
	MOVQ 16(DI), R10
	SBBQ 16(SI), R10
	MOVQ 24(DI), R11
	SBBQ 24(SI), R11
	MOVQ 32(DI), R12
	SBBQ 32(SI), R12
	MOVQ 40(DI), R13
	SBBQ 40(SI), R13

	// |
	MOVQ    ·modulus+0(SB), R14
	CMOVQCC AX, R14
	MOVQ    ·modulus+8(SB), R15
	CMOVQCC AX, R15
	MOVQ    ·modulus+16(SB), CX
	CMOVQCC AX, CX
	MOVQ    ·modulus+24(SB), DX
	CMOVQCC AX, DX
	MOVQ    ·modulus+32(SB), SI
	CMOVQCC AX, SI
	CMOVQCS ·modulus+40(SB), AX
	MOVQ    AX, BX

	// |
	MOVQ c+0(FP), DI
	ADDQ R14, R8
	MOVQ R8, (DI)
	ADCQ R15, R9
	MOVQ R9, 8(DI)
	ADCQ CX, R10
	MOVQ R10, 16(DI)
	ADCQ DX, R11
	MOVQ R11, 24(DI)
	ADCQ SI, R12
	MOVQ R12, 32(DI)
	ADCQ BX, R13
	MOVQ R13, 40(DI)
	RET

// func subn(a *Fe, b *Fe) uint64
TEXT ·subn(SB), NOSPLIT, $0-24
	// |
	MOVQ a+0(FP), DI
	MOVQ b+8(FP), SI
	XORQ AX, AX

	// |
	MOVQ (DI), R8
	SUBQ (SI), R8
	MOVQ 8(DI), R9
	SBBQ 8(SI), R9
	MOVQ 16(DI), R10
	SBBQ 16(SI), R10
	MOVQ 24(DI), R11
	SBBQ 24(SI), R11
	MOVQ 32(DI), R12
	SBBQ 32(SI), R12
	MOVQ 40(DI), R13
	SBBQ 40(SI), R13
	ADCQ $0, AX

	// |
	MOVQ R8, (DI)
	MOVQ R9, 8(DI)
	MOVQ R10, 16(DI)
	MOVQ R11, 24(DI)
	MOVQ R12, 32(DI)
	MOVQ R13, 40(DI)
	MOVQ AX, ret+16(FP)
	RET

// func double(c *Fe, a *Fe)
TEXT ·double(SB), NOSPLIT, $8-16
	// |
	MOVQ a+8(FP), DI
	XORQ AX, AX
	MOVQ (DI), R8
	ADDQ R8, R8
	MOVQ 8(DI), R9
	ADCQ R9, R9
	MOVQ 16(DI), R10
	ADCQ R10, R10
	MOVQ 24(DI), R11
	ADCQ R11, R11
	MOVQ 32(DI), R12
	ADCQ R12, R12
	MOVQ 40(DI), R13
	ADCQ R13, R13

	// |
	MOVQ R8, R14
	SUBQ ·modulus+0(SB), R14
	MOVQ R9, R15
	SBBQ ·modulus+8(SB), R15
	MOVQ R10, CX
	SBBQ ·modulus+16(SB), CX
	MOVQ R11, DX
	SBBQ ·modulus+24(SB), DX
	MOVQ R12, SI
	SBBQ ·modulus+32(SB), SI
	MOVQ R13, BX
	SBBQ ·modulus+40(SB), BX

	// |
	MOVQ    c+0(FP), DI
	CMOVQCC R14, R8
	MOVQ    R8, (DI)
	CMOVQCC R15, R9
	MOVQ    R9, 8(DI)
	CMOVQCC CX, R10
	MOVQ    R10, 16(DI)
	CMOVQCC DX, R11
	MOVQ    R11, 24(DI)
	CMOVQCC SI, R12
	MOVQ    R12, 32(DI)
	CMOVQCC BX, R13
	MOVQ    R13, 40(DI)
	RET

// func neg(c *Fe, a *Fe)
TEXT ·neg(SB), NOSPLIT, $0-16
	// |
	MOVQ a+8(FP), DI

	// |
	MOVQ ·modulus+0(SB), R8
	SUBQ (DI), R8
	MOVQ ·modulus+8(SB), R9
	SBBQ 8(DI), R9
	MOVQ ·modulus+16(SB), R10
	SBBQ 16(DI), R10
	MOVQ ·modulus+24(SB), R11
	SBBQ 24(DI), R11
	MOVQ ·modulus+32(SB), R12
	SBBQ 32(DI), R12
	MOVQ ·modulus+40(SB), R13
	SBBQ 40(DI), R13

	// |
	MOVQ c+0(FP), DI
	MOVQ R8, (DI)
	MOVQ R9, 8(DI)
	MOVQ R10, 16(DI)
	MOVQ R11, 24(DI)
	MOVQ R12, 32(DI)
	MOVQ R13, 40(DI)
	RET


// func montmul(c *Fe, a *Fe, b *Fe)
TEXT ·montmul(SB), NOSPLIT, $56-24
	// |
	// | Multiplication
	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI

	// |
	// |
	XORQ R10, R10
	XORQ R11, R11
	XORQ R12, R12
	XORQ R13, R13
	XORQ R14, R14
	XORQ R15, R15
	MOVQ $0, (SP)
	MOVQ $0, 8(SP)
	MOVQ $0, 16(SP)
	MOVQ $0, 24(SP)

	// |
	// | b0
	MOVQ (SI), CX

	// | a0 * b0
	// | (w0, w1) @ (R8, R9)
	MOVQ (DI), AX
	MULQ CX
	MOVQ AX, R8
	MOVQ DX, R9

	// | a1 * b0
	// | (w1, w2) @ (R9, R10)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10

	// | a2 * b0
	// | (w2, w3) @ (R10, R11)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11

	// | a3 * b0
	// | (w3, w4) @ (R11, R12)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12

	// | a4 * b0
	// | (w4, w5) @ (R12, R13)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13

	// | a5 * b0
	// | (w5, w6) @ (R13, R14)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14

	// |
	// | b1
	MOVQ 8(SI), CX

	// | a0 * b1
	// | (w1, w2, w3, w4) @ (R9, R10, R11, R12)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ $0, R11
	ADCQ $0, R12

	// | a1 * b1
	// | (w2, w3, w4, w5) @ (R10, R11, R12, R13)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ $0, R12
	ADCQ $0, R13

	// | a2 * b1
	// | (w3, w4, w5, w6) @ (R11, R12, R13, R14)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
	ADCQ $0, R13
	ADCQ $0, R14

	// | a3 * b1
	// | (w4, w5, w6, w7) @ (R12, R13, R14, R15)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14
	ADCQ $0, R15

	// | a4 * b1
	// | (w5, w6, w7, w8) @ (R13, R14, R15, (SP))
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, (SP)

	// | a5 * b1
	// | (w6, w7, w8, w9) @ (R14, R15, (SP), 8(SP))
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, (SP)
	ADCQ $0, 8(SP)

	// |
	// | b2
	MOVQ 16(SI), CX

	// | a0 * b2
	// | (w2, w3, w4, w5) @ (R10, R11, R12, R13)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ $0, R12
	ADCQ $0, R13

	// | a1 * b2
	// | (w3, w4, w5, w6) @ (R11, R12, R13, R14)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
	ADCQ $0, R13
	ADCQ $0, R14

	// | a2 * b2
	// | (w4, w5, w6, w7) @ (R12, R13, R14, R15)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14
	ADCQ $0, R15

	// | a3 * b2
	// | (w5, w6, w7, w8) @ (R13, R14, R15, (SP))
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, (SP)

	// | a4 * b2
	// | (w6, w7, w8, w9) @ (R14, R15, (SP), 8(SP))
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, (SP)
	ADCQ $0, 8(SP)

	// | a5 * b2
	// | (w7, w8, w9, w10) @ (R15, (SP), 8(SP), 16(SP))
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R15
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)
	ADCQ $0, 16(SP)

	// |
	// | b3
	MOVQ 24(SI), CX

	// | a0 * b3
	// | (w3, w4, w5, w6) @ (R11, R12, R13, R14)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
	ADCQ $0, R13
	ADCQ $0, R14

	// | a1 * b3
	// | (w4, w5, w6, w7) @ (R12, R13, R14, R15)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14
	ADCQ $0, R15

	// | a2 * b3
	// | (w5, w6, w7, w8) @ (R13, R14, R15, (SP))
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, (SP)

	// | a3 * b3
	// | (w6, w7, w8, w9) @ (R14, R15, (SP), 8(SP))
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, (SP)
	ADCQ $0, 8(SP)

	// | a4 * b3
	// | (w7, w8, w9, w10) @ (R15, (SP), 8(SP), 16(SP))
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R15
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)
	ADCQ $0, 16(SP)

	// | a5 * b3
	// | (w8, w9, w10, w11) @ ((SP), 8(SP), 16(SP), 24(SP))
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, (SP)
	ADCQ DX, 8(SP)
	ADCQ $0, 16(SP)
	ADCQ $0, 24(SP)

	// |
	// | b4
	MOVQ 32(SI), CX

	// | a0 * b4
	// | (w4, w5, w6, w7) @ (R12, R13, R14, R15)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14
	ADCQ $0, R15

	// | a1 * b4
	// | (w5, w6, w7, w8) @ (R13, R14, R15, (SP))
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, (SP)

	// | a2 * b4
	// | (w6, w7, w8, w9) @ (R14, R15, (SP), 8(SP))
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, (SP)
	ADCQ $0, 8(SP)

	// | a3 * b4
	// | (w7, w8, w9, w10) @ (R15, (SP), 8(SP), 16(SP))
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R15
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)
	ADCQ $0, 16(SP)

	// | a4 * b4
	// | (w8, w9, w10, w11) @ ((SP), 8(SP), 16(SP), 24(SP))
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, (SP)
	ADCQ DX, 8(SP)
	ADCQ $0, 16(SP)
	ADCQ $0, 24(SP)

	// | a5 * b4
	// | (w9, w10, w11) @ (8(SP), 16(SP), 24(SP))
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, 8(SP)
	ADCQ DX, 16(SP)
	ADCQ $0, 24(SP)

	// |
	// | b5
	MOVQ 40(SI), CX

	// | a0 * b5
	// | (w5, w6, w7, w8) @ (R13, R14, R15, (SP))
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, (SP)

	// | a1 * b5
	// | (w6, w7, w8, w9) @ (R14, R15, (SP), 8(SP))
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, (SP)
	ADCQ $0, 8(SP)

	// | a2 * b5
	// | (w7, w8, w9, w10) @ (R15, (SP), 8(SP), 16(SP))
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R15
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)
	ADCQ $0, 16(SP)

	// | a3 * b5
	// | (w8, w9, w10, w11) @ ((SP), 8(SP), 16(SP), 24(SP))
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, (SP)
	ADCQ DX, 8(SP)
	ADCQ $0, 16(SP)
	ADCQ $0, 24(SP)

	// | a4 * b5
	// | (w9, w10, w11) @ (8(SP), 16(SP), 24(SP))
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, 8(SP)
	ADCQ DX, 16(SP)
	ADCQ $0, 24(SP)

	// | a5 * b5
	// | (w10, w11) @ (16(SP), 24(SP))
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, 16(SP)
	ADCQ DX, 24(SP)

	// |
	// | Montgomerry Reduction
	MOVQ R15, 32(SP)

	// |
	// | (u @ CX) = (w0 @ R8) * inp
	MOVQ R8, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w0 @ R8
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, DI

	// | w1 @ R9
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ DI, R9
	ADCQ $0, SI

	// | w2 @ R10
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, DI
	ADDQ SI, R10
	ADCQ $0, DI

	// | w3 @ R11
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ DI, R11
	ADCQ $0, SI

	// | w4 @ R12
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, DI
	ADDQ SI, R12
	ADCQ $0, DI

	// | w5 @ R13
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, SI
	ADDQ DI, R13
	ADCQ $0, SI

	// | w6 @ R14
	ADDQ SI, R14
	MOVQ $0, R15
	ADCQ $0, R15

	// |
	MOVQ 32(SP), R8

	// | (u @ CX) = (w1 @ R9) * inp
	MOVQ R9, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w1 @ R9
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, DI

	// | w2 @ R10
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, SI
	ADDQ DI, R10
	ADCQ $0, SI

	// | w3 @ R11
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, DI
	ADDQ SI, R11
	ADCQ $0, DI

	// | w4 @ R12
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, SI
	ADDQ DI, R12
	ADCQ $0, SI

	// | w5 @ R13
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, DI
	ADDQ SI, R13
	ADCQ $0, DI

	// | w6 @ R14
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, SI
	ADDQ DI, R14
	ADCQ $0, SI

	// | w7 @ R8
	ADDQ SI, R15
	ADCQ R15, R8
	MOVQ $0, R15
	ADCQ $0, R15

	// |
	MOVQ (SP), R9

	// | (u @ CX) = (w2 @ R10) * inp
	MOVQ R10, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w2 @ R10
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, DI

	// | w3 @ R11
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ DI, R11
	ADCQ $0, SI

	// | w4 @ R12
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, DI
	ADDQ SI, R12
	ADCQ $0, DI

	// | w5 @ R13
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, SI
	ADDQ DI, R13
	ADCQ $0, SI

	// | w6 @ R14
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, DI
	ADDQ SI, R14
	ADCQ $0, DI

	// | w7 @ R8
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, SI
	ADDQ DI, R8
	ADCQ $0, SI

	// | w8 @ R9
	ADDQ SI, R15
	ADCQ R15, R9
	MOVQ $0, R15
	ADCQ $0, R15

	// |
	MOVQ 8(SP), R10

	// | (u @ CX) = (w3 @ R11) * inp
	MOVQ R11, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w3 @ R11
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, DI

	// | w4 @ R12
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, SI
	ADDQ DI, R12
	ADCQ $0, SI

	// | w5 @ R13
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, DI
	ADDQ SI, R13
	ADCQ $0, DI

	// | w6 @ R14
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, SI
	ADDQ DI, R14
	ADCQ $0, SI

	// | w7 @ R8
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, DI
	ADDQ SI, R8
	ADCQ $0, DI

	// | w8 @ R9
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ DI, R9
	ADCQ $0, SI

	// | w9 @ R10
	ADDQ SI, R15
	ADCQ R15, R10
	MOVQ $0, R15
	ADCQ $0, R15

	// |
	MOVQ 16(SP), R11

	// | (u @ CX) = (w4 @ R12) * inp
	MOVQ R12, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w4 @ R12
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, DI

	// | w5 @ R13
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, SI
	ADDQ DI, R13
	ADCQ $0, SI

	// | w6 @ R14
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, DI
	ADDQ SI, R14
	ADCQ $0, DI

	// | w7 @ R8
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, SI
	ADDQ DI, R8
	ADCQ $0, SI

	// | w8 @ R9
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, DI
	ADDQ SI, R9
	ADCQ $0, DI

	// | w9 @ R10
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, SI
	ADDQ DI, R10
	ADCQ $0, SI

	// | w10 @ R11
	ADDQ SI, R15
	ADCQ R15, R11
	MOVQ $0, R15
	ADCQ $0, R15

	// |
	MOVQ 24(SP), R12

	// | (u @ CX) = (w5 @ R13) * inp
	MOVQ R13, AX
	MULQ ·inp+0(SB)
	MOVQ AX, CX

	// | w5 @ R13
	XORQ DI, DI
	MOVQ ·modulus+0(SB), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, DI

	// | w6 @ R14
	XORQ SI, SI
	MOVQ ·modulus+8(SB), AX
	MULQ CX
	ADDQ AX, R14
	ADCQ DX, SI
	ADDQ DI, R14
	ADCQ $0, SI

	// | w7 @ R8
	XORQ DI, DI
	MOVQ ·modulus+16(SB), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, DI
	ADDQ SI, R8
	ADCQ $0, DI

	// | w8 @ R9
	XORQ SI, SI
	MOVQ ·modulus+24(SB), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ DI, R9
	ADCQ $0, SI

	// | w9 @ R10
	XORQ DI, DI
	MOVQ ·modulus+32(SB), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, DI
	ADDQ SI, R10
	ADCQ $0, DI

	// | w10 @ R11
	XORQ SI, SI
	MOVQ ·modulus+40(SB), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ DI, R11
	ADCQ $0, SI

	// | w11 @ R12
	ADDQ SI, R15
	ADCQ R15, R12
	MOVQ $0, R15
	ADCQ $0, R15

	// | Reduce by modulus
	MOVQ R14, R13
	SUBQ ·modulus+0(SB), R13
	MOVQ R8, CX
	SBBQ ·modulus+8(SB), CX
	MOVQ R9, AX
	SBBQ ·modulus+16(SB), AX
	MOVQ R10, DX
	SBBQ ·modulus+24(SB), DX
	MOVQ R11, BX
	SBBQ ·modulus+32(SB), BX
	MOVQ BX, 40(SP)
	MOVQ R12, BX
	SBBQ ·modulus+40(SB), BX
	MOVQ BX, 48(SP)
	SBBQ $0, R15

	// | Compare & Return
	MOVQ    c+0(FP), DI
	CMOVQCC R13, R14
	MOVQ    R14, (DI)
	CMOVQCC CX, R8
	MOVQ    R8, 8(DI)
	CMOVQCC AX, R9
	MOVQ    R9, 16(DI)
	CMOVQCC DX, R10
	MOVQ    R10, 24(DI)
	CMOVQCC 40(SP), R11
	MOVQ    R11, 32(DI)
	CMOVQCC 48(SP), R12
	MOVQ    R12, 40(DI)
	RET




// func montsquare(c *Fe, a *Fe)
TEXT ·montsquare(SB), NOSPLIT, $40-16
	MOVQ a+8(FP), DI
	XORQ R11, R11
	XORQ R12, R12
	XORQ R13, R13
	XORQ R14, R14
	XORQ R15, R15
	XORQ CX, CX
	XORQ SI, SI
	MOVQ $0, (SP)
	MOVQ $0, 8(SP)
	MOVQ $0, 16(SP)

	// | a0
	// | w0 @ R9
	MOVQ (DI), R8
	MOVQ R8, AX
	MULQ R8
	MOVQ AX, R9
	MOVQ DX, R10

	// | w1 @ R10
	MOVQ 8(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R12
	ADDQ AX, R10
	ADCQ DX, R11

	// | w2 @ R11
	MOVQ 16(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R13
	ADDQ AX, R11
	ADCQ DX, R12
	ADCQ $0, R13

	// | w3 @ R12
	MOVQ 24(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R14
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14

	// | w4 @ R13
	MOVQ 32(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R15
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15

	// | w5 @ R14
	MOVQ 40(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, CX

	// | a1
	// | w2 @ R11
	MOVQ 8(DI), R8
	MOVQ R8, AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, R12
	ADCQ $0, R13
	ADCQ $0, R14

	// | w3 @ R12
	MOVQ 16(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R14
	ADDQ AX, R12
	ADCQ DX, R13
	ADCQ $0, R14

	// | w4 @ R13
	MOVQ 24(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, R15
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15

	// | w5 @ R14
	MOVQ 32(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, CX

	// | w6 @ R15
	MOVQ 40(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, SI
	ADDQ AX, R15
	ADCQ DX, CX
	ADCQ $0, SI

	// | a2
	// | w4 @ R13
	MOVQ 16(DI), R8
	MOVQ R8, AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, R14
	ADCQ $0, R15
	ADCQ $0, CX

	// | w5 @ R14
	MOVQ 24(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, CX
	ADDQ AX, R14
	ADCQ DX, R15
	ADCQ $0, CX

	// | w6 @ R15
	MOVQ 32(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, SI
	ADDQ AX, R15
	ADCQ DX, CX
	ADCQ $0, SI

	// | w7 @ CX
	MOVQ 40(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, (SP)
	ADDQ AX, CX
	ADCQ DX, SI
	ADCQ $0, (SP)

	// | a3
	// | w6 @ R15
	MOVQ 24(DI), R8
	MOVQ R8, AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, CX
	ADCQ $0, SI
	ADCQ $0, (SP)

	// | w7 @ CX
	MOVQ 32(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, (SP)
	ADDQ AX, CX
	ADCQ DX, SI
	ADCQ $0, (SP)

	// | w8 @ SI
	MOVQ 40(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, 8(SP)
	ADDQ AX, SI
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)

	// | a4
	// | w8 @ SI
	MOVQ 32(DI), R8
	MOVQ R8, AX
	MULQ R8
	ADDQ AX, SI
	ADCQ DX, (SP)
	ADCQ $0, 8(SP)

	// | w9 @ (SP)
	MOVQ 40(DI), AX
	MULQ R8
	ADDQ AX, AX
	ADCQ DX, DX
	ADCQ $0, 16(SP)
	ADDQ AX, (SP)
	ADCQ DX, 8(SP)
	ADCQ $0, 16(SP)

	// | a5
	// | w10 @ 8(SP)
	MOVQ 40(DI), R8
	MOVQ R8, AX
	MULQ R8
	ADDQ AX, 8(SP)
	ADCQ DX, 16(SP)

	// |
	// | Montgomerry Reduction
	MOVQ SI, 24(SP)
	MOVQ CX, 32(SP)

	// |
	// | (u @ R8) = (w0 @ R9) * inp
	MOVQ R9, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w0 @ R9
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R9
	ADCQ DX, SI

	// | w1 @ R10
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R10
	ADCQ DX, CX
	ADDQ SI, R10
	ADCQ $0, CX

	// | w2 @ R11
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ CX, R11
	ADCQ $0, SI

	// | w3 @ R12
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R12
	ADCQ DX, CX
	ADDQ SI, R12
	ADCQ $0, CX

	// | w4 @ R13
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, SI
	ADDQ CX, R13
	ADCQ $0, SI

	// | w5 @ R14
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, CX
	ADDQ SI, R14
	ADCQ $0, CX

	// | w6 @ R15
	ADDQ CX, R15
	MOVQ $0, DI
	ADCQ $0, DI

	// |
	MOVQ 32(SP), R9

	// | (u @ R8) = (w1 @ R10) * inp
	MOVQ R10, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w1 @ R10
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R10
	ADCQ DX, SI

	// | w2 @ R11
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, CX
	ADDQ SI, R11
	ADCQ $0, CX

	// | w3 @ R12
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R12
	ADCQ DX, SI
	ADDQ CX, R12
	ADCQ $0, SI

	// | w4 @ R13
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, CX
	ADDQ SI, R13
	ADCQ $0, CX

	// | w5 @ R14
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, SI
	ADDQ CX, R14
	ADCQ $0, SI

	// | w6 @ R15
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, CX
	ADDQ SI, R15
	ADCQ $0, CX

	// | w7 @ R9
	ADDQ CX, DI
	ADCQ DI, R9
	MOVQ $0, DI
	ADCQ $0, DI

	// |
	MOVQ 24(SP), R10

	// | (u @ R8) = (w2 @ R11) * inp
	MOVQ R11, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w2 @ R11
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, SI

	// | w3 @ R12
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R12
	ADCQ DX, CX
	ADDQ SI, R12
	ADCQ $0, CX

	// | w4 @ R13
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, SI
	ADDQ CX, R13
	ADCQ $0, SI

	// | w5 @ R14
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, CX
	ADDQ SI, R14
	ADCQ $0, CX

	// | w6 @ R15
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, SI
	ADDQ CX, R15
	ADCQ $0, SI

	// | w7 @ R9
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R9
	ADCQ DX, CX
	ADDQ SI, R9
	ADCQ $0, CX

	// | w8 @ R10
	ADDQ CX, DI
	ADCQ DI, R10
	MOVQ $0, DI
	ADCQ $0, DI

	// |
	MOVQ (SP), R11

	// | (u @ R8) = (w3 @ R12) * inp
	MOVQ R12, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w3 @ R12
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R12
	ADCQ DX, SI

	// | w4 @ R13
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, CX
	ADDQ SI, R13
	ADCQ $0, CX

	// | w5 @ R14
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, SI
	ADDQ CX, R14
	ADCQ $0, SI

	// | w6 @ R15
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, CX
	ADDQ SI, R15
	ADCQ $0, CX

	// | w7 @ R9
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ CX, R9
	ADCQ $0, SI

	// | w8 @ R10
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R10
	ADCQ DX, CX
	ADDQ SI, R10
	ADCQ $0, CX

	// | w9 @ R11
	ADDQ CX, DI
	ADCQ DI, R11
	MOVQ $0, DI
	ADCQ $0, DI

	// |
	MOVQ 8(SP), R12

	// | (u @ R8) = (w4 @ R13) * inp
	MOVQ R13, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w4 @ R13
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R13
	ADCQ DX, SI

	// | w5 @ R14
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, CX
	ADDQ SI, R14
	ADCQ $0, CX

	// | w6 @ R15
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, SI
	ADDQ CX, R15
	ADCQ $0, SI

	// | w7 @ R9
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R9
	ADCQ DX, CX
	ADDQ SI, R9
	ADCQ $0, CX

	// | w8 @ R10
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R10
	ADCQ DX, SI
	ADDQ CX, R10
	ADCQ $0, SI

	// | w9 @ R11
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, CX
	ADDQ SI, R11
	ADCQ $0, CX

	// | w10 @ R12
	ADDQ CX, DI
	ADCQ DI, R12
	MOVQ $0, DI
	ADCQ $0, DI

	// |
	MOVQ 16(SP), R13

	// | (u @ R8) = (w5 @ R14) * inp
	MOVQ R14, AX
	MULQ ·inp+0(SB)
	MOVQ AX, R8

	// | w5 @ R14
	XORQ SI, SI
	MOVQ ·modulus+0(SB), AX
	MULQ R8
	ADDQ AX, R14
	ADCQ DX, SI

	// | w6 @ R15
	XORQ CX, CX
	MOVQ ·modulus+8(SB), AX
	MULQ R8
	ADDQ AX, R15
	ADCQ DX, CX
	ADDQ SI, R15
	ADCQ $0, CX

	// | w7 @ R9
	XORQ SI, SI
	MOVQ ·modulus+16(SB), AX
	MULQ R8
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ CX, R9
	ADCQ $0, SI

	// | w8 @ R10
	XORQ CX, CX
	MOVQ ·modulus+24(SB), AX
	MULQ R8
	ADDQ AX, R10
	ADCQ DX, CX
	ADDQ SI, R10
	ADCQ $0, CX

	// | w9 @ R11
	XORQ SI, SI
	MOVQ ·modulus+32(SB), AX
	MULQ R8
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ CX, R11
	ADCQ $0, SI

	// | w10 @ R12
	XORQ CX, CX
	MOVQ ·modulus+40(SB), AX
	MULQ R8
	ADDQ AX, R12
	ADCQ DX, CX
	ADDQ SI, R12
	ADCQ $0, CX

	// | w11 @ R13
	ADDQ CX, DI
	ADCQ DI, R13
	MOVQ $0, DI
	ADCQ $0, DI

	// | Compare & Return
	MOVQ    R15, R8
	SUBQ    ·modulus+0(SB), R8
	MOVQ    R9, R14
	SBBQ    ·modulus+8(SB), R14
	MOVQ    R10, CX
	SBBQ    ·modulus+16(SB), CX
	MOVQ    R11, AX
	SBBQ    ·modulus+24(SB), AX
	MOVQ    R12, DX
	SBBQ    ·modulus+32(SB), DX
	MOVQ    R13, SI
	SBBQ    ·modulus+40(SB), SI
	SBBQ    $0, DI
	MOVQ    c+0(FP), DI
	CMOVQCC R8, R15
	MOVQ    R15, (DI)
	CMOVQCC R14, R9
	MOVQ    R9, 8(DI)
	CMOVQCC CX, R10
	MOVQ    R10, 16(DI)
	CMOVQCC AX, R11
	MOVQ    R11, 24(DI)
	CMOVQCC DX, R12
	MOVQ    R12, 32(DI)
	CMOVQCC SI, R13
	MOVQ    R13, 40(DI)
	RET
