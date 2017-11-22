package hash

// Uint64 is a simple hash function for uint64 values
func Uint64(u uint64) uint32 {
	return uint32(u) ^ uint32(u>>32)
}
