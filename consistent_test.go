package consistent

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"unsafe"
)

func TestConsistentBasic(t *testing.T) {
	c := New()
	_, err := c.Hash("abc")
	if err != ErrNoHost {
		t.Fatalf("should returns ErrNoHost")
	}

	c.Add("host1")
	host, err := c.Hash("anykey")
	if host != "host1" {
		t.Fatalf("should returns host1")
	}

	c.Remove("host1")
	_, err = c.Hash("abc")
	if err != ErrNoHost {
		t.Fatalf("should returns ErrNoHost")
	}

	c1 := New()
	c2 := New()
	for i := 0; i < 10; i++ {
		c1.Add(fmt.Sprintf("host%d", i))
		c2.Add(fmt.Sprintf("host%d", 10-i-1))
	}

	strbuf := make([]byte, 256)

	for i := 0; i < 1000; i++ {
		keylen := rand.Intn(16) + 16
		rand.Read(strbuf[:keylen])
		key := string(strbuf[:keylen])
		host1, err := c1.Hash(key)
		if err != nil {
			t.Fatal(err)
		}
		host2, err := c2.Hash(key)
		if err != nil {
			t.Fatal(err)
		}
		if host1 != host2 {
			t.FailNow()
		}
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

func BenchmarkConsistent_AddRemove(b *testing.B) {
	c := New()
	for i := 0; i < 10; i++ {
		c.Add(strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		c.Add("test")
		c.Remove("test")
	}
}

func BenchmarkConsistent_Hash(b *testing.B) {
	c := New()
	c.SetReplica(100)
	for i := 0; i < 10; i++ {
		c.Add(strconv.Itoa(i))
	}
	for i := 0; i < b.N; i++ {
		c.Hash(strconv.Itoa(1000 + i))
	}
}
