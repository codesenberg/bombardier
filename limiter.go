package main

import (
	"math"
	"sync"
	"time"

	"github.com/juju/ratelimit"
)

type token uint64

const (
	brk token = iota
	cont
)

type limiter interface {
	pace(<-chan struct{}) token
}

type nooplimiter struct{}

func (n *nooplimiter) pace(<-chan struct{}) token {
	return cont
}

type bucketlimiter struct {
	limiter   *ratelimit.Bucket
	timerPool *sync.Pool
}

func newBucketLimiter(rate uint64) limiter {
	fillInterval, quantum := estimate(rate, rateLimitInterval)
	return &bucketlimiter{
		ratelimit.NewBucketWithQuantum(
			fillInterval, int64(quantum), int64(quantum),
		),
		&sync.Pool{
			New: func() interface{} {
				return time.NewTimer(math.MaxInt64)
			},
		},
	}
}

func (b *bucketlimiter) pace(done <-chan struct{}) (res token) {
	wd := b.limiter.Take(1)
	if wd <= 0 {
		return cont
	}

	timer := b.timerPool.Get().(*time.Timer)
	timer.Reset(wd)
	select {
	case <-timer.C:
		res = cont
	case <-done:
		res = brk
	}
	b.timerPool.Put(timer)
	return
}
