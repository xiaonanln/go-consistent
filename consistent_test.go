package consistent

import (
	"testing"
	"unsafe"
)

func TestConsistentBasic(t *testing.T) {
	c := NewConsistent()
	_, err := c.Hash("abc")
	if err != ErrNoHost {
		t.Fatalf("should returns ErrNoHost")
	}

	c.Add("host1")
	host, err := c.Hash("anykey")
	if host != "host1" {
		t.Fatalf("should returns host1")
	}

}

func BenchmarkStringToSlice1(b *testing.B) {
	s := "some string LLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL"
	tl := 0
	for i := 0; i < b.N; i++ {
		tl += len([]byte(s))
	}
	if tl != len(s)*b.N {
		b.Fatalf("wrong length")
	}
}

func BenchmarkStringToSlice2(b *testing.B) {
	s := "some string LLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLLL"
	tl := 0
	for i := 0; i < b.N; i++ {
		tl += len(*((*[]byte)(unsafe.Pointer(&s))))
	}
	if tl != len(s)*b.N {
		b.Fatalf("wrong length")
	}
}