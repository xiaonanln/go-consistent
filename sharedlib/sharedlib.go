package main

import "github.com/xiaonanln/go-consistent"

var (
	c = consistent.New()
)

func main() {

}

func Add(host string) {
	c.Add(host)
}

func Remove(host string) {
	c.Remove(host)
}

func Hash(key string) string {
	h, err := c.Hash(key)
	if err != nil {
		return ""
	} else {
		return h
	}
}

func SetReplica(replica int) {
	c.SetReplica(replica)
}
