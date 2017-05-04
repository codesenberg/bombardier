package main

import (
	"flag"
	"runtime"
	"testing"
	"time"
)

var (
	serverPort = flag.String("port", "8080", "port to use for benchmarks")
	clientType = flag.String("client-type", "fasthttp",
		"client to use in benchmarks")
)

var (
	longDuration = 9001 * time.Hour
	highRate     = uint64(1000000)
)

func BenchmarkBombardierSingleReqPerf(b *testing.B) {
	addr := "localhost:" + *serverPort
	benchmarkFireRequest(config{
		numConns:       defaultNumberOfConns,
		numReqs:        nil,
		duration:       &longDuration,
		url:            "http://" + addr,
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: false,
		clientType:     clientTypeFromString(*clientType),
	}, b)
}

func BenchmarkBombardierRateLimitPerf(b *testing.B) {
	addr := "localhost:" + *serverPort
	benchmarkFireRequest(config{
		numConns:       defaultNumberOfConns,
		numReqs:        nil,
		duration:       &longDuration,
		url:            "http://" + addr,
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: false,
		rate:           &highRate,
		clientType:     clientTypeFromString(*clientType),
	}, b)
}

func benchmarkFireRequest(c config, bm *testing.B) {
	b, e := newBombardier(c)
	if e != nil {
		bm.Error(e)
	}
	b.disableOutput()
	bm.SetParallelism(int(defaultNumberOfConns) / runtime.NumCPU())
	bm.ResetTimer()
	bm.RunParallel(func(pb *testing.PB) {
		done := b.barrier.done()
		for pb.Next() {
			b.ratelimiter.pace(done)
			b.performSingleRequest()
		}
	})
}
