package container

import (
	"hash"
	"hash/fnv"
	"math"
)

type BloomFilter struct {
	bitSet []bool
	size   int
	hashes []hash.Hash64
}

func NewBloomFilter(size int, hashCount int) *BloomFilter {
	bitSet := make([]bool, size)
	hashes := make([]hash.Hash64, hashCount)

	// Инициализация хеш-функций
	for i := 0; i < hashCount; i++ {
		hashes[i] = fnv.New64()
	}

	return &BloomFilter{
		bitSet: bitSet,
		size:   size,
		hashes: hashes,
	}
}
func (bf *BloomFilter) Add(item string) {
	for _, hashFunc := range bf.hashes {
		hashFunc.Reset()
		hashFunc.Write([]byte(item))
		index := hashFunc.Sum64() % uint64(bf.size)
		bf.bitSet[index] = true
	}
}
func (bf *BloomFilter) Check(item string) bool {
	for _, hashFunc := range bf.hashes {
		hashFunc.Reset()
		hashFunc.Write([]byte(item))
		index := hashFunc.Sum64() % uint64(bf.size)
		if !bf.bitSet[index] {
			return false
		}
	}
	return true
}

func OptimalParams(n int, p float64) (int, int) {
	m := int(math.Ceil(float64(-n) * math.Log(p) / (math.Pow(math.Log(2), 2))))
	k := int(math.Ceil(math.Log(2) * float64(m) / float64(n)))
	return m, k
}
