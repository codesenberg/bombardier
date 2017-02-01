package main

import (
	"math"
	"testing"
)

func TestShouldInitializeStats(t *testing.T) {
	max := uint64(1000)
	limit := max + 1
	s := newStats(max)
	if s.limit != limit {
		t.Error("limit must be equal to max+1")
	}
	if s.min != math.MaxUint64 {
		t.Error("min must be equal to MaxUint64")
	}
	if s.data == nil {
		t.Error("data shouldn't be nil")
	}
	if len(s.data) != int(limit) {
		t.Error("data's length must be equal to limit")
	}
}

func TestStatsShouldntRecordValuesHigherThanMax(t *testing.T) {
	max := uint64(100)
	s := newStats(max)
	if s.record(max + 1) {
		t.Error("Shouldn't record values higher than max")
	}
}

func TestStatsShouldRecordValues(t *testing.T) {
	max := uint64(100)
	s := newStats(max)
	v := max / 2
	if !s.record(v) {
		t.Fail()
	}
	if s.count != 1 || s.data[v] != 1 || s.min != v || s.max != v {
		t.Fail()
	}
}

func TestStatsMean(t *testing.T) {
	max := uint64(100)
	s := newStats(max)
	if s.mean() != 0 {
		t.Fail()
	}
	for i := uint64(0); i <= max; i++ {
		s.record(i)
	}
	if s.mean() != float64(max/2) {
		t.Fail()
	}
}

func TestStatsStdev(t *testing.T) {
	max := uint64(100)
	s := newStats(max)
	mean := s.mean()
	stdev := s.stdev(mean)
	if stdev != 0 {
		t.Fail()
	}
	for i := uint64(0); i <= max; i++ {
		s.record(max / 2)
	}
	mean = s.mean()
	stdev = s.stdev(mean)
	if stdev != 0 {
		t.Fail()
	}
	for i := uint64(0); i <= max; i++ {
		s.record(0)
	}
	mean = s.mean()
	stdev = s.stdev(mean)
	if !equalFloats(float64(max/4), stdev, 0.5) {
		t.Fail()
	}
}

func TestEmptyStatsPercentile(t *testing.T) {
	max := uint64(10000)
	s := newStats(max)
	expectations := []struct {
		in  float64
		out uint64
	}{
		{25, 0},
		{50, 0},
		{75, 0},
		{99, 0},
		{99.99, 0},
	}
	for _, e := range expectations {
		if s.percentile(e.in) != e.out {
			t.Fail()
		}
	}
}

func TestStatsPercentile(t *testing.T) {
	max := uint64(10000)
	s := newStats(max)
	for i := uint64(1); i <= max; i++ {
		s.record(i)
	}
	expectations := []struct {
		in  float64
		out uint64
	}{
		{25, 25 * (max / 100)},
		{50, 50 * (max / 100)},
		{75, 75 * (max / 100)},
		{99, 99 * (max / 100)},
		{99.99, uint64(float64(99.99) * float64(max/100))},
	}
	for _, e := range expectations {
		if s.percentile(e.in) != e.out {
			t.Fail()
		}
	}
}

var singleValueStats = func(max, val uint64) *stats {
	s := newStats(max)
	s.record(val)
	return s
}

func TestStatsToStringConversions(t *testing.T) {
	s := singleValueStats(100, 50)
	expectations := []struct {
		actual, expected string
	}{
		{rpsString(s), "  Reqs/sec        50.00       0.00         50"},
		{latenciesString(s), "  Latency       50.00us     0.00us    50.00us"},
	}
	for _, exp := range expectations {
		if exp.actual != exp.expected {
			t.Log("Wanted: \"" + exp.expected + "\"")
			t.Log("Got: \"" + exp.actual + "\"")
			t.Fail()
		}
	}
}

func equalFloats(expected, actual, err float64) bool {
	return expected-err < actual && actual < expected+err
}
