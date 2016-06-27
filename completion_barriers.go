package main

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/tevino/abool"
)

type completionBarrier interface {
	grabWork() bool
	jobDone()
	wait()
}

type countingCompletionBarrier struct {
	numReqs, reqsDone uint64
	doneCallback      func()
	wg                sync.WaitGroup
}

func newCountingCompletionBarrier(numReqs uint64, callback func()) completionBarrier {
	c := new(countingCompletionBarrier)
	c.reqsDone, c.numReqs = 0, numReqs
	c.doneCallback = callback
	c.wg.Add(int(numReqs))
	return completionBarrier(c)
}

func (c *countingCompletionBarrier) grabWork() bool {
	return atomic.AddUint64(&c.reqsDone, 1) <= c.numReqs
}

func (c *countingCompletionBarrier) jobDone() {
	c.doneCallback()
	c.wg.Done()
}

func (c *countingCompletionBarrier) wait() {
	c.wg.Wait()
}

type timedCompletionBarrier struct {
	wg           sync.WaitGroup
	tickCallback func()
	done         *abool.AtomicBool
}

func newTimedCompletionBarrier(parties int, duration time.Duration, callback func()) completionBarrier {
	c := new(timedCompletionBarrier)
	c.tickCallback = callback
	c.done = abool.NewBool(false)
	c.wg.Add(parties)
	go func() {
		deadline := time.Now().Add(duration)
		for time.Now().Before(deadline) {
			c.tickCallback()
			time.Sleep(1 * time.Second)
		}
		c.done.Set()
	}()
	return completionBarrier(c)
}

func (c *timedCompletionBarrier) grabWork() bool {
	done := c.done.IsSet()
	if done {
		c.wg.Done()
	}
	return !done
}

func (c *timedCompletionBarrier) jobDone() {
}

func (c *timedCompletionBarrier) wait() {
	c.wg.Wait()
}
