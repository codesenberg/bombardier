package main

import (
	"regexp"
	"sort"
	"strconv"
	"sync"
	"sync/atomic"
)

var (
	// capture 1 should capture the canonicalised error message
	errMsgPatterns = []string{
		// read tcp 10.10.0.62:60682->63.35.24.107:443: read: connection reset by peer
		`read tcp [\d.]+:\d+->[\d.]+:\d+: read: (.*)`,
	}

	errMsgRegexes []*regexp.Regexp
)

func init() {
	for _, pattern := range errMsgPatterns {
		errMsgRegexes = append(errMsgRegexes, regexp.MustCompile(pattern))
	}
}

func canonicalizeError(err error) string {
	msg := err.Error()
	for _, errMsgRegex := range errMsgRegexes {
		matches := errMsgRegex.FindStringSubmatch(msg)
		if matches != nil {
			return matches[1]
		}
	}

	return msg
}

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
	s := canonicalizeError(err)
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
	}
	return *c
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

type errorWithCount struct {
	error string
	count uint64
}

func (ewc *errorWithCount) String() string {
	return "<" + ewc.error + ":" +
		strconv.FormatUint(ewc.count, decBase) + ">"
}

type errorsByFrequency []*errorWithCount

func (ebf errorsByFrequency) Len() int {
	return len(ebf)
}

func (ebf errorsByFrequency) Less(i, j int) bool {
	return ebf[i].count > ebf[j].count
}

func (ebf errorsByFrequency) Swap(i, j int) {
	ebf[i], ebf[j] = ebf[j], ebf[i]
}

func (e *errorMap) byFrequency() errorsByFrequency {
	e.mu.RLock()
	byFreq := make(errorsByFrequency, 0, len(e.m))
	for err, count := range e.m {
		byFreq = append(byFreq, &errorWithCount{err, *count})
	}
	e.mu.RUnlock()
	sort.Sort(byFreq)
	return byFreq
}
