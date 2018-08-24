package main

import "C"

import "github.com/xiaonanln/go-consistent"

var (
	c = consistent.New()
)

func main() {

}

//export Add
func Add(host *C.char) {
	c.Add(C.GoString(host))
}

//export Remove
func Remove(host *C.char) {
	c.Remove(C.GoString(host))
}

//export Hash
func Hash(key *C.char) *C.char {
	h, err := c.Hash(C.GoString(key))
	if err != nil {
		return C.CString("")
	} else {
		return C.CString(h)
	}
}

//export SetReplica
func SetReplica(replica int) {
	c.SetReplica(replica)
}
