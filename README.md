# go-consistent
go consistent hashing library

It is not concurrency safe.

# Install 
```bash
go get github.com/xiaonanln/go-consistent
```

# Usage

```go
import "github.com/xiaonanln/go-consistent"

// Create Consistent Hashing
c := consistent.New()
// Hash returns ErrNoHost if there are no hosts
c.Hash("key") // returns "", consistent.ErrNoHost
// Add adds a new hash
c.Add("host1")
c.Hash("key") // returns "host1", nil 

c.Add("host2")
c.Hash("key") // returns "host1"/"host2", nil

// SetReplica changes replica. The defualt replica is 20
c.SetReplica(100) // reset replica from 20 to 100
```
# Complexity
Assuming `N` is number of hosts:  
* `Add` = `O(N)`
* `Hash` = `O(log N)`
* `SetReplica` = `O(N)`