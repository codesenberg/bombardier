package main

import (
	"sync"
	"sync/atomic"
)

type errorMap struct {
	mu sync.RWMutex
	m  map[string]*uint64
}

func newErrorMap() *errorMap {
	em := new(errorMap)
	em.m = make(map[string]*uint64)
	return em
}

func (e *errorMap) add(err error) {
	s := err.Error()
	e.mu.RLock()
	c, ok := e.m[s]
	e.mu.RUnlock()
	if !ok {
		e.mu.Lock()
		c, ok = e.m[s]
		if !ok {
			c = new(uint64)
			e.m[s] = c
		}
		e.mu.Unlock()
	}
	atomic.AddUint64(c, 1)
}

func (e *errorMap) get(err error) uint64 {
	s := err.Error()
	e.mu.RLock()
	defer e.mu.RUnlock()
	c := e.m[s]
	if c == nil {
		return uint64(0)
	} else {
		return *c
	}
}

func (e *errorMap) sum() uint64 {
	e.mu.RLock()
	defer e.mu.RUnlock()
	sum := uint64(0)
	for _, v := range e.m {
		sum += *v
	}
	return sum
}
