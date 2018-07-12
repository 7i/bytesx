// Package bytesx implements highly optimized byte functions which extends the
// bytes package in the standard library (Currently x86 64-bit only)
package bytesx

// HammingDistance returns the Hamming distance (number if non matching bytes)
// between twostrings of the same length. HammingDistance will return -1 if the
// strings are not of the same length or if the strings are of length 0. This
// function uses sse2 instructions and the POPCNTQ instruction and will fallback
// to a pure GO implementation if any of these CPU fetures are not available.
func HammingDistance(a, b []byte) int

// IndexNotEqual returns the index of the first non matching byte between a and
// b, or -1 if a and b are equal untill the shortest of the two.
func IndexNotEqual(a, b []byte) int

// EqualThreshold returns true if b does not differ in value more than t from
// the corresponding byte in a.
// t may take any value from 0 to 255 where 0 is exact match and 255 will match
// any string. If t is 1 and a is "MNO" and b is "LNP" than EqualThreshold will
// return true while it will return false if b is "LNQ" or "KNO". The equality
// check is only made untill the shortest of a and b.
func EqualThreshold(a, b []byte, t uint8) bool

// Fallback for 386 and arm
func indexNotEqual(a, b []byte) int {
	if &a[0] == &b[0] {
		return -1
	}
	var min int
	if len(a) < len(b) {
		min = len(a)
	} else {
		min = len(b)
	}
	for i := 0; i < min; i++ {
		if a[i] == b[i] {
			continue
		}
		return i
	}
	return -1
}

// Fallback for 386 and arm
func equalThreshold(a, b []byte, t uint8) bool {
	if &a[0] == &b[0] {
		return true
	}
	var min int
	if len(a) < len(b) {
		min = len(a)
	} else {
		min = len(b)
	}
	for i := 0; i < min; i++ {
		ai := int(a[i])
		ti := int(t)
		bi := int(b[i])
		if bi < ai+ti && bi > ai-ti {
			continue
		}
		return false
	}
	return true
}
