package main

import (
	"bytes"
	"container/ring"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"
)

func TestBombardierShouldFireSpecifiedNumberOfRequests(t *testing.T) {
	reqsReceived := uint64(0)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&reqsReceived, 1)
	}))
	numReqs := uint64(100)
	noHeaders := new(headersList)
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  &numReqs,
		duration: nil,
		url:      s.URL,
		headers:  noHeaders,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	if reqsReceived != numReqs {
		t.Fail()
	}
}

func TestBombardierShouldRunTestForSpecifiedDuration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode")
	}
	reqsReceived := uint64(0)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&reqsReceived, 1)
	}))
	noHeaders := new(headersList)
	desiredTestDuration := 1 * time.Second
	desiredError := 100 * time.Millisecond
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: &desiredTestDuration,
		url:      s.URL,
		headers:  noHeaders,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	start := time.Now()
	b.bombard()
	testDuration := time.Since(start)
	if !approximatelyEqual(desiredTestDuration, testDuration, desiredError) {
		t.Log(desiredTestDuration, testDuration)
		t.Fail()
	}
	if reqsReceived == 0 {
		t.Fail()
	}
}

func TestBombardierShouldSendHeaders(t *testing.T) {
	requestHeaders := headersList([]header{
		{"Header1", "Value1"},
		{"Header-Two", "value-two"},
	})
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		for _, h := range requestHeaders {
			if r.Header.Get(h.key) != h.value {
				t.Fail()
			}
		}
	}))
	numReqs := uint64(1)
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  &numReqs,
		duration: nil,
		url:      s.URL,
		headers:  &requestHeaders,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
}

func TestBombardierErrorsCodeRecording(t *testing.T) {
	n := 7
	codes := ring.New(n)
	for i := 0; i < n; i++ {
		codes.Value = i*100 + 1
		codes = codes.Next()
	}
	codes = codes.Next()
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		nextCode := codes.Value.(int)
		rw.WriteHeader(nextCode)
		codes = codes.Next()
	}))
	eachCodeCount := uint64(10)
	numReqs := uint64(n) * eachCodeCount
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  &numReqs,
		duration: nil,
		url:      s.URL,
		headers:  new(headersList),
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	expectation := []struct {
		name     string
		reqsGot  uint64
		expected uint64
	}{
		{"errored", b.errored, eachCodeCount * 2},
		{"1xx", b.req1xx, eachCodeCount},
		{"2xx", b.req2xx, eachCodeCount},
		{"3xx", b.req3xx, eachCodeCount},
		{"4xx", b.req4xx, eachCodeCount},
		{"5xx", b.req5xx, eachCodeCount},
	}
	for _, e := range expectation {
		if e.reqsGot != e.expected {
			t.Log(e.name, e.reqsGot, e.expected)
			t.Fail()
		}
	}
}

func TestBombardierTimeoutRecoding(t *testing.T) {
	shortTimeout := 10 * time.Millisecond
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		time.Sleep(shortTimeout * 2)
	}))
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  &numReqs,
		duration: nil,
		url:      s.URL,
		headers:  new(headersList),
		timeout:  shortTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	if b.errored != numReqs {
		t.Fail()
	}
}

func TestBombardierThroughputRecording(t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(response)
	}))
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns: defaultNumberOfConns,
		numReqs:  &numReqs,
		duration: nil,
		url:      s.URL,
		headers:  new(headersList),
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	bytesExpected := uint64(responseSize) * numReqs
	if uint64(b.bytesData) != bytesExpected {
		t.Log(b.bytesData, bytesExpected)
		t.Fail()
	}
	if a, e := b.throughput(), float64(b.bytesTotal)/b.timeTaken.Seconds(); a != e {
		t.Log(a, e)
		t.Fail()
	}
}

func TestBombardierStatsPrinting(t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		rw.Write(response)
	}))
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns:       defaultNumberOfConns,
		numReqs:        &numReqs,
		duration:       nil,
		url:            s.URL,
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: true,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	out := new(bytes.Buffer)
	b.redirectOutputTo(out)
	b.printStats()
	l := out.Len()
	// Here we only test if anything is written
	if l == 0 {
		t.Fail()
	}
}

// REMOVE
// func BenchmarkFireRequest(bm *testing.B) {
// 	responseSize := 1024
// 	response := bytes.Repeat([]byte{'a'}, responseSize)
// 	s := httptest.NewServer(http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
// 		rw.Write(response)
// 	}))
// 	longDuration := 9001 * time.Hour
// 	b, e := newBombardier(config{
// 		numConns:       defaultNumberOfConns,
// 		numReqs:        nil,
// 		duration:       &longDuration,
// 		url:            s.URL,
// 		headers:        new(headersList),
// 		timeout:        defaultTimeout,
// 		method:         "GET",
// 		body:           "",
// 		printLatencies: true,
// 	})
// 	if e != nil {
// 		bm.Error(e)
// 	}
// 	b.disableOutput()
// 	bm.ResetTimer()
// 	bm.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			req, resp := b.prepareRequest()
// 			bytesData, bytesTotal, code, msTaken := b.fireRequest(req, resp)
// 			b.releaseRequest(req, resp)
// 			b.writeStatistics(bytesData, bytesTotal, code, msTaken)
// 		}
// 	})
// }
