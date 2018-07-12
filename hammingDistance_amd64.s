//; ============================================================================
//; ===[HammingDistance]========================================================
//; ============================================================================

//; SI = Pointer to a
//; DI = Pointer to b
//; R8 = Bytes left to check
//; AX = From hdHuge and onwards AX is used to count the number of mismatching bytes
//; BX, CX, DX, R9 are used as temporary registers

//; func HammingDistance(a, b []byte) int
TEXT ·HammingDistance(SB),7,$0
	MOVQ a+0(FP), SI  //; pointer to a
	MOVQ a+8(FP), CX  //; length of a
	MOVQ b+24(FP), DI //; pointer to b
	MOVQ b+32(FP), R8 //; length of b
//; return value should be saved in res+48(FP)

//; Return if not equal length
	CMPQ CX, R8
	JNE error
//; return if length is zero
	CMPQ R8 ,$0
	JE error

	XORQ AX, AX
//; If the length of a and b is less than 192 bytes jump to hdSmallEntry
	CMPQ R8, $192
	JB hdSmallEntry

//; TODO: Make another version of this HammingDistance function without this
//; check and with a good comment that explicitly warns that we need to check if
//; we have all CPU fetures that we need before calling the new function.

//; Check that we have access to all instructions we need
	MOVQ $1, AX	//; AX=1: Processor Info and Feature Bits
	CPUID				//; MOVOU (MOVDQU) needs sse2, bit 26 in EDX
						//; PCMPEQB needs sse2, bit 26 in EDX
						//; PMOVMSKB needs sse2, bit 26 in EDX
						//; sse2 came with the Pentium 4 in 2000
						//; POPCNTQ needs popcnt bit 23 ECX
						//; The popcnt CPUID feture is the newest feture that we need.
						//; popcnt came with the Nehalem architecture in 2008.

	ANDQ $0x4000000, DX	//; AND against bit 26 in EDX
	JZ missingSSE2
	ANDQ $0x800000, CX	//; AND against bit 23 in ECX
	JZ missingPOPCNTQ

//; AX is used to count the number of mismatching bytes
	XORQ AX, AX //; set AX to zero

hdHuge:
//; Copy in 64 bytes of a and b in to the xmm registers X0 to X7
	MOVOU (SI), X0			//;						MOVDQU on Sandy Bridge (m128, xmm) 1 cycles/instructions and 3 in latency
	MOVOU 16(SI), X2
	MOVOU 32(SI), X4
	MOVOU 48(SI), X6
	MOVOU (DI), X1
	MOVOU 16(DI), X3
	MOVOU 32(DI), X5
	MOVOU 48(DI), X7

//; PCMPEQB will set X0, X2, X4 and X6 to only contain 0xFF if corresponding
//; bytes in a and b are the same.
	PCMPEQB X1, X0			//;						Sandy Bridge (xmm xmm) 0.5 cycles/instructions and 1 in latency
	PCMPEQB X3, X2
	PCMPEQB X5, X4
	PCMPEQB X7, X6

//; PMOVMSKB takes the highest bit from every byte in X0, X2, X4 and X6
//; respectivly and copys them in to the 16 least significant bits of
//; R9, BX, CX and DX respectivly.
	PMOVMSKB X0, R9		//;						Sandy Bridge (reg im) 1 cycles/instructions and 2 in latency
	PMOVMSKB X2, BX
	PMOVMSKB X4, CX
	PMOVMSKB X6, DX

//; Move all 64 resulting bits in to R9
	SHLQ $48, R9		//;	XX______				Sandy Bridge (reg im) 0.5 cycles/instructions and 1 in latency
	SHRQ $16, R9:BX	//;	XXXX____				Sandy Bridge (reg im im) 0.5 cycles/instructions
	SHRQ $16, R9:CX	//;	XXXXXX__
	SHRQ $16, R9:DX	//;	XXXXXXXX

//; Currently all set bits in R9 represents equal bytes in a and b in the
//; 64 bytes currently being processed
//; After we NOT R9 then all set bits in R9 represents bytes that are not equal.
	NOTQ R9
//; Clear BX so we can add the POPCNTQ result directly to AX (also break
//; dependency chain as POPCNTQ has a false dependency on dest)
	XORQ BX, BX
//; Count the number of set bits and save to R9
	POPCNTQ	R9, BX
	ADDQ BX, AX

//; Update the pointers and counter then jump back to the start of the loop.
	ADDQ $64, SI 			//;						Sandy Bridge (reg im) 0.33 cycles/instructions and 1 in latency
	ADDQ $64, DI
	SUBQ $64, R8

	CMPQ R8, $64
	JNB hdHuge
//; End of hdHuge loop, there is now 63 bytes or less left to check

	JMP hdSmallEntry

//; Check one byte at a time and ADD 1 for every byte that is not equal
hdSmallnoAdd:
	ADDQ $1, DI
	ADDQ $1, SI
	SUBQ $1, R8
hdSmallEntry:
	CMPB R8, $0
	JE hdRet
hdSmall:
	MOVB (SI), CX
	CMPB CX, (DI)
	JE hdSmallnoAdd
	ADDQ $1, AX
	ADDQ $1, DI
	ADDQ $1, SI
	SUBQ $1, R8
	CMPB R8, $0
	JNE hdSmall

hdRet:
	MOVQ AX, res+48(FP)
	RET

missingPOPCNTQ:
//; Not yet implemented

missingSSE2:
//; Not yet implemented

error:
	MOVQ $-1, res+48(FP)
	RET
//; ============================================================================
//; ===[End of HammingDistance]=================================================
//; ============================================================================
