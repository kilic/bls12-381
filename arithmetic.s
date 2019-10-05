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
TEXT ·montmul(SB), NOSPLIT, $16-24

/*   | inputs               */	

	MOVQ a+8(FP), DI
	MOVQ b+16(FP), SI

/*   | multiplication phase */

  // | w = a * b
  // | a = (a0, a1, a2, a3, a4, a5)
  // | b = (b0, b1, b2, b3, b4, b5)
  // | w = (w0, w1, w2, w3, w4, w5, w6, w7, w8, w9, w10, w11)

	MOVQ $0, R9
  MOVQ $0, R10
  MOVQ $0, R11
  MOVQ $0, R12
  MOVQ $0, R13
  MOVQ $0, BX

	// | b0 @ CX
	MOVQ (SI), CX

	// | a0 * b0
	// | (w0, w1) @ (SP, R8)
	MOVQ (DI), AX
	MULQ CX
	MOVQ AX, 0(SP)
	MOVQ DX, R8

	// | a1 * b0
	// | (w1, w2) @ (R8, R9)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9

	// | a2 * b0
	// | (w2, w3) @ (R9, R10)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10

	// | a3 * b0
	// | (w3, w4) @ (R10, R11)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11

	// | a4 * b0
	// | (w4, w5) @ (R11, R12)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12

	// | a5 * b0
	// | (w5, w6) @ (R12, R13)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13

	// | b1 @ CX
	MOVQ 8(SI), CX

	// | a0 * b1
	// | (w1, w2, w3, w4) @ (R8, R9, R10, BX)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9
	ADCQ $0, R10
	ADCQ $0, BX
  // | w1 @ 8(SP)
  MOVQ R8, 8(SP)

	// | a1 * b1
	// | (w2, w3, w4, w5) @ (R9, R10, R11, BX)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ BX, R11
  MOVQ $0, BX
	ADCQ $0, BX
	// | w2 @ R8
  MOVQ R9, R8

	// | a2 * b1
	// | (w3, w4, w5, w6) @ (R10, R11, R12, BX)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ BX, R12
  MOVQ $0, BX
	ADCQ $0, BX
	// | w3 @ R9
  MOVQ R10, R9

	// | a3 * b1
	// | (w4, w5, w6, w7) @ (R11, R12, R13, BX)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
  ADCQ BX, R13
  MOVQ $0, BX
	ADCQ $0, BX
	// | w4 @ R10
  MOVQ R11, R10

	// | a4 * b1
	// | (w5, w6, w7) @ (R12, R13, BX)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
  ADCQ $0, BX
	// | w5 @ R11
  MOVQ R12, R11

	// | a5 * b1
	// | (w6, w7) @ (R13, BX)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, BX
	// | w6 @ R12
  MOVQ R13, R12
	// | w7 @ R13
  MOVQ BX, R13

	// | b2 @ CX
	MOVQ 16(SI), CX
  MOVQ $0, BX

	// | a0 * b2
	// | (w2, w3, w4, w5) @ (R8, R9, R10, BX)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9
	ADCQ $0, R10
	ADCQ $0, BX
  // | w2 @ 8(SP)
  MOVQ R8, 16(SP)

	// | a1 * b2
	// | (w3, w4, w5, w6) @ (R9, R10, R11, BX)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ BX, R11
  MOVQ $0, BX
	ADCQ $0, BX
	// | w3 @ R8
  MOVQ R9, R8

	// | a2 * b2
	// | (w4, w5, w6, w7) @ (R10, R11, R12, BX)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ BX, R12
  MOVQ $0, BX
	ADCQ $0, BX
	// | w4 @ R9
  MOVQ R10, R9

	// | a3 * b2
	// | (w5, w6, w7, w8) @ (R11, R12, R13, BX)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
  ADCQ BX, R13
  MOVQ $0, BX
	ADCQ $0, BX
	// | w5 @ R10
  MOVQ R11, R10

	// | a4 * b2
	// | (w6, w7, w8) @ (R12, R13, BX)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
  ADCQ $0, BX
	// | w6 @ R11
  MOVQ R12, R11

	// | a5 * b1
	// | (w7, w8) @ (R13, BX)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, BX
	// | w7 @ R12
  MOVQ R13, R12
	// | w8 @ R13
  MOVQ BX, R13

	// | b3 @ CX
	MOVQ 24(SI), CX
  MOVQ $0, BX

	// | a0 * b3
	// | (w3, w4, w5, w6) @ (R8, R9, R10, BX)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9
	ADCQ $0, R10
	ADCQ $0, BX
  // | w3 @ 8(SP)
	MOVQ R8, R14

	// | a1 * b3
	// | (w4, w5, w6, w7) @ (R9, R10, R11, BX)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ BX, R11
  MOVQ $0, BX
	ADCQ $0, BX
  // | w4 @ R8
  MOVQ R9, R8

	// | a2 * b3
	// | (w5, w6, w7, w8) @ (R10, R11, R12, BX)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ BX, R12
  MOVQ $0, BX
	ADCQ $0, BX
	// | w5 @ R9
  MOVQ R10, R9

	// | a3 * b3
	// | (w6, w7, w8, w9) @ (R11, R12, R13, BX)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
  ADCQ BX, R13
  MOVQ $0, BX
	ADCQ $0, BX
	// | w6 @ R10
  MOVQ R11, R10

	// | a4 * b3
	// | (w7, w8, w9) @ (R12, R13, BX)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
  ADCQ $0, BX
	// | w7 @ R11
  MOVQ R12, R11

	// | a5 * b3
	// | (w8, w9) @ (R13, BX)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, BX
	// | w8 @ R12
  MOVQ R13, R12
	// | w9 @ R13
  MOVQ BX, R13

	// | b4 @ CX
	MOVQ 32(SI), CX
  MOVQ $0, BX

	// | a0 * b4
	// | (w4, w5, w6, w7) @ (R8, R9, R10, BX)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9
	ADCQ $0, R10
	ADCQ $0, BX
  // | w4 @ 8(SP)
  MOVQ R8, R15

	// | a1 * b4
	// | (w5, w6, w7, w8) @ (R9, R10, R11, BX)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ BX, R11
  MOVQ $0, BX
	ADCQ $0, BX
	// | w5 @ R8
  MOVQ R9, R8

	// | a2 * b4
	// | (w6, w7, w8, w9) @ (R10, R11, R12, BX)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ BX, R12
  MOVQ $0, BX
	ADCQ $0, BX
	// | w6 @ R9
  MOVQ R10, R9

	// | a3 * b4
	// | (w7, w8, w9, w10) @ (R11, R12, R13, BX)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
  ADCQ BX, R13
  MOVQ $0, BX
	ADCQ $0, BX
	// | w7 @ R10
  MOVQ R11, R10

	// | a4 * b4
	// | (w8, w9, w10) @ (R12, R13, BX)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
  ADCQ $0, BX
	// | w8 @ R11
  MOVQ R12, R11

	// | a5 * b4
	// | (w9, w10) @ (R13, BX)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, BX
	// | w9 @ R12
  MOVQ R13, R12
	// | w10 @ R13
  MOVQ BX, R13

	// | b5 @ CX
	MOVQ 40(SI), CX
  MOVQ $0, BX

	// | a0 * b5
	// | (w5, w6, w7, w8) @ (R8, R9, R10, BX)
	MOVQ (DI), AX
	MULQ CX
	ADDQ AX, R8
	ADCQ DX, R9
	ADCQ $0, R10
	ADCQ $0, BX

	// | a1 * b5
	// | (w6, w7, w8, w9) @ (R9, R10, R11, BX)
	MOVQ 8(DI), AX
	MULQ CX
	ADDQ AX, R9
	ADCQ DX, R10
	ADCQ BX, R11
  MOVQ $0, BX
	ADCQ $0, BX

	// | a2 * b5
	// | (w7, w8, w9, w10) @ (R10, R11, R12, BX)
	MOVQ 16(DI), AX
	MULQ CX
	ADDQ AX, R10
	ADCQ DX, R11
	ADCQ BX, R12
  MOVQ $0, BX
	ADCQ $0, BX

	// | a3 * b5
	// | (w8, w9, w10, w11) @ (R11, R12, R13, BX)
	MOVQ 24(DI), AX
	MULQ CX
	ADDQ AX, R11
	ADCQ DX, R12
  ADCQ BX, R13
  MOVQ $0, BX
	ADCQ $0, BX

	// | a4 * b5
	// | (w9, w10, w11) @ (R12, R13, BX)
	MOVQ 32(DI), AX
	MULQ CX
	ADDQ AX, R12
	ADCQ DX, R13
  ADCQ $0, BX

	// | a5 * b5
	// | (w10, w11) @ (R13, BX)
	MOVQ 40(DI), AX
	MULQ CX
	ADDQ AX, R13
	ADCQ DX, BX

	// |  w0,  w1,  w2,   w3,  w4,  w5,
	// | 	 0,   8,  16,  R14, R15,  R8,
	// |  w6,  w7,  w8,  w9,  w10, w11,
	// |  R9,  R10, R11, R12, R13,  BX,

	// | 
	// | Montgomerry Reduction Phase
	// | c = w % p

/*   | swap                 */

	MOVQ 0(SP), SI
	MOVQ 8(SP), DI
	MOVQ 16(SP), CX
	MOVQ R13, 0(SP)
	MOVQ BX, 8(SP)
	// | R13 will be the carry register

/*   | mont_i0              */

	// |  w0,  w1,  w2,  w3,  w4,  w5,  w6,  w7,  w8,  w9, w10, w11
	// |  SI,  DI,  CX, R14, R15,  R8,  R9, R10, R11, R12,   0,   8

	// | i = 0
	// | (u @ BX) = (w0 @ SI) * inverse_p
	MOVQ SI, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// | SI is idle now

	// | w1 @ DI
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, DI
	ADCQ $0, DX
	ADDQ R13, DI
	MOVQ $0, R13
	ADCQ DX, R13

	// | w2 @ CX
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, CX
	ADCQ $0, DX
	ADDQ R13, CX
	MOVQ $0, R13
	ADCQ DX, R13

	// | w3 @ R14
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R14
	ADCQ $0, DX
	ADDQ R13, R14
	MOVQ $0, R13
	ADCQ DX, R13

	// | w4 @ R15
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R15
	ADCQ $0, DX
	ADDQ R13, R15
	MOVQ $0, R13
	ADCQ DX, R13

	// | w5 @ R8
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, R8
	ADCQ $0, DX
	ADDQ R13, R8
	// | w6 @ R9
	ADCQ DX, R9
	
	// | long_carry @ SI should be added to w7
	MOVQ $0, SI
	ADCQ $0, SI

/*   | mont_i1:             */

	// |  lc,  w1,  w2,  w3,  w4,  w5,  w6,  w7,  w8,  w9, w10, w11,
	// |  SI,  DI,  CX, R14, R15,  R8,  R9, R10, R11, R12,   0,   8,

	// | i = 1
	// | (u @ BX) = (w1 @ DI) * inverse_p
	MOVQ DI, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// | DI is idle now

	// | w2 @ CX
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, CX
	ADCQ $0, DX
	ADDQ R13, CX
	MOVQ $0, R13
	ADCQ DX, R13

	// | w3 @ R14
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, R14
	ADCQ $0, DX
	ADDQ R13, R14
	MOVQ $0, R13
	ADCQ DX, R13

	// | w4 @ R15
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R15
	ADCQ $0, DX
	ADDQ R13, R15
	MOVQ $0, R13
	ADCQ DX, R13

	// | w5 @ R8
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R8
	ADCQ $0, DX
	ADDQ R13, R8
	MOVQ $0, R13
	ADCQ DX, R13

	// | w6 @ R9
	// | in the last round of the iteration
	// | we don't use the short carry @ R13
	// | instead we bring back long_carry @ SI
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, R9
	ADCQ DX, SI
	ADDQ R13, R9
	// | w7 @ R10
	ADCQ SI, R10
	// | long_carry @ DI should be added to w8
	MOVQ $0, SI
	ADCQ $0, SI
	
/*   | mont_i2              */

	// |  lc,  - ,  w2,  w3,  w4,  w5,  w6,  w7,  w8,  w9, w10, w11
	// |  SI,  DI,  CX, R14, R15,  R8,  R9, R10, R11, R12,   0,   8

	// | i = 2
	// | (u @ BX) = (w2 @ CX) * inverse_p
	MOVQ CX, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// CX is idle now

	// | w3 @ R14
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, R14
	ADCQ $0, DX
	ADDQ R13, R14
	MOVQ $0, R13
	ADCQ DX, R13

	// | w4 @ R15
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, R15
	ADCQ $0, DX
	ADDQ R13, R15
	MOVQ $0, R13
	ADCQ DX, R13

	// | w5 @ R8
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R8
	ADCQ $0, DX
	ADDQ R13, R8
	MOVQ $0, R13
	ADCQ DX, R13

	// | w6 @ R9
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R9
	ADCQ $0, DX
	ADDQ R13, R9
	MOVQ $0, R13
	ADCQ DX, R13

	// | w7 @ R10
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, R10
	ADCQ DX, SI
	ADDQ R13, R10
	// | w8 @ R11
	ADCQ SI, R11
	// | long_carry @ SI should be added to w9
	MOVQ $0, SI
	ADCQ $0, SI

/*   | mont_i3:             */

	// |  lc,  - ,  - ,  w3,  w4,  w5,  w6,  w7,  w8,  w9, w10, w11
	// |  SI,  DI,  CX, R14, R15,  R8,  R9, R10, R11, R12,   0,  8

	// | i = 3
	// | (u @ BX) = (w3 @ R14) * inverse_p
	MOVQ R14, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// R14 is idle now

	// | w4 @ R15
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, R15
	ADCQ $0, DX
	ADDQ R13, R15
	MOVQ $0, R13
	ADCQ DX, R13

	// | w5 @ R8
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, R8
	ADCQ $0, DX
	ADDQ R13, R8
	MOVQ $0, R13
	ADCQ DX, R13

	// | w6 @ R9
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R9
	ADCQ $0, DX
	ADDQ R13, R9
	MOVQ $0, R13
	ADCQ DX, R13

	// | w7 @ R10
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R10
	ADCQ $0, DX
	ADDQ R13, R10
	MOVQ $0, R13
	ADCQ DX, R13

	// | w8 @ R11
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, R11
	ADCQ DX, SI
	ADDQ R13, R11
	// | w9 @ R12
	ADCQ SI, R12
	// | long_carry @ SI should be added to w10
	MOVQ $0, SI
	ADCQ $0, SI
 	
/*   | mont_i4:             */

	// |  lc,  - ,  - ,  - ,  w4,  w5,  w6,  w7,  w8,  w9, w10, w11
	// |  SI,  DI,  CX, R14, R15,  R8,  R9, R10, R11, R12,   0,  8

  // | i = 4
	// | (u @ BX) = (w4 @ R15) * inverse_p
	MOVQ R15, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// R15 is idle now

	// | w5 @ R8
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, R8
	ADCQ $0, DX
	ADDQ R13, R8
	MOVQ $0, R13
	ADCQ DX, R13

	// | w6 @ R9
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, R9
	ADCQ $0, DX
	ADDQ R13, R9
	MOVQ $0, R13
	ADCQ DX, R13

	// | w7 @ R10
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R10
	ADCQ $0, DX
	ADDQ R13, R10
	MOVQ $0, R13
	ADCQ DX, R13

	// | w8 @ R11
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R11
	ADCQ $0, DX
	ADDQ R13, R11
	MOVQ $0, R13
	ADCQ DX, R13

	// | w9 @ R12
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, R12
	ADCQ DX, SI
	ADDQ R13, R12

/*   | swap                 */

  // | from stack to available registers
	// | w10 @ CX
	// | w11 @ R14
	MOVQ 0(SP), CX
	MOVQ 8(SP), R14

	// | w10 @ DI
	ADCQ SI, CX
	// | long_carry @ SI should be added to w11
	ADCQ $0, R14

/*   | mont_i5:             */

	// |  lc,  - ,  - ,  w5,  w6,  w7,  w8,  w9, w10, w11
	// |  SI,  DI, R15,  R8,  R9, R10, R11, R12,  CX, R14

	// | i = 5
	// | (u @ BX) = (w5 @ R8) * inverse_p
	MOVQ R8, AX
	MULQ ·inp+0(SB)
	MOVQ AX, BX
	MOVQ $0, R13
	MOVQ ·modulus+0(SB), AX
	MULQ BX
	ADCQ DX, R13
	// R8 is idle now

	// | w6 @ R9
	MOVQ ·modulus+8(SB), AX
	MULQ BX
	ADDQ AX, R9
	ADCQ $0, DX
	ADDQ R13, R9
	MOVQ $0, R13
	ADCQ DX, R13

		// | w7 @ R10
	MOVQ ·modulus+16(SB), AX
	MULQ BX
	ADDQ AX, R10
	ADCQ $0, DX
	ADDQ R13, R10
	MOVQ $0, R13
	ADCQ DX, R13

	// | w8 @ R11
	MOVQ ·modulus+24(SB), AX
	MULQ BX
	ADDQ AX, R11
	ADCQ $0, DX
	ADDQ R13, R11
	MOVQ $0, R13
	ADCQ DX, R13

	// | w9 @ R12
	MOVQ ·modulus+32(SB), AX
	MULQ BX
	ADDQ AX, R12
	ADCQ $0, DX
	ADDQ R13, R12
	ADCQ DX, CX
	ADCQ $0, R14

	// | (w10, w11) @ (CX, R14)
	MOVQ ·modulus+40(SB), AX
	MULQ BX
	ADDQ AX, CX
	ADCQ DX, R14

/*   | reduction            */

  // | c = (w6, w7, w8, w9, w10, w11) @ (R9, R10, R11, DI, CX, R14)
	MOVQ R9, AX
	MOVQ R10, BX
	MOVQ R11, DX
	MOVQ R12, R8
	MOVQ CX, R15
	MOVQ R14, R13
	SUBQ ·modulus+0(SB), AX
	SBBQ ·modulus+8(SB), BX
	SBBQ ·modulus+16(SB), DX
	SBBQ ·modulus+24(SB), R8
	SBBQ ·modulus+32(SB), R15
	SBBQ ·modulus+40(SB), R13
	CMOVQCC AX, R9
	CMOVQCC BX, R10
	CMOVQCC DX, R11
	CMOVQCC R8, R12
	CMOVQCC R15, CX
	CMOVQCC R13, R14

/*   | return               */

	MOVQ c+0(FP), SI
	MOVQ R9, (SI)
	MOVQ R10, 8(SI)
	MOVQ R11, 16(SI)
	MOVQ R12, 24(SI)
	MOVQ CX, 32(SI)
	MOVQ R14, 40(SI)
  RET
