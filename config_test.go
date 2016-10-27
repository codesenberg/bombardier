package main

import (
	"testing"
	"time"
)

func TestCanHaveBody(t *testing.T) {
	expectations := []struct {
		in  string
		out bool
	}{
		{"GET", false},
		{"POST", true},
		{"PUT", true},
		{"DELETE", true},
		{"HEAD", false},
		{"OPTIONS", true},
	}
	for _, e := range expectations {
		if r := canHaveBody(e.in); r != e.out {
			t.Log(e.in, e.out, r)
			t.Fail()
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
	bigTimeoutDuration := 900 * time.Second
	noHeaders := new(headersList)
	expectations := []struct {
		in  config
		out error
	}{
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
			},
			errInvalidURL,
		},
		{
			config{
				numConns: 0,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
			},
			errInvalidNumberOfConns,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &invalidNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
			},
			errInvalidNumberOfRequests,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  nil,
				duration: &smallTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
			},
			errInvalidTestDuration,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  negativeTimeoutDuration,
				method:   "GET",
				body:     "",
			},
			errNegativeTimeout,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  bigTimeoutDuration,
				method:   "GET",
				body:     "",
			},
			errLargeTimeout,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "BODY",
			},
			errBodyNotAllowed,
		},
		{
			config{
				numConns: defaultNumberOfConns,
				numReqs:  &defaultNumberOfReqs,
				duration: &defaultTestDuration,
				url:      "http://localhost:8080",
				headers:  noHeaders,
				timeout:  defaultTimeout,
				method:   "GET",
				body:     "",
			},
			nil,
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

func TestCheckArgsGarbageUrl(t *testing.T) {
	c := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      "8080",
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if c.checkArgs() == nil {
		t.Fail()
	}
}

func TestCheckArgsInvalidRequestMethod(t *testing.T) {
	c := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      "http://localhost:8080",
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
		url:      "http://localhost:8080",
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	timedConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: &defaultTestDuration,
		url:      "http://localhost:8080",
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	both := config{
		numConns: defaultNumberOfConns,
		numReqs:  &defaultNumberOfReqs,
		duration: &defaultTestDuration,
		url:      "http://localhost:8080",
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	defaultConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      "http://localhost:8080",
		headers:  nil,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if err := countedConfig.checkArgs(); err != nil || countedConfig.testType != counted {
		t.Fail()
	}
	if err := timedConfig.checkArgs(); err != nil || timedConfig.testType != timed {
		t.Fail()
	}
	if err := both.checkArgs(); err != nil || both.testType != counted {
		t.Fail()
	}
	if err := defaultConfig.checkArgs(); err != nil || defaultConfig.testType != timed || defaultConfig.duration != &defaultTestDuration {
		t.Fail()
	}
}

func TestTimeoutMillis(t *testing.T) {
	defaultConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      "http://localhost:8080",
		headers:  nil,
		timeout:  2 * time.Second,
		method:   "GET",
		body:     "",
	}
	if defaultConfig.timeoutMillis() != 2000000 {
		t.Fail()
	}
}

func TestConfigRequestHeaders(t *testing.T) {
	emptyHeaders := headersList([]header{})
	singetonHeader := headersList([]header{
		{"Key", "Value"}})
	defaultConfig := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      "http://localhost:8080",
		headers:  &emptyHeaders,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if defaultConfig.requestHeaders() != nil {
		t.Fail()
	}
	defaultConfigWithHeaders := config{
		numConns: defaultNumberOfConns,
		numReqs:  nil,
		duration: nil,
		url:      "http://localhost:8080",
		headers:  &singetonHeader,
		timeout:  defaultTimeout,
		method:   "GET",
		body:     "",
	}
	if defaultConfigWithHeaders.requestHeaders() == nil {
		t.Fail()
	}
}

func TestInvalidHTTPMethodError(t *testing.T) {
	invalidMethod := "NOSUCHMETHOD"
	want := "Unknown HTTP method: " + invalidMethod
	err := &invalidHTTPMethodError{invalidMethod}
	if got := err.Error(); got != want {
		t.Log(got, want)
		t.Fail()
	}
}
