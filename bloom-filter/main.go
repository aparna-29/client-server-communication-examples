package main

import (
	"fmt"
	"hash/fnv"
)

// BloomFilter struct
type BloomFilter struct {
	bitset []bool
	k      int // number of hash functions
	m      int // size of bitset
}

// NewBloomFilter creates a new Bloom Filter
func NewBloomFilter(size int, hashCount int) *BloomFilter {
	return &BloomFilter{
		bitset: make([]bool, size),
		k:      hashCount,
		m:      size,
	}
}

// hash1 uses FNV-1a hash
func hash1(data string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(data))
	return h.Sum32()
}

// hash2 uses FNV-1 hash (different initial parameters)
func hash2(data string) uint32 {
	h := fnv.New32()
	h.Write([]byte(data))
	return h.Sum32()
}

// combinedHash implements double hashing: h_i(x) = h1(x) + i * h2(x)
func combinedHash(data string, i int) uint32 {
	return hash1(data) + uint32(i)*hash2(data)
}

// Add inserts an item into the Bloom Filter
func (bf *BloomFilter) Add(item string) {
	for i := 0; i < bf.k; i++ {
		pos := combinedHash(item, i) % uint32(bf.m)
		bf.bitset[pos] = true
	}
}

// Contains checks if an item is possibly in the set
func (bf *BloomFilter) Contains(item string) bool {
	for i := 0; i < bf.k; i++ {
		pos := combinedHash(item, i) % uint32(bf.m)
		if !bf.bitset[pos] {
			return false // definitely not in the set
		}
	}
	return true // possibly in the set
}

func main() {
	// Example: Bloom Filter with 100 bits, 4 hash functions
	bf := NewBloomFilter(100, 4)

	// Add some items
	bf.Add("apple")
	bf.Add("banana")
	bf.Add("orange")

	// Check membership
	fmt.Println("Contains 'apple'?  ", bf.Contains("apple"))  // true (very likely)
	fmt.Println("Contains 'banana'? ", bf.Contains("banana")) // true (very likely)
	fmt.Println("Contains 'grape'?  ", bf.Contains("grape"))  // false (definitely not)
	fmt.Println("Contains 'kiwi'?   ", bf.Contains("kiwi"))   // false (definitely not)
}
