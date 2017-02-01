// +build go1.7

package main

import (
	"flag"
	"runtime"
	"testing"
	"time"
)

var serverPort = flag.String("port", "8080", "port to use for benchmarks")

func BenchmarkBombardier(b *testing.B) {
	addr := "localhost:" + *serverPort
	b.Run("single-req-perf", benchmarkFireRequest(addr))
}

func benchmarkFireRequest(addr string) func(bm *testing.B) {
	return func(bm *testing.B) {
		longDuration := 9001 * time.Hour
		b, e := newBombardier(config{
			numConns:       defaultNumberOfConns,
			numReqs:        nil,
			duration:       &longDuration,
			url:            "http://" + addr,
			headers:        new(headersList),
			timeout:        defaultTimeout,
			method:         "GET",
			body:           "",
			printLatencies: false,
		})
		if e != nil {
			bm.Error(e)
		}
		b.disableOutput()
		bm.SetParallelism(int(defaultNumberOfConns) / runtime.NumCPU())
		bm.ResetTimer()
		bm.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				b.performSingleRequest()
			}
		})
	}
}
