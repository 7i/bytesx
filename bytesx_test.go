package bytesx_test

import (
	"fmt"
	"github.com/7i/bytesx"
	"syscall"
	"testing"
	"unsafe"
)

var test1 = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")
var test2 = []byte("AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA")

func TestIndexNotEqual(t *testing.T) {
	// Test special cases
	if -1 != bytesx.IndexNotEqual([]byte(""), []byte("X")) {
		t.Errorf("IndexNotEqual test failed. String1: \"\" String2: \"X\"")
	}
	if -1 != bytesx.IndexNotEqual([]byte(""), []byte("")) {
		t.Errorf("IndexNotEqual test failed. String1: \"\" String2: \"\"")
	}
	if -1 != bytesx.IndexNotEqual([]byte("a"), []byte("aX")) {
		t.Errorf("IndexNotEqual test failed. String1: \"a\" String2: \"aX\"")
	}
	if -1 != bytesx.IndexNotEqual([]byte("aX"), []byte("a")) {
		t.Errorf("IndexNotEqual test failed. String1: \"aX\" String2: \"a\"")
	}

	// Test all length of data from 1 to 128 bytes with a non matching byte in every
	// possibe possition and all length of strings from 1 to 128 where the data is
	// the same.
	for i := range test1 {
		for j := 0; j < i; j++ {
			test2[j] = 'X'
			got := bytesx.IndexNotEqual(test1[:i], test2[:i])
			test2[j] = 'A'
			if got != j {
				t.Errorf("IndexNotEqual test failed. \nGot:%d\nString1: %d \"A\"s\nString2: %d \"A\"s plus a \"X\" at possition %d \nExpected: %d ", got, i, i-1, j, j)
			}
		}
		got := bytesx.IndexNotEqual(test1[:i], test2[:i])
		if got != -1 {
			t.Errorf("IndexNotEqual test failed. Got:%d\nString1: %d \"A\"s\nString2: %d \"A\"s\nExpected: -1 ", got, i, i)
		}
	}
}

func TestEqualThresholdNearPageBoundary(t *testing.T) {
	pagesize := syscall.Getpagesize()
	b := make([]byte, 4*pagesize)
	i := pagesize
	for ; uintptr(unsafe.Pointer(&b[i]))%uintptr(pagesize) != 0; i++ {
	}
	syscall.Mprotect(b[i-pagesize:i], 0)
	syscall.Mprotect(b[i+pagesize:i+2*pagesize], 0)
	defer syscall.Mprotect(b[i-pagesize:i], syscall.PROT_READ|syscall.PROT_WRITE)
	defer syscall.Mprotect(b[i+pagesize:i+2*pagesize], syscall.PROT_READ|syscall.PROT_WRITE)

	// both of these should fault
	//pagesize += int(b[i-1])
	//pagesize += int(b[i+pagesize])

	for j := 0; j < pagesize; j++ {
		b[i+j] = 'A'
	}
	for j := 0; j <= pagesize; j++ {
		bytesx.EqualThreshold(b[i:i+j], b[i+pagesize-j:i+pagesize], 0)
		bytesx.EqualThreshold(b[i+pagesize-j:i+pagesize], b[i:i+j], 0)
	}
}

func BenchmarkEqualThreshold0(b *testing.B) {
	var buf [4]byte
	buf1 := buf[0:0]
	buf2 := buf[1:1]
	for i := 0; i < b.N; i++ {
		eq := bytesx.EqualThreshold(buf1, buf2, 1)
		if !eq {
			b.Fatal("bad equal")
		}
	}
}

var bmbuf []byte

func BenchmarkEqualThreshold1(b *testing.B)    { bmEqualThreshold(b, bytesx.EqualThreshold, 1) }
func BenchmarkEqualThreshold6(b *testing.B)    { bmEqualThreshold(b, bytesx.EqualThreshold, 6) }
func BenchmarkEqualThreshold9(b *testing.B)    { bmEqualThreshold(b, bytesx.EqualThreshold, 9) }
func BenchmarkEqualThreshold15(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 15) }
func BenchmarkEqualThreshold16(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 16) }
func BenchmarkEqualThreshold20(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 20) }
func BenchmarkEqualThreshold32(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 32) }
func BenchmarkEqualThreshold64(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 64) }
func BenchmarkEqualThreshold128(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 128) }
func BenchmarkEqualThreshold256(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 256) }
func BenchmarkEqualThreshold512(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 512) }
func BenchmarkEqualThreshold1K(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 1<<10) }
func BenchmarkEqualThreshold2K(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 2<<10) }
func BenchmarkEqualThreshold4K(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 4<<10) }
func BenchmarkEqualThreshold8K(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 8<<10) }
func BenchmarkEqualThreshold16K(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 16<<10) }
func BenchmarkEqualThreshold64K(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 64<<10) }
func BenchmarkEqualThreshold1M(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 1<<20) }
func BenchmarkEqualThreshold2M(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 2<<20) }
func BenchmarkEqualThreshold4M(b *testing.B)   { bmEqualThreshold(b, bytesx.EqualThreshold, 4<<20) }
func BenchmarkEqualThreshold64M(b *testing.B)  { bmEqualThreshold(b, bytesx.EqualThreshold, 64<<20) }
func BenchmarkEqualThreshold128M(b *testing.B) { bmEqualThreshold(b, bytesx.EqualThreshold, 128<<20) }
func BenchmarkEqualThreshold512M(b *testing.B) { bmEqualThreshold(b, bytesx.EqualThreshold, 512<<20) }

func bmEqualThreshold(b *testing.B, equal func([]byte, []byte, uint8) bool, n int) {
	if len(bmbuf) < 2*n {
		bmbuf = make([]byte, 2*n)
	}
	b.SetBytes(int64(n))
	buf1 := bmbuf[0:n]
	buf2 := bmbuf[n : 2*n]
	buf1[n-1] = 'x'
	buf2[n-1] = 'x'
	for i := 0; i < b.N; i++ {
		eq := equal(buf1, buf2, 1)
		if !eq {
			b.Fatal("bad equal threshold")
		}
	}
	buf1[n-1] = '\x00'
	buf2[n-1] = '\x00'
}

func TestEqualThreshold(t *testing.T) {
	// Temp test
	bmbuf = make([]byte, 2*(1<<32)+1)
	// End Temp test

	return
	size := 128
	if testing.Short() {
		size = 32
	}

	// Test special cases
	if true != bytesx.EqualThreshold([]byte(""), []byte("X"), 0) {
		t.Errorf("EqualThreshold test failed. String1: \"\" String2: \"X\" Threshold: 0")
	}
	if true != bytesx.EqualThreshold([]byte("X"), []byte(""), 0) {
		t.Errorf("EqualThreshold test failed. String1: \"\" String2: \"\"")
	}
	if true != bytesx.EqualThreshold([]byte(""), []byte(""), 0) {
		t.Errorf("EqualThreshold test failed. String1: \"a\" String2: \"aX\"")
	}

	fmt.Println("Testing all threshold values for strings up to 128 bytes")
	fmt.Println("Total number of tests: 2'147'483'648")
	fmt.Println("This will take 1-2 min")
	for i := 1; i < size; i++ {
		for a := 0; a < 256; a++ {
			for b := 0; b < 256; b++ {
				for th := 0; th < 256; th++ {
					// set test1 and test2 to the same data at all positions
					for neq := bytesx.IndexNotEqual(test1, test2); -1 != neq; neq = bytesx.IndexNotEqual(test1, test2) {
						test1[neq] = 'A'
						test2[neq] = 'A'
					}

					test1[i-1] = byte(a)
					test2[i-1] = byte(b)

					diff := 0
					if a > b {
						diff = a - b
					} else {
						diff = b - a
					}

					got := bytesx.EqualThreshold(test1[:i], test2[:i], byte(th))
					ans := th >= diff

					if ans != got {
						t.Errorf("\nEqualThreshold test failed.\nGot:  %v\nAns:  %v\nStr1: %s\nStr2: %s\nTH:   %d\nDiff: %d\nLen:  %d", got, ans, test1[:i], test2[:i], byte(th), diff, i)
					}
				}
			}
		}
	}
}

//// Benchmark Hamming distance
func BenchmarkHammingDistance2x2G(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<31) }
func BenchmarkHammingDistance2x1G(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<30) }
func BenchmarkHammingDistance2x512M(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<29) }
func BenchmarkHammingDistance2x256M(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<28) }
func BenchmarkHammingDistance2x128M(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<27) }
func BenchmarkHammingDistance2x64M(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<26) }
func BenchmarkHammingDistance2x32M(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<25) }
func BenchmarkHammingDistance2x16M(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<24) }
func BenchmarkHammingDistance2x8M(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<23) }
func BenchmarkHammingDistance2x4M(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<22) }
func BenchmarkHammingDistance2x2M(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<21) }
func BenchmarkHammingDistance2x1M(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<20) }
func BenchmarkHammingDistance2x512K(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<19) }
func BenchmarkHammingDistance2x256K(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<18) }
func BenchmarkHammingDistance2x128K(b *testing.B) { bmHammingDistance(b, bytesx.HammingDistance, 1<<17) }
func BenchmarkHammingDistance2x64K(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<16) }
func BenchmarkHammingDistance2x32K(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<15) }
func BenchmarkHammingDistance2x16K(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<14) }
func BenchmarkHammingDistance2x8K(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<13) }
func BenchmarkHammingDistance2x4K(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<12) }
func BenchmarkHammingDistance2x2K(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<11) }
func BenchmarkHammingDistance2x1K(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 1<<10) }
func BenchmarkHammingDistance2x512(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 1<<9) }
func BenchmarkHammingDistance2x256(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 256) }
func BenchmarkHammingDistance2x128(b *testing.B)  { bmHammingDistance(b, bytesx.HammingDistance, 196) }
func BenchmarkHammingDistance2x64(b *testing.B)   { bmHammingDistance(b, bytesx.HammingDistance, 64) }
func BenchmarkHammingDistance2x8(b *testing.B)    { bmHammingDistance(b, bytesx.HammingDistance, 8) }
func BenchmarkHammingDistance2x2(b *testing.B)    { bmHammingDistance(b, bytesx.HammingDistance, 2) }

func bmHammingDistance(b *testing.B, hd func([]byte, []byte) int, n int) {
	if len(bmbuf) < 2*n {
		bmbuf = make([]byte, 2*n)
	}
	b.SetBytes(int64(2 * n))
	buf1 := bmbuf[0:n]
	buf2 := bmbuf[n : 2*n]
	buf1[0] = 'y'
	buf1[n-1] = 'y'
	buf2[0] = 'x'
	buf2[n-1] = 'x'
	for i := 0; i < b.N; i++ {
		count := hd(buf1, buf2)
		if count != 2 {
			b.Fatal("bad equal threshold")
		}
	}
	buf1[0] = '\x00'
	buf1[n-1] = '\x00'
	buf2[0] = '\x00'
	buf2[n-1] = '\x00'
}
