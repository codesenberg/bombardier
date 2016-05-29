package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
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
	numReqs        int
	numConns       int
	url            string
	requestHeaders *fasthttp.RequestHeader
	timeout        time.Duration

	reqsDone uint64

	bytesWritten int64
	timeTaken    time.Duration
	latencies    *stats
	requests     *stats

	jobs   sync.WaitGroup
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

func newBombardier(numConns, numReqs int, url string, headers *headersList, timeout time.Duration) (*bombardier, error) {
	b := new(bombardier)
	b.numReqs = numReqs
	b.numConns = numConns
	b.url = url
	b.timeout = timeout
	if err := b.checkArgs(); err != nil {
		return nil, err
	}
	b.latencies = newStats(b.timeout.Nanoseconds() / 1000)
	b.requests = newStats(maxRps)
	b.jobs.Add(b.numReqs)
	b.client = &fasthttp.Client{
		MaxConnsPerHost: b.numConns,
	}
	b.done = make(chan bool)
	b.requestHeaders = headers.toRequestHeader()
	return b, nil
}

func (b *bombardier) checkArgs() error {
	if b.numReqs < 1 {
		return errors.New("Invalid number of requests(must be > 0)")
	}
	if b.numConns < 1 {
		return errors.New("Invalid number of connections(must be > 0)")
	}
	if b.timeout < 0 {
		return errors.New("Timeout can't be negative")
	}
	if b.timeout > 10*time.Second {
		return errors.New("Timeout is too big(more that 10s)")
	}
	return nil
}

func (b *bombardier) prepareRequest(headers *fasthttp.RequestHeader) (*fasthttp.Request, *fasthttp.Response) {
	req := fasthttp.AcquireRequest()
	resp := fasthttp.AcquireResponse()
	if headers != nil {
		headers.CopyTo(&req.Header)
	}
	req.Header.SetMethod("GET")
	req.SetRequestURI(b.url)
	return req, resp
}

func (b *bombardier) fireRequest(req *fasthttp.Request, resp *fasthttp.Response) (bytesWritten int64, code int, msTaken uint64) {
	start := time.Now()
	err := b.client.DoTimeout(req, resp, b.timeout)
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

func (b *bombardier) grabWork() bool {
	reqID := atomic.AddUint64(&b.reqsDone, 1)
	return reqID <= uint64(b.numReqs)
}

func (b *bombardier) reportDone() {
	b.bar.Increment()
	b.jobs.Done()
}

func (b *bombardier) Worker(headers *fasthttp.RequestHeader) {
	for b.grabWork() {
		req, resp := b.prepareRequest(headers)
		bytesWritten, code, msTaken := b.fireRequest(req, resp)
		b.releaseRequest(req, resp)
		b.writeStatistics(bytesWritten, code, msTaken)
		b.reportDone()
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
	fmt.Printf("Bombarding %v with %v requests using %v connections\n",
		b.url, b.numReqs, b.numConns)
	b.bar = pb.StartNew(b.numReqs)
	bombardmentBegin := time.Now()
	b.start = time.Now()
	for i := 0; i < b.numConns; i++ {
		var headers *fasthttp.RequestHeader
		if b.requestHeaders != nil {
			headers = new(fasthttp.RequestHeader)
			b.requestHeaders.CopyTo(headers)
		}
		go b.Worker(headers)
	}
	go b.rateMeter()
	b.jobs.Wait()
	b.timeTaken = time.Since(bombardmentBegin)
	b.done <- true
	<-b.done
	b.bar.Finish()
}

func (b *bombardier) throughput() float64 {
	return float64(b.bytesWritten) / b.timeTaken.Seconds()
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

var headers = new(headersList)
var numConns = flag.Int("c", 200, "Maximum number of concurrent connections")
var numReqs = flag.Int("n", 10000, "Number of requests")
var timeout = flag.Duration("timeout", 2*time.Second, "Socket/request timeout")
var latencies = flag.Bool("latencies", false, "Print latency statistics")

func main() {
	flag.Var(headers, "H", "HTTP headers to use")
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
	rawurl := flag.Args()[0]
	url, err := url.ParseRequestURI(rawurl)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	if url.Host == "" || (url.Scheme != "http" && url.Scheme != "https") {
		fmt.Println("No hostname or invalid scheme")
		os.Exit(1)
	}
	bombardier, err := newBombardier(
		*numConns, *numReqs,
		url.String(), headers,
		*timeout)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	bombardier.bombard()
	bombardier.printStats()
}
