package main

import (
	"fmt"
	"math"
	"sync/atomic"
)

type stats struct {
	count uint64
	limit uint64
	min   uint64
	max   uint64
	data  []uint64
}

func newStats(max uint64) *stats {
	s := new(stats)
	s.limit = uint64(max + 1)
	s.min = math.MaxUint64
	s.data = make([]uint64, s.limit)
	return s
}

func (s *stats) record(val uint64) bool {
	if val >= s.limit {
		return false
	}
	atomic.AddUint64(&s.data[val], 1)
	atomic.AddUint64(&s.count, 1)
	min := atomic.LoadUint64(&s.min)
	for ; val < min; min = atomic.LoadUint64(&s.min) {
		atomic.CompareAndSwapUint64(&s.min, min, val)
	}
	max := atomic.LoadUint64(&s.max)
	for ; val > max; max = atomic.LoadUint64(&s.max) {
		atomic.CompareAndSwapUint64(&s.max, max, val)
	}
	return true
}

func (s *stats) mean() float64 {
	if s.count == 0 {
		return 0.0
	}
	sum := uint64(0)
	for i := s.min; i <= s.max; i++ {
		sum += s.data[i] * i
	}
	return float64(sum / s.count)
}

func (s *stats) stdev(mean float64) float64 {
	sum := 0.0
	if s.count < 2 {
		return 0.0
	}
	for i := s.min; i <= s.max; i++ {
		if s.data[i] != 0 {
			sum += math.Pow(float64(i)-mean, 2) * float64(s.data[i])
		}
	}
	return math.Sqrt(sum / float64(s.count-1))
}

func (s *stats) percentile(p float64) uint64 {
	rank := uint64((p/100.0)*float64(s.count) + 0.5)
	total := uint64(0)
	for i := s.min; i <= s.max; i++ {
		total += s.data[i]
		if total >= rank {
			return i
		}
	}
	return 0
}

func rpsString(s *stats) string {
	mean := s.mean()
	stdev := s.stdev(mean)
	max := s.max
	return fmt.Sprintf("  %-10v %10.2f %10.2f %10d",
		"Reqs/sec", mean, stdev, max)
}

func latenciesString(s *stats) string {
	mean := s.mean()
	stdev := s.stdev(mean)
	max := s.max
	return fmt.Sprintf("  %-10v %10v %10v %10v",
		"Latency",
		formatTimeUs(mean),
		formatTimeUs(stdev),
		formatTimeUs(float64(max)))
}
