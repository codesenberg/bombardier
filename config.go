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
	defaultTestDuration  = 10 * time.Second
	defaultNumberOfConns = uint64(125)
	defaultNumberOfReqs  = uint64(10000)
	httpMethods          = []string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}
	cantHaveBody         = []string{"GET", "HEAD"}

	errInvalidURL              = errors.New("No hostname or invalid scheme")
	errInvalidNumberOfConns    = errors.New("Invalid number of connections(must be > 0)")
	errInvalidNumberOfRequests = errors.New("Invalid number of requests(must be > 0)")
	errInvalidTestDuration     = errors.New("Invalid test duration(must be >= 1s)")
	errNegativeTimeout         = errors.New("Timeout can't be negative")
	errLargeTimeout            = errors.New("Timeout is too big(more that 10s)")
	errBodyNotAllowed          = errors.New("GET and HEAD requests cannot have body")
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
	printLatencies    bool
}

type invalidHTTPMethodError struct {
	method string
}

func (i *invalidHTTPMethodError) Error() string {
	return fmt.Sprintf("Unknown HTTP method: %v", i.method)
}

func (c *config) checkArgs() error {
	url, err := url.ParseRequestURI(c.url)
	if err != nil {
		return err
	}
	if url.Host == "" || (url.Scheme != "http" && url.Scheme != "https") {
		return errInvalidURL
	}
	if c.numConns < uint64(1) {
		return errInvalidNumberOfConns
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
		return errInvalidNumberOfRequests
	}
	if c.testType == timed && *c.duration < time.Second {
		return errInvalidTestDuration
	}
	if c.timeout < 0 {
		return errNegativeTimeout
	}
	if c.timeout > 10*time.Second {
		return errLargeTimeout
	}
	if !allowedHTTPMethod(c.method) {
		return &invalidHTTPMethodError{method: c.method}
	}
	if !canHaveBody(c.method) && len(c.body) > 0 {
		return errBodyNotAllowed
	}
	return nil
}

func (c *config) timeoutMillis() uint64 {
	return uint64(c.timeout.Nanoseconds() / 1000)
}

func (c *config) requestHeaders() *fasthttp.RequestHeader {
	return c.headers.toRequestHeader()
}

func allowedHTTPMethod(method string) bool {
	i := sort.SearchStrings(httpMethods, method)
	return i < len(httpMethods) && httpMethods[i] == method
}

func canHaveBody(method string) bool {
	i := sort.SearchStrings(cantHaveBody, method)
	return !(i < len(cantHaveBody) && cantHaveBody[i] == method)
}
