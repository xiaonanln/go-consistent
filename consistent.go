package consistent

import (
	"encoding/binary"
	"unsafe"

	"sort"

	"github.com/pkg/errors"
)

const DefaultReplica = 20

var (
	ErrNoHost = errors.New("no host")
)

type circleEntry struct {
	hash uint32
	host string
}
type circle []circleEntry

// Len returns the length of the circle array.
func (x circle) Len() int { return len(x) }

// Less returns true if element i is less than element j.
func (x circle) Less(i, j int) bool {
	h1, h2 := x[i].hash, x[j].hash
	if h1 != h2 {
		return h1 < h2
	}

	return x[i].host < x[j].host
}

// Swap exchanges elements i and j.
func (x circle) Swap(i, j int) { x[i], x[j] = x[j], x[i] }

type Consistent struct {
	replica int
	circle  circle
	hosts   map[string]struct{}
}

func NewConsistent() *Consistent {
	c := &Consistent{
		replica: DefaultReplica,
		hosts:   map[string]struct{}{},
	}
	return c
}

func (c *Consistent) SetReplica(replica int) {
	if replica <= 0 {
		panic(errors.Errorf("wrong replica value"))
	}

	c.replica = replica
	c.rebuildCircle()
}

func (c *Consistent) Add(host string) {
	if _, ok := c.hosts[host]; ok {
		return
	}

	c.hosts[host] = struct{}{}
	c.rebuildCircle()
}

func (c *Consistent) Remove(host string) {
	if _, ok := c.hosts[host]; !ok {
		return
	}
	delete(c.hosts, host)
	c.rebuildCircle()
}

func (c *Consistent) Hash(key string) (host string, err error) {
	if len(c.circle) == 0 {
		return "", ErrNoHost
	}

	keyhash := c.hashKey(key)
	i := sort.Search(len(c.circle), func(i int) bool {
		return c.circle[i].hash > keyhash
	})
	if i >= len(c.circle) {
		i = 0
	}

	return c.circle[i].host, nil
}

func (c *Consistent) hashHost(host string, replicaIdx uint32) uint32 {
	return hashPrefix(replicaIdx, unsafeStringToByteSlice(host))
}

func (c *Consistent) hashKey(key string) uint32 {
	return hash(unsafeStringToByteSlice(key))
}

func (c *Consistent) rebuildCircle() {
	c.circle = make(circle, 0, len(c.hosts)*c.replica)
	for host := range c.hosts {
		for i := 0; i < c.replica; i++ {
			replicaIdx := uint32(i)
			c.circle = append(c.circle, circleEntry{
				hash: c.hashHost(host, replicaIdx),
				host: host,
			})
		}
	}
	sort.Sort(c.circle)
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
