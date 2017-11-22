package hash

import "math"

// Float64 is a simple hash function for float64 values.
func Float64(f float64) uint32 {
	result := math.Float64bits(f)
	return uint32(result ^ (result >> 32))
}

// Float32 is a simple hash function for float32 values.
func Float32(f float32) uint32 {
	return math.Float32bits(f)
}
