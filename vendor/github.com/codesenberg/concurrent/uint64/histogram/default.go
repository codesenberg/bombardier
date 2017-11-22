package histogram

import "github.com/codesenberg/concurrent/uint64/hash"

// WithDefaultHash creates histogram with specified shardCount and
// reasonable sharding function.
func WithDefaultHash(shardsCount uint32) (*Histogram, error) {
	return New(shardsCount, hash.Uint64)
}

// Default creates histogram with reasonable defaults.
func Default() *Histogram {
	// We can safely ignore the error in this case
	h, _ := WithDefaultHash(32)
	return h
}
