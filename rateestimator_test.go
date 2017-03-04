package main

import (
	"testing"
	"time"
)

func TestRateEstimatorPanicWithZeroRate(t *testing.T) {
	defer func() {
		pv, ok := recover().(string)
		if !ok {
			t.Error("expected string value")
			return
		}
		if pv != panicZeroRate {
			t.Error(panicZeroRate, pv)
		}
	}()
	_, _ = estimate(0, 10*time.Second)
	t.Error("should fail with rate == 0")
}

func TestRateEstimatorPanicWithNegativeAdjustTo(t *testing.T) {
	defer func() {
		pv, ok := recover().(string)
		if !ok {
			t.Error("expected string value")
			return
		}
		if pv != panicNegativeAdjustTo {
			t.Error(panicNegativeAdjustTo, pv)
		}
	}()
	_, _ = estimate(10, -10*time.Second)
	t.Error("should fail with adjustTo <= 0")
}

func TestRateEstimatorAccuracy(t *testing.T) {
	defer func() {
		rv := recover()
		if rv != nil {
			t.Error(rv)
		}
	}()
	expectations := []struct {
		rate                 uint64
		adjustTo             time.Duration
		expectedQuantum      uint64
		expectedFillInterval time.Duration
	}{
		{1, 100 * time.Millisecond, 1, 1 * time.Second},
		{1, 1000 * time.Millisecond, 1, 1 * time.Second},
		{1, 2000 * time.Millisecond, 2, 2 * time.Second},
		{1, 3000 * time.Millisecond, 3, 3 * time.Second},
		{4, 3000 * time.Millisecond, 12, 3 * time.Second},
		{10000, 100 * time.Millisecond, 1000, 100 * time.Millisecond},
		{100000, 100 * time.Millisecond, 10000, 100 * time.Millisecond},
		{1000000, 100 * time.Millisecond, 100000, 100 * time.Millisecond},
	}
	for _, exp := range expectations {
		actualFillInterval, actualQuantum := estimate(exp.rate, exp.adjustTo)
		if actualFillInterval != exp.expectedFillInterval ||
			actualQuantum != exp.expectedQuantum {
			t.Log("Expected: ", exp.expectedQuantum, exp.expectedFillInterval)
			t.Log("Actual: ", actualQuantum, actualFillInterval)
			t.Fail()
		}
	}
}
