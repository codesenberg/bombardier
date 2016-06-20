package main

import (
	"errors"
	"fmt"
	"net/url"
	"sort"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	none = iota
	timed
	counted
)

var (
	defaultTestDuration = 10 * time.Second
	httpMethods         = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
	cantHaveBody        = []string{"GET", "HEAD"}
)

func init() {
	sort.Strings(httpMethods)
	sort.Strings(cantHaveBody)
}

type config struct {
	numConns          uint64
	numReqs           *uint64
	duration          *time.Duration
	url, method, body string
	headers           *headersList
	timeout           time.Duration
	testType          int
}

func (c *config) checkArgs() error {
	url, err := url.ParseRequestURI(c.url)
	if err != nil {
		return err
	}
	if url.Host == "" || (url.Scheme != "http" && url.Scheme != "https") {
		return errors.New("No hostname or invalid scheme")
	}
	if c.numConns < uint64(1) {
		return errors.New("Invalid number of connections(must be > 0)")
	}
	testType := none
	if c.numReqs != nil {
		testType = counted
	} else if c.duration != nil {
		testType = timed
	}
	c.testType = testType
	if c.testType == none {
		c.testType = timed
		c.duration = &defaultTestDuration
	}
	if c.testType == counted && *c.numReqs < uint64(1) {
		return errors.New("Invalid number of requests(must be > 0)")
	}
	if c.testType == timed && *c.duration < time.Second {
		return errors.New("Invalid test duration(must be >= 1s)")
	}
	if c.timeout < 0 {
		return errors.New("Timeout can't be negative")
	}
	if c.timeout > 10*time.Second {
		return errors.New("Timeout is too big(more that 10s)")
	}
	if allowedHttpMethod(c.method) {
		return errors.New(fmt.Sprintf("Unknown HTTP method: %v", c.method))
	}
	if !canHaveBody(c.method) && len(c.body) > 0 {
		return errors.New("GET and HEAD requests cannot have body")
	}
	return nil
}

func (c *config) timeoutMillis() int64 {
	return c.timeout.Nanoseconds() / 1000
}

func (c *config) requestHeaders() *fasthttp.RequestHeader {
	return c.headers.toRequestHeader()
}

func allowedHttpMethod(method string) bool {
	i := sort.SearchStrings(httpMethods, method)
	return !(i < len(httpMethods) && httpMethods[i] == method)
}

func canHaveBody(method string) bool {
	i := sort.SearchStrings(cantHaveBody, method)
	return !(i < len(cantHaveBody) && cantHaveBody[i] == method)
}
