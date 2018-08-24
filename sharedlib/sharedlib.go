package main

import "github.com/xiaonanln/go-consistent"

var (
	c = consistent.New()
)

func main() {

}

//export Add
func Add(host string) {
	c.Add(host)
}

//export Remove
func Remove(host string) {
	c.Remove(host)
}

//export Hash
func Hash(key string) string {
	h, err := c.Hash(key)
	if err != nil {
		return ""
	} else {
		return h
	}
}

//export SetReplica
func SetReplica(replica int) {
	c.SetReplica(replica)
}
