//; func IndexNotEqual(a, b []byte) int
TEXT ·IndexNotEqual(SB),NOSPLIT,$0
	JMP ·indexNotEqual(SB)

//; func EqualThreshold(a, b []byte, t uint8) bool
TEXT ·EqualThreshold(SB),NOSPLIT,$0
	JMP ·equalThreshold(SB)
