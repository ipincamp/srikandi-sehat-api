package inmemory

import (
	"github.com/ipincamp/srikandi-sehat/internal/core/ports"
	"github.com/ipincamp/srikandi-sehat/pkg/bloomfilter"
)

// NewUserBloomCache creates a new cache adapter that implements
// the ports.UserCache interface using our pkg.BloomUserCache.
// This function acts as a bridge between the concrete implementation
// and the port.
func NewUserBloomCache(m uint64, k uint) *bloomfilter.BloomUserCache {
	return bloomfilter.NewBloomUserCache(m, k)
}

// Compile-time check to ensure our adapter (via pkg) satisfies the port.
var _ ports.UserCache = (*bloomfilter.BloomUserCache)(nil)
