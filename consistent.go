package consistent

import (
	"encoding/binary"
	"unsafe"

	"github.com/pkg/errors"
)

const DefaultReplica = 20

var (
	ErrNoHost = errors.New("no host")
)

type hashentry struct {
	hash       uint32
	host       string
	replicaIdx uint32
}
type hashlist []struct{}

// Len returns the length of the hashlist array.
func (x hashlist) Len() int { return len(x) }

// Less returns true if element i is less than element j.
func (x hashlist) Less(i, j int) bool { return x[i] < x[j] }

// Swap exchanges elements i and j.
func (x hashlist) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type Consistent struct {
	replica      int
	sortedHashes hashlist
	host2hash    map[string]uint32
	hash2host    map[uint32]string
}

func NewConsistent() *Consistent {
	c := &Consistent{
		replica: DefaultReplica,
	}
	return c
}

func (c *Consistent) SetReplica(replica int) {
	if replica <= 0 {
		panic(errors.Errorf("wrong replica value"))
	}

	c.replica = replica
}

func (c *Consistent) Add(host string) {

}

func (c *Consistent) Remove(host string) {

}

func (c *Consistent) Hash(key string) (host string, err error) {
	keyhash := c.hashKey(key)
	idx = c.search(keyhash)
	return "", ErrNoHost
}

func (c *Consistent) hashHost(host string, replicaIdx uint32) uint32 {
	return hashPrefix(replicaIdx, unsafeStringToByteSlice(host))
}

func (c *Consistent) hashKey(key string) uint32 {
	return hash(unsafeStringToByteSlice(key))
}

func unsafeStringToByteSlice(s string) []byte {
	return *((*[]byte)(unsafe.Pointer(&s)))
}

// hash return hash of the given data.
func hash(data []byte) uint32 {
	// Similar to murmur hash
	const (
		m = uint32(0xc6a4a793)
		r = uint32(24)
	)
	var (
		h = uint32(0xbc9f1d34) ^ (uint32(len(data)) * m)
		i int
	)

	for n := len(data) - len(data)%4; i < n; i += 4 {
		h += binary.LittleEndian.Uint32(data[i:])
		h *= m
		h ^= (h >> 16)
	}

	switch len(data) - i {
	default:
		panic("not reached")
	case 3:
		h += uint32(data[i+2]) << 16
		fallthrough
	case 2:
		h += uint32(data[i+1]) << 8
		fallthrough
	case 1:
		h += uint32(data[i])
		h *= m
		h ^= (h >> r)
	case 0:
	}

	return h
}

// hash return hash of the given data.
func hashPrefix(prefix uint32, data []byte) uint32 {
	// Similar to murmur hash
	const (
		m = uint32(0xc6a4a793)
		r = uint32(24)
	)
	var (
		h = uint32(0xbc9f1d34) ^ (uint32(len(data)) * m)
		i int
	)

	h += prefix
	h *= m
	h ^= (h >> 16)

	for n := len(data) - len(data)%4; i < n; i += 4 {
		h += binary.LittleEndian.Uint32(data[i:])
		h *= m
		h ^= (h >> 16)
	}

	switch len(data) - i {
	default:
		panic("not reached")
	case 3:
		h += uint32(data[i+2]) << 16
		fallthrough
	case 2:
		h += uint32(data[i+1]) << 8
		fallthrough
	case 1:
		h += uint32(data[i])
		h *= m
		h ^= (h >> r)
	case 0:
	}

	return h
}
