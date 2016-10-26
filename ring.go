// Package ring is a Fast Consistent Hashing module
// It is based on a paper by John Lamping and Eric Veach called "A Fast, Minimal Memory, Consistent Hash Algorithm"
// The paper can be found here: http://arxiv.org/pdf/1406.2294v1.pdf
package ring

// package imports
import (
	"fmt"
	"hash/fnv"
	"sort"
	"sync"

	jump "github.com/dgryski/go-jump"
)

// --------------------
//      Interfaces
// --------------------

// Node is an interface representing a physical host. Each node has a host, a capacity and a hash.
type Node interface {

	// Returns the host for the node.
	GetHost() string

	// Returns the capacity of the node. This number determines how many virtual nodes belong to the host.
	GetSize() int

	// Returns the hash of the node. This 64-bit number symbolizes where a node falls on the ring.
	GetHash() uint64
}

// Ring is the main interface for this package. It comprises of methods used to hash keys into buckets which
// will be evenly divided among all virtual nodes in the ring.
// All values are hashed using the FNV algorithm into an unsigned 64-bit integer. The Jump Hash
// algorithm then determines which bucket a hash falls into.
type Ring interface {

	// Adds a host to the ring. The first arg
	Add(host string, size int)

	// Returns the size of the ring. Virtual nodes are included.
	Size() int

	// Returns a node for the given bucket number
	GetNode(data []byte) Node
}

// --------------------
//      Objects
// --------------------

// Node struct
type node struct {
	host string
	size int
	hash uint64
}

// Getter for a Node's host
func (n node) GetHost() string {
	return n.host
}

// Getter for a Node's size
func (n node) GetSize() int {
	return n.size
}

// Getter for a Node's hash
func (n node) GetHash() uint64 {
	return n.hash
}

// NewNode creates a new Node with a hostname and a capacity.
func NewNode(host string, size int) Node {
	return node{host: host, size: size}
}

// --------------------
//      Hash Ring
// --------------------

type nodeList []node

type hashRing struct {
	sync.RWMutex
	nodes nodeList
}

// Len is the number of elements in the collection.
func (h nodeList) Len() int {
	return len(h)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (h nodeList) Less(i, j int) bool {
	return h[i].hash < h[j].hash
}

// Swap swaps the elements with indexes i and j.
func (h nodeList) Swap(i, j int) {
	h[i], h[j] = h[j], h[i]
}

// sorts the existing nodes by hash
func (h nodeList) sort() {
	sort.Sort(h)
}

// adds a host (+virtual hosts to the ring)
func (h *hashRing) Add(host string, size int) {
	h.Lock()
	defer h.Unlock()
	var hasher = fnv.New64a()
	hlen := len(h.nodes)
	cap := hlen + size

	// resize node list
	nodes := make([]node, cap)
	copy(nodes, h.nodes)
	h.nodes = nodes

	// insert new nodes at the end
	for i := hlen; i < cap; i++ {
		// hash: 0:localhost:7000:0
		// adding the index at the start and end seemed to give better distribution...
		hasher.Write([]byte(fmt.Sprint(i, ":", host, ":", i)))

		// hash value
		value := hasher.Sum64()

		// create node
		n := node{hash: value, host: host, size: size}

		// insert node
		h.nodes[i] = n

		// reset hash
		hasher.Reset()
	}

	// sort nodes around ring based on hash
	h.nodes.sort()
}

// returns the size of the ring
func (h *hashRing) Size() int {
	h.RLock()
	defer h.RUnlock()
	return len(h.nodes)
}

// returns a particular index
func (h *hashRing) GetNode(data []byte) Node {
	h.RLock()
	defer h.RUnlock()
	return h.nodes[h.calculateJumpHash(hash(data))]
}

// returns a particular index
func (h *hashRing) calculateJumpHash(hash uint64) int {
	h.RLock()
	defer h.RUnlock()
	return int(jump.Hash(hash, len(h.nodes)))
}

// NewHashRing creates a new hash ring.
func NewHashRing() Ring {
	return &hashRing{nodes: make([]node, 0, 16)}
}

// CalculateBucketGivenSize calculates a Jump hash for the key provided
func CalculateBucketGivenSize(data []byte, size int) int {
	var hasher = fnv.New64a()
	hasher.Write(data)
	return int(jump.Hash(hasher.Sum64(), size))
}

// FNV-64a hash
func hash(data []byte) uint64 {
	var hasher = fnv.New64a()
	hasher.Write(data)
	return hasher.Sum64()
}
