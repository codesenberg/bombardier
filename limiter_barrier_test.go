package main

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestNoopLimiterCounterBarrierCombination(t *testing.T) {
	expectations := []uint64{
		1, 15, 50, 100, 150, 500, 1000, 1500, 5000,
	}
	done := make(chan struct{})
	for _, count := range expectations {
		b := newCountingCompletionBarrier(count)
		var lim limiter = &nooplimiter{}
		counter := uint64(0)
		numParties := 10
		var wg sync.WaitGroup
		wg.Add(numParties)
		for i := 0; i < numParties; i++ {
			go func() {
				defer wg.Done()
				for b.tryGrabWork() {
					lim.pace(done)
					atomic.AddUint64(&counter, 1)
					b.jobDone()
				}
			}()
		}
		wg.Wait()
		if counter != count {
			t.Error(count, counter)
		}
	}
}

func TestBucketLimiterCounterBarrierCombination(t *testing.T) {
	expectations := []struct {
		count, rate uint64
	}{
		{10, 100},
		{10, 1000},
		{100, 1000},
		{100, 10000},
		{1000, 10000},
		{1000, 100000},
	}
	done := make(chan struct{})
	var expWg sync.WaitGroup
	expWg.Add(len(expectations))
	for i := range expectations {
		exp := expectations[i]
		go func() {
			defer expWg.Done()
			b := newCountingCompletionBarrier(exp.count)
			lim := newBucketLimiter(exp.rate)
			counter := uint64(0)
			numParties := 10
			var wg sync.WaitGroup
			wg.Add(numParties)
			for i := 0; i < numParties; i++ {
				go func() {
					defer wg.Done()
					for b.tryGrabWork() {
						lim.pace(done)
						atomic.AddUint64(&counter, 1)
						b.jobDone()
					}
				}()
			}
			wg.Wait()
			if counter != exp.count {
				t.Error(exp.count, counter)
			}
		}()
	}
	expWg.Wait()
}
