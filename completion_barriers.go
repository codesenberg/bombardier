package main

import (
	"sync"
	"sync/atomic"
	"time"
)

type completionBarrier interface {
	completed() float64
	tryGrabWork() bool
	done() <-chan struct{}
	cancel()
}

type countingCompletionBarrier struct {
	numReqs, reqsDone uint64
	doneChan          chan struct{}
	closeOnce         sync.Once
}

func newCountingCompletionBarrier(numReqs uint64) completionBarrier {
	c := new(countingCompletionBarrier)
	c.reqsDone, c.numReqs = 0, numReqs
	c.doneChan = make(chan struct{})
	return completionBarrier(c)
}

func (c *countingCompletionBarrier) tryGrabWork() bool {
	select {
	case <-c.doneChan:
		return false
	default:
		reqsDone := atomic.AddUint64(&c.reqsDone, 1)
		canGrabWork := reqsDone <= c.numReqs
		if !canGrabWork {
			c.closeOnce.Do(func() {
				close(c.doneChan)
			})
		}
		return canGrabWork
	}
}

func (c *countingCompletionBarrier) done() <-chan struct{} {
	return c.doneChan
}

func (c *countingCompletionBarrier) cancel() {
	c.closeOnce.Do(func() {
		close(c.doneChan)
	})
}

func (c *countingCompletionBarrier) completed() float64 {
	select {
	case <-c.doneChan:
		return 1.0
	default:
		reqsDone := atomic.LoadUint64(&c.reqsDone)
		return float64(reqsDone) / float64(c.numReqs)
	}
}

type timedCompletionBarrier struct {
	doneChan  chan struct{}
	closeOnce sync.Once
	start     time.Time
	duration  time.Duration
}

func newTimedCompletionBarrier(duration time.Duration) completionBarrier {
	if duration < 0 {
		panic("timedCompletionBarrier: negative duration")
	}
	c := new(timedCompletionBarrier)
	c.doneChan = make(chan struct{})
	c.start = time.Now()
	c.duration = duration
	go func() {
		time.AfterFunc(duration, func() {
			c.closeOnce.Do(func() {
				close(c.doneChan)
			})
		})
	}()
	return completionBarrier(c)
}

func (c *timedCompletionBarrier) tryGrabWork() bool {
	select {
	case <-c.doneChan:
		return false
	default:
		return true
	}
}

func (c *timedCompletionBarrier) done() <-chan struct{} {
	return c.doneChan
}

func (c *timedCompletionBarrier) cancel() {
	c.closeOnce.Do(func() {
		close(c.doneChan)
	})
}

func (c *timedCompletionBarrier) completed() float64 {
	select {
	case <-c.doneChan:
		return 1.0
	default:
		return float64(time.Since(c.start).Nanoseconds()) /
			float64(c.duration.Nanoseconds())
	}
}
