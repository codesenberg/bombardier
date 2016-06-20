package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"sync/atomic"
	"time"

	pb "github.com/cheggaaa/pb"
	"github.com/valyala/fasthttp"
)

const (
	maxRps           = 10000000
	requestsInterval = 100 * time.Millisecond
)

type bombardier struct {
	conf           config
	requestHeaders *fasthttp.RequestHeader
	barrier        completionBarrier

	bytesWritten int64
	timeTaken    time.Duration
	latencies    *stats
	requests     *stats

	client *fasthttp.Client
	done   chan bool

	// RPS metrics
	rpl   sync.Mutex
	reqs  int64
	start time.Time

	// HTTP codes
	req1xx  uint64
	req2xx  uint64
	req3xx  uint64
	req4xx  uint64
	req5xx  uint64
	errored uint64

	// Progress bar
	bar *pb.ProgressBar
}

func newBombardier(c config) (*bombardier, error) {
	if err := c.checkArgs(); err != nil {
		return nil, err
	}
	b := new(bombardier)
	b.conf = c
	b.latencies = newStats(c.timeoutMillis())
	b.requests = newStats(maxRps)
	if b.conf.testType == counted {
		b.bar = pb.New64(int64(*b.conf.numReqs))
		b.barrier = newCountingCompletionBarrier(*c.numReqs, func() {
			b.bar.Increment()
		})
	} else if b.conf.testType == timed {
		b.bar = pb.New(int(b.conf.duration.Seconds()))
		b.bar.ShowCounters = false
		b.bar.ShowPercent = false
		b.barrier = newTimedCompletionBarrier(int(c.numConns), *c.duration, func() {
			b.bar.Increment()
		})
	}
	b.client = &fasthttp.Client{
		MaxConnsPerHost: int(c.numConns),
	}
	b.done = make(chan bool)
	b.requestHeaders = c.requestHeaders()
	return b, nil
}

func (b *bombardier) prepareRequest() (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	if b.requestHeaders != nil {
		b.requestHeaders.CopyTo(&req.Header)
	}
	req.Header.SetMethod(b.conf.method)
	req.SetRequestURI(b.conf.url)
	req.SetBodyString(b.conf.body)
	return req, resp
}

func (b *bombardier) fireRequest(req *fasthttp.Request, resp *fasthttp.Response) (bytesWritten int64, code int, msTaken uint64) {
	start := time.Now()
	err := b.client.DoTimeout(req, resp, b.conf.timeout)
	if err != nil {
		code = 0
	} else {
		code = resp.StatusCode()
	}
	bytesWritten, _ = resp.WriteTo(ioutil.Discard)
	msTaken = uint64(time.Since(start).Nanoseconds() / 1000)
	return
}

func (b *bombardier) releaseRequest(req *fasthttp.Request, resp *fasthttp.Response) {
	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}

func (b *bombardier) writeStatistics(bytesWritten int64, code int, msTaken uint64) {
	b.latencies.record(msTaken)
	atomic.AddInt64(&b.bytesWritten, bytesWritten)
	b.rpl.Lock()
	b.reqs++
	b.rpl.Unlock()
	var counter *uint64
	switch code / 100 {
	case 0:
		counter = &b.errored
	case 1:
		counter = &b.req1xx
	case 2:
		counter = &b.req2xx
	case 3:
		counter = &b.req3xx
	case 4:
		counter = &b.req4xx
	case 5:
		counter = &b.req5xx
	default:
		counter = &b.errored
	}
	atomic.AddUint64(counter, 1)
}

func (b *bombardier) Worker() {
	for b.barrier.grabWork() {
		req, resp := b.prepareRequest()
		bytesWritten, code, msTaken := b.fireRequest(req, resp)
		b.releaseRequest(req, resp)
		b.writeStatistics(bytesWritten, code, msTaken)
		b.barrier.jobDone()
	}
}

func (b *bombardier) rateMeter() {
	tick := time.Tick(requestsInterval)
	for {
		select {
		case <-tick:
			b.recordRps()
			continue
		case <-b.done:
			b.recordRps()
			b.done <- true
			return
		}
	}
}

func (b *bombardier) recordRps() {
	b.rpl.Lock()
	duration := time.Since(b.start)
	reqs := b.reqs
	b.reqs = 0
	b.start = time.Now()
	b.rpl.Unlock()
	b.requests.record(uint64(float64(reqs) / duration.Seconds()))
}

func (b *bombardier) bombard() {
	b.printIntro()
	b.bar.Start()
	bombardmentBegin := time.Now()
	b.start = time.Now()
	for i := uint64(0); i < b.conf.numConns; i++ {
		go b.Worker()
	}
	go b.rateMeter()
	b.barrier.wait()
	b.timeTaken = time.Since(bombardmentBegin)
	b.done <- true
	<-b.done
	b.bar.FinishPrint("Done!")
}

func (b *bombardier) throughput() float64 {
	return float64(b.bytesWritten) / b.timeTaken.Seconds()
}

func (b *bombardier) printIntro() {
	if b.conf.testType == counted {
		fmt.Printf("Bombarding %v with %v requests using %v connections\n",
			b.conf.url, *b.conf.numReqs, b.conf.numConns)
	} else if b.conf.testType == timed {
		fmt.Printf("Bombarding %v for %v using %v connections\n",
			b.conf.url, *b.conf.duration, b.conf.numConns)
	}
}

func (b *bombardier) printLatencyStats() {
	percentiles := []float64{50.0, 75.0, 90.0, 99.0}
	fmt.Println("  Latency Distribution")
	for i := 0; i < len(percentiles); i++ {
		p := percentiles[i]
		n := b.latencies.percentile(p)
		fmt.Printf("     %2.0f%% %10s", p, formatUnits(float64(n), timeUnitsUs, 2))
		fmt.Printf("\n")
	}
}

func (b *bombardier) printStats() {
	fmt.Printf("%10v %10v %10v %10v\n", "Statistics", "Avg", "Stdev", "Max")
	fmt.Println(rpsString(b.requests))
	fmt.Println(latenciesString(b.latencies))
	if *latencies {
		b.printLatencyStats()
	}
	fmt.Println("  HTTP codes:")
	fmt.Printf("    1xx - %v, 2xx - %v, 3xx - %v, 4xx - %v, 5xx - %v\n",
		b.req1xx, b.req2xx, b.req3xx, b.req4xx, b.req5xx)
	fmt.Printf("    errored - %v\n", b.errored)
	fmt.Printf("  %-10v %10v/s\n", "Throughput:", formatBinary(b.throughput()))
}

var (
	numReqs   = new(nullableUint64)
	duration  = new(nullableDuration)
	headers   = new(headersList)
	numConns  = flag.Uint64("c", 200, "Maximum number of concurrent connections")
	timeout   = flag.Duration("timeout", 2*time.Second, "Socket/request timeout")
	latencies = flag.Bool("latencies", false, "Print latency statistics")
	method    = flag.String("m", "GET", "Request method")
	body      = flag.String("data", "", "Request body")
)

func main() {
	flag.Var(headers, "H", "HTTP headers to use")
	flag.Var(numReqs, "n", "Number of requests")
	flag.Var(duration, "d", "Duration of test")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Println("No URL supplied")
		flag.Usage()
		os.Exit(1)
	}
	if flag.NArg() > 1 {
		fmt.Println("Too many arguments are supplied")
		os.Exit(1)
	}
	bombardier, err := newBombardier(
		config{
			numConns: *numConns,
			numReqs:  numReqs.val,
			duration: duration.val,
			url:      flag.Arg(0),
			headers:  headers,
			timeout:  *timeout,
			method:   *method,
			body:     *body,
		})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bombardier.bombard()
	bombardier.printStats()
}
