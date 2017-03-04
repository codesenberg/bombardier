package main

import "math"

func round(f float64) uint64 {
	if math.Abs(f) < 0.5 {
		return 0.0
	}
	return uint64(f + math.Copysign(0.5, f))
}
