TEXT ·hasPopCount(SB),NOSPLIT,$0
    MOVQ    $1, AX
    CPUID
    SHRQ    $23, CX // POPCNT bit
    ANDQ    $1, CX
    MOVB    CX, ret+0(FP)
    RET

TEXT ·hasAvx(SB),NOSPLIT,$0
    MOVQ    $1, AX
    CPUID
    SHRQ    $28, CX // AVX bit
    ANDQ    $1, CX
    MOVB    CX, ret+0(FP)
    RET

TEXT ·popcountSliceAsm(SB),NOSPLIT,$0-32
    MOVQ    $0, AX      // var result = 0
    MOVQ    $0, BX      // var i = 0
    MOVQ	s+0(FP), SI // &s
    MOVQ	s+8(FP), CX // len(s)

    TESTQ	CX, CX      // if (len(s) == 0)
    JZ		end         // goto end

loop:
    // POPCNTQ (SI)(BX*8), DX
    BYTE $0xF3; BYTE $0x48; BYTE $0x0F; BYTE $0xB8; BYTE $0x14; BYTE $0xDE
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

#define bitOpSliceAvx(BITOP, AVX_BO)                                          \
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
    /* VMOVDQU   (BX)(DI*8), Y0 */                                            \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x6F; BYTE $0x04; BYTE $0xFB                \
    /* VMOVDQU   8(BX)(DI*8), Y1 */                                           \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x6F; BYTE $0x4C; BYTE $0xFB; BYTE $0x08    \
    /* VMOVDQU   16(BX)(DI*8), Y2 */                                          \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x6F; BYTE $0x54; BYTE $0xFB; BYTE $0x10    \
    /* VMOVDQU   24(BX)(DI*8), Y3 */                                          \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x6F; BYTE $0x5C; BYTE $0xFB; BYTE $0x18    \
    /* AVX_BO (DX)(DI*8), Y0 */                                               \
    BYTE $0xC5; BYTE $0xFD; BYTE AVX_BO; BYTE $0x04; BYTE $0xFA               \
    /* AVX_BO 8(DX)(DI*8), Y1 */                                              \
    BYTE $0xC5; BYTE $0xF5; BYTE AVX_BO; BYTE $0x4C; BYTE $0xFA; BYTE $0x08   \
    /* AVX_BO 16(DX)(DI*8), Y2 */                                             \
    BYTE $0xC5; BYTE $0xED; BYTE AVX_BO; BYTE $0x54; BYTE $0xFA; BYTE $0x10   \
    /* AVX_BO 24(DX)(DI*8), Y3 */                                             \
    BYTE $0xC5; BYTE $0xE5; BYTE AVX_BO; BYTE $0x5C; BYTE $0xFA; BYTE $0x18   \
    /* VMOVDQU   Y0, (AX)(DI*8) */                                            \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x7F; BYTE $0x04; BYTE $0xF8                \
    /* VMOVDQU   Y1, 8(AX)(DI*8) */                                           \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x7F; BYTE $0x4C; BYTE $0xF8; BYTE $0x08    \
    /* VMOVDQU   Y2, 16(AX)(DI*8) */                                          \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x7F; BYTE $0x54; BYTE $0xF8; BYTE $0x10    \
    /* VMOVDQU   Y3, 24(AX)(DI*8) */                                          \
    BYTE $0xC5; BYTE $0xFE; BYTE $0x7F; BYTE $0x5C; BYTE $0xF8; BYTE $0x18    \
    ADDQ    $16, DI         /* i += 16 */                                     \
    SUBQ    $16, R8         /* j -= 16 */                                     \
    CMPQ    DI, CX          /* if (i != len(a)) */                            \
    JNE     _128ByteLoop    /* goto _128ByteLoop */                           \
                                                                              \
_64BitLoop:                                                                   \
    CMPQ    DI, CX          /* if (i == len(a)) */                            \
    JEQ     _end            /* goto _end */                                   \
    MOVQ    (BX)(DI*8), SI  /* tmp := a[i] */                                 \
    BITOP   (DX)(DI*8), SI  /* tmp = bitop(tmp, b[i]) */                      \
    MOVQ    SI, (AX)(DI*8)  /* dest[i] = tmp */                               \
    INCQ    DI              /* i++ */                                         \
    JMP     _64BitLoop      /* goto _64BitLoop */                             \
_end:                                                                         \
    RET

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
    MOVOU   8(BX)(DI*8),  X1                                                  \
    MOVOU   16(BX)(DI*8), X2                                                  \
    MOVOU   24(BX)(DI*8), X3                                                  \
    MOVOU   32(BX)(DI*8), X4                                                  \
    MOVOU   40(BX)(DI*8), X5                                                  \
    MOVOU   48(BX)(DI*8), X6                                                  \
    MOVOU   56(BX)(DI*8), X7                                                  \
    MOVOU   (DX)(DI*8),   X8                                                  \
    MOVOU   8(DX)(DI*8),  X9                                                  \
    MOVOU   16(DX)(DI*8), X10                                                 \
    MOVOU   24(DX)(DI*8), X11                                                 \
    MOVOU   32(DX)(DI*8), X12                                                 \
    MOVOU   40(DX)(DI*8), X13                                                 \
    MOVOU   48(DX)(DI*8), X14                                                 \
    MOVOU   56(DX)(DI*8), X15                                                 \
    SSE_BO  X8,  X0                                                           \
    SSE_BO  X9,  X1                                                           \
    SSE_BO  X10, X2                                                           \
    SSE_BO  X11, X3                                                           \
    SSE_BO  X12, X4                                                           \
    SSE_BO  X13, X5                                                           \
    SSE_BO  X14, X6                                                           \
    SSE_BO  X15, X7                                                           \
    MOVOU   X0, (AX)(DI*8)                                                    \
    MOVOU   X1, 8(AX)(DI*8)                                                   \
    MOVOU   X2, 16(AX)(DI*8)                                                  \
    MOVOU   X3, 24(AX)(DI*8)                                                  \
    MOVOU   X4, 32(AX)(DI*8)                                                  \
    MOVOU   X5, 40(AX)(DI*8)                                                  \
    MOVOU   X6, 48(AX)(DI*8)                                                  \
    MOVOU   X7, 56(AX)(DI*8)                                                  \
    ADDQ    $16, DI         /* i += 16 */                                     \
    SUBQ    $16, R8         /* j -= 16 */                                     \
    CMPQ    DI, CX          /* if (i != len(a)) */                            \
    JNE     _128ByteLoop    /* goto _128ByteLoop */                           \
                                                                              \
_64BitLoop:                                                                   \
    CMPQ    DI, CX          /* if (i == len(a)) */                            \
    JEQ     end             /* goto _end */                                   \
    MOVQ    (BX)(DI*8), SI  /* tmp := a[i] */                                 \
    BITOP   (DX)(DI*8), SI  /* tmp = bitop(tmp, b[i]) */                      \
    MOVQ    SI, (AX)(DI*8)  /* dest[i] = tmp */                               \
    INCQ    DI              /* i++ */                                         \
    JMP     _64BitLoop      /* goto _64BitLoop */                             \
                                                                              \
end:                                                                          \
    RET

TEXT ·andSliceAvx(SB),NOSPLIT,$0-72
    bitOpSliceAvx(ANDQ, $0x54)

TEXT ·andSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(ANDQ, ANDPS)

TEXT ·orSliceAvx(SB),NOSPLIT,$0-72
    bitOpSliceAvx(ORQ, $0x56)

TEXT ·orSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(ORQ, ORPS)

TEXT ·xorSliceAvx(SB),NOSPLIT,$0-72
    bitOpSliceAvx(XORQ, $0x57)

TEXT ·xorSliceSse2(SB),NOSPLIT,$0-72
    bitOpSliceSse2(XORQ, XORPS)

// TODO: create AVX version of this one
TEXT ·notSliceSse2(SB),NOSPLIT,$0-48
    MOVQ    dest+0(FP), AX                                          // &dest
    MOVQ    a+24(FP), BX                                            // &src
    MOVQ    a+32(FP), CX                                            // len(src)
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
    MOVOU   8(BX)(DI*8),  X1
    MOVOU   16(BX)(DI*8), X2
    MOVOU   24(BX)(DI*8), X3
    MOVOU   32(BX)(DI*8), X4
    MOVOU   40(BX)(DI*8), X5
    MOVOU   48(BX)(DI*8), X6
    MOVOU   56(BX)(DI*8), X7
    XORPS   X8, X0
    XORPS   X8, X1
    XORPS   X8, X2
    XORPS   X8, X3
    XORPS   X8, X4
    XORPS   X8, X5
    XORPS   X8, X6
    XORPS   X8, X7
    MOVOU   X0, (AX)(DI*8)
    MOVOU   X1, 8(AX)(DI*8)
    MOVOU   X2, 16(AX)(DI*8)
    MOVOU   X3, 24(AX)(DI*8)
    MOVOU   X4, 32(AX)(DI*8)
    MOVOU   X5, 40(AX)(DI*8)
    MOVOU   X6, 48(AX)(DI*8)
    MOVOU   X7, 56(AX)(DI*8)
    ADDQ    $16, DI                                                 // i += 16
    SUBQ    $16, R8                                                 // j -= 16
    CMPQ    DI, CX                                                  // if (i != len(a))
    JNE     _128ByteLoop                                            // goto _128ByteLoop

_64BitLoop:
    CMPQ    DI, CX                                                  // if (i == len(a))
    JEQ     end                                                     // goto end
    MOVQ    (BX)(DI*8), SI                                          // tmp := src[i]
    XORQ    $0xffffffffffffffff, SI                                 // tmp = ^tmp
    MOVQ    SI, (AX)(DI*8)                                          // dest[i] = tmp
    INCQ    DI                                                      // i++
    JMP     _64BitLoop                                              // goto _64BitLoop

end:
    RET
