package main

import (
	"fmt"
	"net/url"
	"sort"
	"time"
)

type config struct {
	numConns                  uint64
	numReqs                   *uint64
	disableKeepAlives         bool
	duration                  *time.Duration
	url                       *url.URL
	method, certPath, keyPath string
	body, bodyFilePath        string
	stream                    bool
	headers                   *headersList
	timeout                   time.Duration
	// TODO(codesenberg): printLatencies should probably be
	// re(named&maked) into printPercentiles or even let
	// users provide their own percentiles and not just
	// calculate for [0.5, 0.75, 0.9, 0.99]
	printLatencies, insecure bool
	rate                     *uint64
	clientType               clientTyp

	printIntro, printProgress, printResult bool

	format format
}

type testTyp int

const (
	none testTyp = iota
	timed
	counted
)

type invalidHTTPMethodError struct {
	method string
}

func (i *invalidHTTPMethodError) Error() string {
	return fmt.Sprintf("Unknown HTTP method: %v", i.method)
}

func (c *config) checkArgs() error {
	c.checkOrSetDefaultTestType()

	checks := []func() error{
		c.checkURL,
		c.checkRate,
		c.checkRunParameters,
		c.checkTimeoutDuration,
		c.checkHTTPParameters,
		c.checkCertPaths,
	}

	for _, check := range checks {
		if err := check(); err != nil {
			return err
		}
	}

	return nil
}

func (c *config) checkOrSetDefaultTestType() {
	if c.testType() == none {
		c.duration = &defaultTestDuration
	}
}

func (c *config) testType() testTyp {
	typ := none
	if c.numReqs != nil {
		typ = counted
	} else if c.duration != nil {
		typ = timed
	}
	return typ
}

func (c *config) checkURL() error {
	if c.url.Scheme != "http" && c.url.Scheme != "https" {
		return errUnsupportedScheme
	}
	return nil
}

func (c *config) checkRate() error {
	if c.rate != nil && *c.rate < 1 {
		return errZeroRate
	}
	return nil
}

func (c *config) checkRunParameters() error {
	if c.numConns < uint64(1) {
		return errInvalidNumberOfConns
	}
	if c.testType() == counted && *c.numReqs < uint64(1) {
		return errInvalidNumberOfRequests
	}
	if c.testType() == timed && *c.duration < time.Second {
		return errInvalidTestDuration
	}
	return nil
}

func (c *config) checkTimeoutDuration() error {
	if c.timeout < 0 {
		return errNegativeTimeout
	}
	return nil
}

func (c *config) checkHTTPParameters() error {
	if !allowedHTTPMethod(c.method) {
		return &invalidHTTPMethodError{method: c.method}
	}
	if !canHaveBody(c.method) && (c.body != "" || c.bodyFilePath != "") {
		return errBodyNotAllowed
	}
	if c.body != "" && c.bodyFilePath != "" {
		return errBodyProvidedTwice
	}
	return nil
}

func (c *config) checkCertPaths() error {
	if c.certPath != "" && c.keyPath == "" {
		return errNoPathToKey
	} else if c.certPath == "" && c.keyPath != "" {
		return errNoPathToCert
	}
	return nil
}

func (c *config) timeoutMillis() uint64 {
	return uint64(c.timeout.Nanoseconds() / 1000)
}

func allowedHTTPMethod(method string) bool {
	i := sort.SearchStrings(httpMethods, method)
	return i < len(httpMethods) && httpMethods[i] == method
}

func canHaveBody(method string) bool {
	i := sort.SearchStrings(cantHaveBody, method)
	return !(i < len(cantHaveBody) && cantHaveBody[i] == method)
}

type clientTyp int

const (
	fhttp clientTyp = iota
	nhttp1
	nhttp2
)

func (ct clientTyp) String() string {
	switch ct {
	case fhttp:
		return "FastHTTP"
	case nhttp1:
		return "net/http v1.x"
	case nhttp2:
		return "net/http v2.0"
	}
	return "unknown client"
}
