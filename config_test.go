package main

import (
	"testing"
	"time"
)

var (
	defaultNumberOfReqs = uint64(10000)
)

func TestCanHaveBody(t *testing.T) {
	expectations := []struct {
		in  string
		out bool
	}{
		{"HEAD", false},
		{"GET", true},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"OPTIONS", true},
	}
	for _, e := range expectations {
		if r := canHaveBody(e.in); r != e.out {
			t.Error(e.in, e.out, r)
		}
	}
}

func TestAllowedHttpMethod(t *testing.T) {
	expectations := []struct {
		in  string
		out bool
	}{
		{"GET", true},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"HEAD", true},
		{"OPTIONS", true},
		{"TRUNCATE", false},
	}
	for _, e := range expectations {
		if r := allowedHTTPMethod(e.in); r != e.out {
			t.Logf("Expected f(%v) = %v, but got %v", e.in, e.out, r)
			t.Fail()
		}
	}
}

func TestCheckArgs(t *testing.T) {
	invalidNumberOfReqs := uint64(0)
	smallTestDuration := 99 * time.Millisecond
	negativeTimeoutDuration := -1 * time.Second
	noHeaders := new(headersList)
	zeroRate := uint64(0)
	expectations := []struct {
		in  config
		out error
	}{
		{
			config{
				numConns: 0,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				format:   knownFormat("plain-text"),
			},
			errInvalidNumberOfConns,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &invalidNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				format:   knownFormat("plain-text"),
			},
			errInvalidNumberOfRequests,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  nil,
				duration: &smallTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				format:   knownFormat("plain-text"),
			},
			errInvalidTestDuration,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  negativeTimeoutDuration,
				method:   "GET",
				body:     "",
				format:   knownFormat("plain-text"),
			},
			errNegativeTimeout,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "HEAD",
				body:     "BODY",
				format:   knownFormat("plain-text"),
			},
			errBodyNotAllowed,
		},
		{
			config{
				numConns:     defaultNumberOfConns,
				numReqs:      &defaultNumberOfReqs,
				duration:     &defaultTestDuration,
				url:          ParseURLOrPanic("http://localhost:8080"),
				headers:      noHeaders,
				timeout:      defaultTimeout,
				method:       "HEAD",
				bodyFilePath: "testbody.txt",
				format:       knownFormat("plain-text"),
			},
			errBodyNotAllowed,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "BODY",
				format:   knownFormat("plain-text"),
			},
			nil,
		},
		{
			config{
				numConns:     defaultNumberOfConns,
				numReqs:      &defaultNumberOfReqs,
				duration:     &defaultTestDuration,
				url:          ParseURLOrPanic("http://localhost:8080"),
				headers:      noHeaders,
				timeout:      defaultTimeout,
				method:       "GET",
				bodyFilePath: "testbody.txt",
				format:       knownFormat("plain-text"),
			},
			nil,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				format:   knownFormat("plain-text"),
			},
			nil,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				certPath: "test_cert.pem",
				keyPath:  "",
				format:   knownFormat("plain-text"),
			},
			errNoPathToKey,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
				certPath: "",
				keyPath:  "test_key.pem",
				format:   knownFormat("plain-text"),
			},
			errNoPathToCert,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      ParseURLOrPanic("http://localhost:8080"),
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				rate:     &zeroRate,
				format:   knownFormat("plain-text"),
			},
			errZeroRate,
		},
		{
			config{
				numConns:     defaultNumberOfConns,
				numReqs:      &defaultNumberOfReqs,
				duration:     &defaultTestDuration,
				url:          ParseURLOrPanic("http://localhost:8080"),
				headers:      noHeaders,
				timeout:      defaultTimeout,
				method:       "POST",
				body:         "abracadabra",
				bodyFilePath: "testbody.txt",
				format:       knownFormat("plain-text"),
			},
			errBodyProvidedTwice,
		},
	}
	for _, e := range expectations {
		if r := e.in.checkArgs(); r != e.out {
			t.Logf("Expected (%v).checkArgs to return %v, but got %v", e.in, e.out, r)
			t.Fail()
		}
		if _, r := newBombardier(e.in); r != e.out {
			t.Logf("Expected newBombardier(%v) to return %v, but got %v", e.in, e.out, r)
			t.Fail()
		}
	}
}

func TestCheckArgsUnsupportedURLScheme(t *testing.T) {
	c := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      ParseURLOrPanic("ftp://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if c.checkArgs() != errUnsupportedScheme {
		t.Fail()
	}
}

func TestCheckArgsInvalidRequestMethod(t *testing.T) {
	c := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "ABRACADABRA",
		body:     "",
	}
	e := c.checkArgs()
	if e == nil {
		t.Fail()
	}
	if _, ok := e.(*invalidHTTPMethodError); !ok {
		t.Fail()
	}
}

func TestCheckArgsTestType(t *testing.T) {
	countedConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: nil,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	timedConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: &defaultTestDuration,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	both := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	defaultConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if err := countedConfig.checkArgs(); err != nil ||
		countedConfig.testType() != counted {
		t.Fail()
	}
	if err := timedConfig.checkArgs(); err != nil ||
		timedConfig.testType() != timed {
		t.Fail()
	}
	if err := both.checkArgs(); err != nil ||
		both.testType() != counted {
		t.Fail()
	}
	if err := defaultConfig.checkArgs(); err != nil ||
		defaultConfig.testType() != timed ||
		defaultConfig.duration != &defaultTestDuration {
		t.Fail()
	}
}

func TestTimeoutMillis(t *testing.T) {
	defaultConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      ParseURLOrPanic("http://localhost:8080"),
		headers:  nil,
		timeout:  2 * time.Second,
		method:   "GET",
		body:     "",
	}
	if defaultConfig.timeoutMillis() != 2000000 {
		t.Fail()
	}
}

func TestInvalidHTTPMethodError(t *testing.T) {
	invalidMethod := "NOSUCHMETHOD"
	want := "Unknown HTTP method: " + invalidMethod
	err := &invalidHTTPMethodError{invalidMethod}
	if got := err.Error(); got != want {
		t.Error(got, want)
	}
}

func TestClientTypToStringConversion(t *testing.T) {
	expectations := []struct {
		in  clientTyp
		out string
	}{
		{fhttp, "FastHTTP"},
		{nhttp1, "net/http v1.x"},
		{nhttp2, "net/http v2.0"},
		{42, "unknown client"},
	}
	for _, exp := range expectations {
		act := exp.in.String()
		if act != exp.out {
			t.Errorf("Expected %v, but got %v", exp.out, act)
		}
	}
}

func clientTypeFromString(s string) clientTyp {
	switch s {
	case "fasthttp":
		return fhttp
	case "http1":
		return nhttp1
	case "http2":
		return nhttp2
	default:
		return fhttp
	}
}
