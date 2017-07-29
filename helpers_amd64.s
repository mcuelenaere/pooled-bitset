#include "textflag.h"

TEXT ·hasPopCount(SB),NOSPLIT,$0
    MOVQ    $1, AX
    CPUID
    SHRQ    $23, CX // POPCNT bit
    ANDQ    $1, CX
    MOVB    CX, ret+0(FP)
    RET

TEXT ·hasAvx2(SB),NOSPLIT,$0
    MOVQ    $7, AX
    CPUID
    SHRQ    $5, BX // AVX2 bit
    ANDQ    $1, BX
    MOVB    BX, ret+0(FP)
    RET

TEXT ·popcountSliceAsm(SB),NOSPLIT,$0-32
    MOVQ    $0, AX      // var result = 0
    MOVQ    $0, BX      // var i = 0
    MOVQ	s_base+0(FP), SI // &s
    MOVQ	s_len+8(FP), CX // len(s)

    TESTQ	CX, CX      // if (len(s) == 0)
    JZ		end         // goto end

loop:
    POPCNTQ (SI)(BX*8), DX
    ADDQ	DX, AX      // result += tmp
    INCQ    BX          // i++
    CMPQ    BX, CX      // if (i != len(s))
    JNE     loop        // goto loop

end:
    MOVQ	AX, ret+24(FP)
    RET

TEXT ·findFirstSetBitAsm(SB),NOSPLIT,$0-16
    BSFQ    v+0(FP), AX
    MOVQ	AX, ret+8(FP)
    RET

#define bitOpSliceGeneric(BITOP)                                              \
_1:                                                                           \
    CMPQ    DI, CX          /* if (i == len(a)) */                            \
    JEQ     _end            /* goto _end */                                   \
    MOVQ    (BX)(DI*8), SI  /* tmp := a[i] */                                 \
    BITOP   (DX)(DI*8), SI  /* tmp = bitop(tmp, b[i]) */                      \
    MOVQ    SI, (AX)(DI*8)  /* dest[i] = tmp */                               \
    INCQ    DI              /* i++ */                                         \
    JMP     _1              /* goto _1 */                                     \
_end:                                                                         \
    RET

#define bitOpSliceAvx(BITOP, AVX_BO)                                          \
    MOVQ    dest_base+0(FP), AX /* &dest */                                   \
    MOVQ    a_base+24(FP), BX   /* &a */                                      \
    MOVQ    a_len+32(FP), CX   /* len(a) */                                   \
    MOVQ    b_base+48(FP), DX   /* &b */                                      \
    MOVQ    $0, DI         /* var i = 0 */                                    \
    MOVQ    CX, R8         /* var j = len(a) */                               \
                                                                              \
_128ByteLoop:                                                                 \
    CMPQ    R8, $16        /* if (j < 16) */                                  \
    JB      _64BitLoop     /* goto _64BitLoop */                              \
    VMOVDQU (BX)(DI*8),   Y0                                                  \
    VMOVDQU 32(BX)(DI*8), Y1                                                  \
    VMOVDQU 64(BX)(DI*8), Y2                                                  \
    VMOVDQU 96(BX)(DI*8), Y3                                                  \
    AVX_BO  (DX)(DI*8),   Y0, Y0                                              \
    AVX_BO  32(DX)(DI*8), Y1, Y1                                              \
    AVX_BO  64(DX)(DI*8), Y2, Y2                                              \
    AVX_BO  96(DX)(DI*8), Y3, Y3                                              \
    VMOVDQU Y0, (AX)(DI*8)                                                    \
    VMOVDQU Y1, 32(AX)(DI*8)                                                  \
    VMOVDQU Y2, 64(AX)(DI*8)                                                  \
    VMOVDQU Y3, 96(AX)(DI*8)                                                  \
    ADDQ    $16, DI         /* i += 16 */                                     \
    SUBQ    $16, R8         /* j -= 16 */                                     \
    CMPQ    DI, CX          /* if (i != len(a)) */                            \
    JNE     _128ByteLoop    /* goto _128ByteLoop */                           \
                                                                              \
_64BitLoop:                                                                   \
    bitOpSliceGeneric(BITOP)

#define bitOpSliceSse2(BITOP, SSE_BO)                                         \
    MOVQ    dest+0(FP), AX /* &dest */                                        \
    MOVQ    a+24(FP), BX   /* &a */                                           \
    MOVQ    a+32(FP), CX   /* len(a) */                                       \
    MOVQ    b+48(FP), DX   /* &b */                                           \
    MOVQ    $0, DI         /* var i = 0 */                                    \
    MOVQ    CX, R8         /* var j = len(a) */                               \
                                                                              \
_128ByteLoop:                                                                 \
    CMPQ    R8, $16        /* if (j < 16) */                                  \
    JB      _64BitLoop     /* goto _64BitLoop */                              \
    MOVOU   (BX)(DI*8),   X0                                                  \
    MOVOU   16(BX)(DI*8), X1                                                  \
    MOVOU   32(BX)(DI*8), X2                                                  \
    MOVOU   48(BX)(DI*8), X3                                                  \
    MOVOU   64(BX)(DI*8), X4                                                  \
    MOVOU   80(BX)(DI*8), X5                                                  \
    MOVOU   96(BX)(DI*8), X6                                                  \
    MOVOU   112(BX)(DI*8), X7                                                 \
    MOVOU   (DX)(DI*8),   X8                                                  \
    MOVOU   16(DX)(DI*8), X9                                                  \
    MOVOU   32(DX)(DI*8), X10                                                 \
    MOVOU   48(DX)(DI*8), X11                                                 \
    MOVOU   64(DX)(DI*8), X12                                                 \
    MOVOU   80(DX)(DI*8), X13                                                 \
    MOVOU   96(DX)(DI*8), X14                                                 \
    MOVOU   112(DX)(DI*8), X15                                                \
    SSE_BO  X8,  X0                                                           \
    SSE_BO  X9,  X1                                                           \
    SSE_BO  X10, X2                                                           \
    SSE_BO  X11, X3                                                           \
    SSE_BO  X12, X4                                                           \
    SSE_BO  X13, X5                                                           \
    SSE_BO  X14, X6                                                           \
    SSE_BO  X15, X7                                                           \
    MOVOU   X0, (AX)(DI*8)                                                    \
    MOVOU   X1, 16(AX)(DI*8)                                                  \
    MOVOU   X2, 32(AX)(DI*8)                                                  \
    MOVOU   X3, 48(AX)(DI*8)                                                  \
    MOVOU   X4, 64(AX)(DI*8)                                                  \
    MOVOU   X5, 80(AX)(DI*8)                                                  \
    MOVOU   X6, 96(AX)(DI*8)                                                  \
    MOVOU   X7, 112(AX)(DI*8)                                                 \
    ADDQ    $16, DI         /* i += 16 */                                     \
    SUBQ    $16, R8         /* j -= 16 */                                     \
    CMPQ    DI, CX          /* if (i != len(a)) */                            \
    JNE     _128ByteLoop    /* goto _128ByteLoop */                           \
                                                                              \
_64BitLoop:                                                                   \
    bitOpSliceGeneric(BITOP)

TEXT ·andSliceAvx2(SB),NOSPLIT,$0-72
    bitOpSliceAvx(ANDQ, VPAND)

TEXT ·andSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(ANDQ, ANDPS)

TEXT ·orSliceAvx2(SB),NOSPLIT,$0-72
    bitOpSliceAvx(ORQ, VPOR)

TEXT ·orSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(ORQ, ORPS)

TEXT ·xorSliceAvx2(SB),NOSPLIT,$0-72
    bitOpSliceAvx(XORQ, VPXOR)

TEXT ·xorSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(XORQ, XORPS)

#define notSliceGeneric()                                                     \
_1:                                                                           \
    CMPQ    DI, CX          /* if (i == len(a)) */                            \
    JEQ     _end            /* goto _end */                                   \
    MOVQ    (BX)(DI*8), SI  /* tmp := a[i] */                                 \
    NOTQ    SI              /* tmp = ^tmp */                                  \
    MOVQ    SI, (AX)(DI*8)  /* dest[i] = tmp */                               \
    INCQ    DI              /* i++ */                                         \
    JMP     _1              /* goto _1 */                                     \
_end:                                                                         \
    RET

TEXT ·notSliceSse2(SB),NOSPLIT,$0-48
    MOVQ    dest_base+0(FP), AX                                     // &dest
    MOVQ    src_base+24(FP), BX                                     // &src
    MOVQ    src_len+32(FP), CX                                      // len(src)
    MOVQ    $0, DI                                                  // var i = 0
    MOVQ    CX, R8                                                  // var j = len(src)

    // setup X8 register
    MOVQ    $0xffffffffffffffff, SI
    MOVQ    SI, X8
    MOVLHPS X8, X8

_128ByteLoop:
    CMPQ    R8, $16                                                 // if (j < 16)
    JB      _64BitLoop                                              // goto _64BitLoop
    MOVOU   (BX)(DI*8),   X0                                                  
    MOVOU   16(BX)(DI*8), X1
    MOVOU   32(BX)(DI*8), X2
    MOVOU   48(BX)(DI*8), X3
    MOVOU   64(BX)(DI*8), X4
    MOVOU   80(BX)(DI*8), X5
    MOVOU   96(BX)(DI*8), X6
    MOVOU   112(BX)(DI*8), X7
    XORPS   X8, X0
    XORPS   X8, X1
    XORPS   X8, X2
    XORPS   X8, X3
    XORPS   X8, X4
    XORPS   X8, X5
    XORPS   X8, X6
    XORPS   X8, X7
    MOVOU   X0, (AX)(DI*8)                                                    
    MOVOU   X1, 16(AX)(DI*8)
    MOVOU   X2, 32(AX)(DI*8)
    MOVOU   X3, 48(AX)(DI*8)
    MOVOU   X4, 64(AX)(DI*8)
    MOVOU   X5, 80(AX)(DI*8)
    MOVOU   X6, 96(AX)(DI*8)
    MOVOU   X7, 112(AX)(DI*8)
    ADDQ    $16, DI                                                 // i += 16
    SUBQ    $16, R8                                                 // j -= 16
    CMPQ    DI, CX                                                  // if (i != len(a))
    JNE     _128ByteLoop                                            // goto _128ByteLoop

_64BitLoop:
    notSliceGeneric()

TEXT ·notSliceAvx2(SB),NOSPLIT,$0-48
    MOVQ    dest_base+0(FP), AX                                     // &dest
    MOVQ    src_base+24(FP), BX                                     // &src
    MOVQ    src_len+32(FP), CX                                      // len(src)
    MOVQ    $0, DI                                                  // var i = 0
    MOVQ    CX, R8                                                  // var j = len(src)

    // setup Y4 register
    VPCMPEQB Y4, Y4, Y4

_128ByteLoop:
    CMPQ     R8, $16                                                 // if (j < 16)
    JB       _64BitLoop                                              // goto _64BitLoop
    VMOVDQU  (BX)(DI*8),   Y0
    VMOVDQU  32(BX)(DI*8),  Y1
    VMOVDQU  64(BX)(DI*8), Y2
    VMOVDQU  96(BX)(DI*8), Y3
    VPXOR    Y0, Y4, Y0
    VPXOR    Y1, Y4, Y1
    VPXOR    Y2, Y4, Y2
    VPXOR    Y3, Y4, Y3
    VMOVDQU  Y0, (AX)(DI*8)
    VMOVDQU  Y1, 32(AX)(DI*8)
    VMOVDQU  Y2, 64(AX)(DI*8)
    VMOVDQU  Y3, 96(AX)(DI*8)
    ADDQ     $16, DI                                                 // i += 16
    SUBQ     $16, R8                                                 // j -= 16
    CMPQ     DI, CX                                                  // if (i != len(a))
    JNE      _128ByteLoop                                            // goto _128ByteLoop

_64BitLoop:
    notSliceGeneric()
