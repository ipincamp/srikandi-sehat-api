package bloomfilter

import (
	"fmt"
	"hash/fnv"
	"math"
	"sync"

	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
)

// This compile-time check ensures our struct implements the port interface.
var _ ports.UserCache = (*BloomUserCache)(nil)

// hashFunc defines the signature for our hash functions.
type hashFunc func(data []byte) uint64

// BloomUserCache implements the ports.UserCache interface using a Bloom Filter.
// It is safe for concurrent use.
type BloomUserCache struct {
	m         uint64       // Number of bits in the filter
	k         uint         // Number of hash functions
	bits      []uint64     // The bit array (using uint64 for efficiency)
	hashFuncs []hashFunc   // The k hash functions
	mu        sync.RWMutex // Mutex to protect the bit array
}

// CalculateParams estimates m (bits) and k (hash functions) for a given
// n (expected items) and p (false positive rate).
// m = - (n * ln(p)) / (ln(2)^2)
// k = (m / n) * ln(2)
func CalculateParams(n uint64, p float64) (uint64, uint) {
	if n == 0 {
		n = 1000 // Default to 1000 items if 0 is passed
	}
	if p <= 0 || p >= 1 {
		p = 0.001 // Default to 0.1%
	}

	m := math.Ceil(-1 * (float64(n) * math.Log(p)) / math.Pow(math.Log(2), 2))
	k := math.Ceil((m / float64(n)) * math.Log(2))

	// Ensure k is at least 1
	if k < 1 {
		k = 1
	}

	// m bits, but we use uint64 array, so round m up to nearest 64
	mRounded := (uint64(m) + 63) &^ 63

	return mRounded, uint(k)
}

// NewBloomUserCache creates a new Bloom filter optimized for user email caching.
// It uses FNV-1a hashing with different seeds to create k hash functions.
func NewBloomUserCache(m uint64, k uint) *BloomUserCache {
	if k < 1 {
		k = 1
	}
	// Ensure m is a multiple of 64 for our uint64 array
	if m%64 != 0 {
		m = (m + 63) &^ 63
	}
	if m == 0 {
		m = 64
	}

	// We use uint64, so the array size is m / 64
	bitArraySize := m / 64
	bits := make([]uint64, bitArraySize)
	hashFuncs := make([]hashFunc, k)

	for i := uint(0); i < k; i++ {
		// Create k different hash functions by seeding FNV
		seed := uint64(i)
		hashFuncs[i] = func(data []byte) uint64 {
			h := fnv.New64a()
			// Add the seed to the data to create a different hash
			h.Write(data)
			h.Write([]byte(fmt.Sprintf("%d", seed)))
			return h.Sum64()
		}
	}

	return &BloomUserCache{
		m:         m,
		k:         k,
		bits:      bits,
		hashFuncs: hashFuncs,
		// mu is automatically initialized
	}
}

// set sets the bit at the given index to 1.
func (b *BloomUserCache) set(idx uint64) {
	arrayIdx := idx / 64              // Which uint64 in the array
	bitIdx := idx % 64                // Which bit in that uint64
	b.bits[arrayIdx] |= (1 << bitIdx) // Set the bit
}

// test checks if the bit at the given index is 1.
func (b *BloomUserCache) test(idx uint64) bool {
	arrayIdx := idx / 64
	bitIdx := idx % 64
	return (b.bits[arrayIdx] & (1 << bitIdx)) != 0
}

// Add adds an email to the filter.
func (b *BloomUserCache) Add(email string) {
	data := []byte(email)
	b.mu.Lock()
	defer b.mu.Unlock()

	for i := uint(0); i < b.k; i++ {
		hash := b.hashFuncs[i](data)
		idx := hash % b.m
		b.set(idx)
	}
}

// Test checks if an email is *possibly* in the filter.
// Returns false if the email is *definitely not* present.
// Returns true if the email is *possibly* present.
func (b *BloomUserCache) Test(email string) bool {
	data := []byte(email)
	b.mu.RLock() // Use Read-lock for checking
	defer b.mu.RUnlock()

	for i := uint(0); i < b.k; i++ {
		hash := b.hashFuncs[i](data)
		idx := hash % b.m
		if !b.test(idx) {
			// If any bit is 0, it's definitely not in the set
			return false
		}
	}

	// All bits were 1, so it *might* be in the set
	return true
}

// Populate efficiently adds a large list of emails,
// acquiring the lock only once.
func (b *BloomUserCache) Populate(emails []string) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for _, email := range emails {
		data := []byte(email)
		for i := uint(0); i < b.k; i++ {
			hash := b.hashFuncs[i](data)
			idx := hash % b.m
			b.set(idx)
		}
	}
}
