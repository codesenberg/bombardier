package histogram

import (
	"errors"
	"sync"
)

// Errors returned by New.
var (
	ErrNoShardFn      = errors.New("no shard function provided in constructor")
	ErrZeroShardCount = errors.New("can't create histogram with zero shards")
)

// Histogram is a generic histogram implementation intended for use
// in concurrent environment. All methods are goroutine-safe.
// Zero value and nil are not valid Histograms.
type Histogram struct {
	shards  []*histogramShard
	shardFn func(float64) uint32
}

type histogramShard struct {
	sync.RWMutex
	counters map[float64]uint64
}

// New creates a new histogram. It can panic
// if the shardsCount is more than runtime can handle
// (i.e. slice size limit).
func New(
	shardsCount uint32, shardFn func(float64) uint32,
) (*Histogram, error) {
	if shardFn == nil {
		return nil, ErrNoShardFn
	}
	if shardsCount == 0 {
		return nil, ErrZeroShardCount
	}
	shards := make([]*histogramShard, shardsCount)
	for i := range shards {
		shards[i] = &histogramShard{
			counters: make(map[float64]uint64),
		}
	}
	return &Histogram{
		shards:  shards,
		shardFn: shardFn,
	}, nil
}

func (h *Histogram) getShardFor(key float64) *histogramShard {
	shardNum := int(h.shardFn(key)) % len(h.shards)
	return h.shards[shardNum]
}

// Increment increments counter associated with the specified key by 1.
func (h *Histogram) Increment(key float64) {
	h.Add(key, 1)
}

// Add increments counter associated with the specified key by
// the amount specified.
func (h *Histogram) Add(key float64, amount uint64) {
	shard := h.getShardFor(key)
	shard.Lock()
	shard.counters[key] += amount
	shard.Unlock()
}

// Get returns current value for the given key.
func (h *Histogram) Get(key float64) uint64 {
	shard := h.getShardFor(key)
	shard.RLock()
	result := shard.counters[key]
	shard.RUnlock()
	return result
}

// Count calculates the number of elements in histogram in a more
// efficient way that can be implemented with VisitAll
func (h *Histogram) Count() uint64 {
	for _, sh := range h.shards {
		sh.RLock()
	}
	count := uint64(0)
	for _, sh := range h.shards {
		count += uint64(len(sh.counters))
	}
	for _, sh := range h.shards {
		sh.RUnlock()
	}
	return count
}

// VisitAll function applies function to each key-value pair in histogram.
// If function returns false iteration stops. Locks the entire histogram.
func (h *Histogram) VisitAll(fn func(float64, uint64) bool) {
	for _, sh := range h.shards {
		sh.RLock()
	}
	defer func() {
		for _, sh := range h.shards {
			sh.RUnlock()
		}
	}()
	for _, sh := range h.shards {
		for k, v := range sh.counters {
			if !fn(k, v) {
				return
			}
		}
	}
}
