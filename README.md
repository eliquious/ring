# Ring [![Coverage Status](https://coveralls.io/repos/eliquious/ring/badge.png)](https://coveralls.io/r/eliquious/ring)


    go get -u github.com/eliquious/ring

Package `ring` is a very fast consistent hashing module. It is based on a [paper](http://arxiv.org/pdf/1406.2294v1.pdf) by John
Lamping and Eric Veach called "A Fast, Minimal Memory, Consistent Hash Algorithm".

The hash ring is thread safe so it can be used by multiple goroutines.

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

	// Add adds a host to the ring.
	Add(host string, size int)

	// Size Returns the size of the ring. Virtual nodes are included.
	Size() int

	// GetNode returns a node for the given input
	GetNode(data []byte) Node
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
BenchmarkGetNode_5_Nodes    	10000000	       124 ns/op
BenchmarkGetNode_25_Nodes   	10000000	       139 ns/op
BenchmarkGetNode_100_Nodes  	10000000	       151 ns/op
BenchmarkGetNode_1000_Nodes 	10000000	       170 ns/op
BenchmarkGetNode_10000_Nodes	10000000	       192 ns/op
```
