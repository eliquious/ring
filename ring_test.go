package ring

import (
	"fmt"
	"testing"

	"github.com/GaryBoone/GoStats/stats"
	jump "github.com/dgryski/go-jump"
	"github.com/stretchr/testify/assert"
)

// tests the GetHost function on the Node interface
func TestNodeGetHost(t *testing.T) {
	host := "localhost:7000"
	n := NewNode(host, 5)

	// validate host name
	assert.Equal(t, host, n.GetHost())
}

// tests the ring size
func TestNodeGetSize(t *testing.T) {
	host := "localhost:7000"

	// create a new host with 5 virtual nodes
	n := NewNode(host, 5)

	// size of new node should be 5
	assert.Equal(t, 5, n.GetSize())
}

// validates nodes (+ virtual nodes have been created)
func TestNodeAddNode(t *testing.T) {
	host := "localhost:7000"

	// create new ring
	r := NewHashRing()

	// add 5 virtual nodes
	r.Add(host, 5)

	// there should now be 5 virtual nodes
	assert.Equal(t, 5, r.Size())
}

// make sure nodes have been sorted by hash
func TestNodesHaveBeenSorted(t *testing.T) {
	host := "localhost:7000"

	// create new ring
	r := NewHashRing().(*hashRing)

	// add 5 virtual nodes
	r.Add(host, 5)

	// test sort order of hash values
	last := uint64(0)
	for i := 0; i < r.Size(); i++ {
		assert.True(t, r.nodes[i].GetHash() > last)
		last = r.nodes[i].GetHash()
	}
}

// ensures all virtual nodes are distributed evenly within e=0.0001
func TestVirtualNodeDistribution(t *testing.T) {
	r := NewHashRing().(*hashRing)

	// add 5 virtual nodes
	r.Add("localhost:7000", 20)
	r.Add("localhost:7001", 20)
	r.Add("localhost:7002", 20)
	r.Add("localhost:7003", 20)
	r.Add("localhost:7004", 20)

	var d stats.Stats
	counts := make([]int, r.Size())
	var COUNT int = 10e6
	for i := 0; i < COUNT; i++ {
		pos := r.calculateJumpHash(uint64(i))
		counts[pos]++
	}

	var avg = float64(COUNT) / float64(r.Size())
	for _, node := range counts {
		d.Update(float64(node))
	}
	// fmt.Println("StDev: ", d.PopulationStandardDeviation())
	// fmt.Println("Var: ", d.PopulationVariance())
	// fmt.Println("Mean: ", d.Mean())
	assert.InEpsilon(t, avg, d.Mean(), 0.0001)
}

// ensures all physical nodes are distributed evenly within e=0.0001
func TestNodeDistribution(t *testing.T) {
	r := NewHashRing().(*hashRing)

	// add 5 virtual nodes
	r.Add("localhost:7000", 20)
	r.Add("localhost:7001", 20)
	r.Add("localhost:7002", 20)
	r.Add("localhost:7003", 20)
	r.Add("localhost:7004", 20)

	var d stats.Stats
	nodes := make(map[string]int)
	nodes["localhost:7000"] = 0
	nodes["localhost:7001"] = 0
	nodes["localhost:7002"] = 0
	nodes["localhost:7003"] = 0
	nodes["localhost:7004"] = 0

	var COUNT int = 10e6
	size := r.Size()
	for i := 0; i < COUNT; i++ {
		pos := int(jump.Hash(uint64(i), size))
		nodes[r.nodes[pos].GetHost()]++
	}

	var avg = float64(COUNT) / float64(5)
	for _, value := range nodes {
		// fmt.Println(key, value)
		d.Update(float64(value))
	}
	// fmt.Println("StDev: ", d.PopulationStandardDeviation())
	// fmt.Println("Var: ", d.PopulationVariance())
	// fmt.Println("Mean: ", d.Mean())
	// fmt.Println("Expected: ", avg)
	assert.InEpsilon(t, avg, d.Mean(), 0.0001)
}

// closure function for benchmarking multiple clusters
func baselineBenchmark(hosts, vnodes int) func(b *testing.B) {
	ring := NewHashRing().(*hashRing)
	var startPort = 7000
	for i := startPort; i < hosts+startPort; i++ {
		ring.Add(fmt.Sprint("localhost:", i), vnodes)
	}

	return func(b *testing.B) {
		b.ResetTimer()

		// use the ring hash a number
		for n := 0; n < b.N; n++ {
			ring.calculateJumpHash(uint64(n))
		}
	}
}

// 5 Nodes
func BenchmarkGetNode_5_Nodes(b *testing.B) {
	baselineBenchmark(5, 1)(b)
}

// 5 Nodes with 5 Virtual Nodes each
func BenchmarkGetNode_25_Nodes(b *testing.B) {
	baselineBenchmark(5, 5)(b)
}

// 20 Nodes with 5 Virtual Nodes each
func BenchmarkGetNode_100_Nodes(b *testing.B) {
	baselineBenchmark(20, 5)(b)
}

// 250 Nodes with 5 Virtual Nodes each
func BenchmarkGetNode_1000_Nodes(b *testing.B) {
	baselineBenchmark(250, 5)(b)
}

// 2500 Nodes with 5 Virtual Nodes each
func BenchmarkGetNode_10000_Nodes(b *testing.B) {
	baselineBenchmark(2500, 5)(b)
}

func TestHashing(t *testing.T) {
	var count = 50
	lastNode := count + 1
	data := ([]byte("input"))
	for size := count; size > 1; size-- {
		assert.True(t, CalculateBucketGivenSize(data, size) <= lastNode)
	}
}

func TestHashCorrectness(t *testing.T) {
	r := NewHashRing().(*hashRing)

	// // add 5 virtual nodes
	r.Add("localhost:7000", 20)
	r.Add("localhost:7001", 20)
	r.Add("localhost:7002", 20)
	r.Add("localhost:7003", 20)
	r.Add("localhost:7004", 20)

	// The bucket should not change until the ring size
	// decreases enough to where the bucket no longer exists.
	// Then the hashed value should be remapped to another bucket.
	var bucket int
	for i := uint64(0); i < 10; i++ {
		bucket = r.calculateJumpHash(i)

		// Bucket should not change wilth the ring is
		// larger than the bucket index
		for j := (r.Size()); j > bucket; j-- {
			assert.Equal(t, int(bucket), int(jump.Hash(i, j)))
		}

		// Make sure the bucket is remapped after the ring size no
		// longer includes the bucket
		assert.NotEqual(t, bucket, jump.Hash(i, bucket))
	}
}

func TestFindBucketWithBytesAndString(t *testing.T) {

	// create new ring
	r := NewHashRing().(*hashRing)

	// add 100 virtual nodes
	r.Add("localhost:7000", 20)
	r.Add("localhost:7001", 20)
	r.Add("localhost:7002", 20)
	r.Add("localhost:7003", 20)
	r.Add("localhost:7004", 20)

	// create byte array
	data := []byte("golang")

	// Find bucket
	expected := r.nodes[r.calculateJumpHash(hash(data))]

	// ensure the bucket is correct for byte arrays
	assert.Equal(t, expected, r.GetNode(data))
}
