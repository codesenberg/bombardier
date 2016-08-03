package main

import (
	"sync/atomic"
	"testing"
	"time"
)

func TestCountingCompletionBarrier(t *testing.T) {
	callsCount := uint64(0)
	expectedCallsCount := uint64(10)
	b := newCountingCompletionBarrier(expectedCallsCount, func() {
		atomic.AddUint64(&callsCount, 1)
	})
	for b.grabWork() {
		b.jobDone()
	}
	if callsCount != expectedCallsCount {
		t.Errorf("Expected to get %v calls, but got %v instead",
			expectedCallsCount, callsCount)
	}
}

func TestCouintingCompletionBarrierWait(t *testing.T) {
	b := newCountingCompletionBarrier(100, func() {})
	go func() {
		for b.grabWork() {
			b.jobDone()
		}
	}()
	wc := make(chan struct{})
	go func() {
		b.wait()
		wc <- struct{}{}
	}()
	select {
	case <-wc:
		return
	case <-time.After(100 * time.Millisecond):
		t.Fail()
	}
}

func TestTimedCompletionBarrier(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	tickCount := uint64(0)
	parties := uint64(10)
	duration := 100 * time.Millisecond
	tickDuration := 10 * time.Millisecond
	b := newTimedCompletionBarrier(parties, tickDuration, duration, func() {
		atomic.AddUint64(&tickCount, 1)
	})
	callsCount := uint64(0)
	for i := uint64(0); i < parties; i++ {
		go func() {
			for b.grabWork() {
				atomic.AddUint64(&callsCount, 1)
				b.jobDone()
			}
		}()
	}
	b.wait()
	if callsCount == 0 {
		t.Errorf("Workers haven't done any work")
	}
	if tickCount == 0 {
		t.Errorf("Tick callback wasn't called even once")
	}
}

func TestTimedCompletionBarrierWait(t *testing.T) {
	parties := uint64(10)
	duration := 100 * time.Millisecond
	timeout := duration * 2
	err := 15 * time.Millisecond
	sleepDuration := 2 * time.Millisecond
	tickDuration := 5 * time.Millisecond
	b := newTimedCompletionBarrier(parties, tickDuration, duration, func() {})
	for i := uint64(0); i < parties; i++ {
		go func() {
			for b.grabWork() {
				b.jobDone()
				time.Sleep(sleepDuration)
			}
		}()
	}
	wc := make(chan time.Duration)
	go func() {
		start := time.Now()
		b.wait()
		wc <- time.Since(start)
	}()
	select {
	case actual := <-wc:
		if !approximatelyEqual(duration, actual, sleepDuration+err) {
			t.Errorf("Expected to run %v, but ran %v instead", duration, actual)
		}
	case <-time.After(timeout):
		t.Error("Barrier hanged")
	}
}

func approximatelyEqual(expected, actual, err time.Duration) bool {
	return expected-err < actual && actual < expected+err
}
