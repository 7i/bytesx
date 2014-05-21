//; func IndexNotEqual(a, b []byte) int
TEXT ·IndexNotEqual(SB),NOSPLIT,$0
	B ·indexNotEqual(SB)
   
//; func EqualThreshold(a, b []byte, t uint8) bool
TEXT ·EqualThreshold(SB),NOSPLIT,$0
	B ·equalThreshold(SB)
