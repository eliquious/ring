# Ring [![Coverage Status](https://coveralls.io/repos/eliquious/ring/badge.png)](https://coveralls.io/r/eliquious/ring)


    import "github.com/swiftkick-io/ring"

Package ring is a very fast consistent hashing module. It is based on a [paper](http://arxiv.org/pdf/1406.2294v1.pdf) by John
Lamping and Eric Veach called "A Fast, Minimal Memory, Consistent Hash Algorithm".

## Usage

#### type Node

```go
type Node interface {

	// Returns the host for the node.
	GetHost() string

	// Returns the capacity of the node. This number determines how many virtual nodes belong to the host.
	GetSize() int

	// Returns the hash of the node. This 64-bit number symbolizes where a node falls on the ring.
	GetHash() uint64
}
```

Node is an interface representing a physical host. Each node has a host, a
capacity and a hash.

#### func  NewNode

```go
func NewNode(host string, size int) Node
```
NewNode creates a new Node with a hostname and a capacity.

#### type Ring

```go
type Ring interface {

	// Adds a host to the ring. The first arg
	Add(host string, size int)

	// Determines the bucket of an unsigned 64-bit integer
	FindBucket(key uint64) int

	// Hashes the bytes given with FNV and then returns the result of FindBucket(key uint64)
	FindBucketWithBytes(data []byte) int

	// Hashes the string given with FNV and then returns the result of FindBucket(key uint64)
	FindBucketWithString(data string) int

	// Finds a bucket for a given key based on the size of the ring given.
	FindBucketGivenSize(key uint64, size int) int

	// Hashes the data using FNV
	Hash(data []byte) uint64

	// Returns the size of the ring. Virtual nodes are included.
	Size() int

	// Returns a node for the given bucket number
	GetNode(index int) Node
}
```

Ring is the main interface for this package. It comprises of methods used to
hash keys into buckets which will be evenly divided among all virtual nodes in
the ring. All values are hashed using the FNV algorithm into an unsigned 64-bit
integer. The Jump Hash algorithm then determines which bucket a hash falls into.

#### func  NewHashRing

```go
func NewHashRing() Ring
```
NewHashRing creates a new hash ring.

## Benchmarks

The number implies the total virtual nodes in the hash ring.

```
Benchmark_5_NodeHashRing	100000000	        23.7 ns/op
Benchmark_25_NodeHashRing	50000000	        45.0 ns/op
Benchmark_100_NodeHashRing	50000000	        58.4 ns/op
Benchmark_1000_NodeHashRing	20000000	        81.4 ns/op
```
