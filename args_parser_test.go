package main

import (
	"reflect"
	"strconv"
	"testing"
	"time"
)

const (
	programName = "bombardier"
)

func TestInvalidArgsParsing(t *testing.T) {
	expectations := []struct {
		in  []string
		out string
	}{
		{
			[]string{programName},
			"required argument 'url' not provided",
		},
		{
			[]string{programName, "http://google.com", "http://yahoo.com"},
			"unexpected http://yahoo.com",
		},
	}
	for _, e := range expectations {
		p := newKingpinParser()
		if _, err := p.parse(e.in); err == nil ||
			err.Error() != e.out {
			t.Error(err, e.out)
		}
	}
}

func TestUnspecifiedArgParsing(t *testing.T) {
	p := newKingpinParser()
	args := []string{programName, "--someunspecifiedflag"}
	_, err := p.parse(args)
	if err == nil {
		t.Fail()
	}
}

func TestArgsParsing(t *testing.T) {
	ten := uint64(10)
	expectations := []struct {
		in  [][]string
		out config
	}{
		{
			[][]string{{programName, "https://somehost.somedomain"}},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				method:   "GET",
				url:      "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"-c", "10",
					"-n", strconv.FormatUint(defaultNumberOfReqs, decBase),
					"-t", "10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-c10",
					"-n" + strconv.FormatUint(defaultNumberOfReqs, decBase),
					"-t10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--connections", "10",
					"--requests", strconv.FormatUint(defaultNumberOfReqs, decBase),
					"--timeout", "10s",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--connections=10",
					"--requests=" + strconv.FormatUint(defaultNumberOfReqs, decBase),
					"--timeout=10s",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: 10,
				timeout:  10 * time.Second,
				headers:  new(headersList),
				method:   "GET",
				numReqs:  &defaultNumberOfReqs,
				url:      "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--latencies",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-l",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:       defaultNumberOfConns,
				timeout:        defaultTimeout,
				headers:        new(headersList),
				printLatencies: true,
				method:         "GET",
				url:            "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--insecure",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-k",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				insecure: true,
				method:   "GET",
				url:      "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--key", "testclient.key",
					"--cert", "testclient.cert",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--key=testclient.key",
					"--cert=testclient.cert",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				method:   "GET",
				keyPath:  "testclient.key",
				certPath: "testclient.cert",
				url:      "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--method", "POST",
					"--body", "reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--method=POST",
					"--body=reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-m", "POST",
					"-b", "reqbody",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-mPOST",
					"-breqbody",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				method:   "POST",
				body:     "reqbody",
				url:      "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--header", "One: Value one",
					"--header", "Two: Value two",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-H", "One: Value one",
					"-H", "Two: Value two",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--header=One: Value one",
					"--header=Two: Value two",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers: &headersList{
					{"One", "Value one"},
					{"Two", "Value two"},
				},
				method: "GET",
				url:    "https://somehost.somedomain"},
		},
		{
			[][]string{
				{
					programName,
					"--rate", "10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-r", "10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--rate=10",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-r10",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				method:   "GET",
				url:      "https://somehost.somedomain",
				rate:     &ten,
			},
		},
		{
			[][]string{
				{
					programName,
					"--fasthttp",
					"https://somehost.somedomain",
				},
				{
					programName,
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:   defaultNumberOfConns,
				timeout:    defaultTimeout,
				headers:    new(headersList),
				method:     "GET",
				url:        "https://somehost.somedomain",
				clientType: fhttp,
			},
		},
		{
			[][]string{
				{
					programName,
					"--http1",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:   defaultNumberOfConns,
				timeout:    defaultTimeout,
				headers:    new(headersList),
				method:     "GET",
				url:        "https://somehost.somedomain",
				clientType: nhttp1,
			},
		},
		{
			[][]string{
				{
					programName,
					"--http2",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:   defaultNumberOfConns,
				timeout:    defaultTimeout,
				headers:    new(headersList),
				method:     "GET",
				url:        "https://somehost.somedomain",
				clientType: nhttp2,
			},
		},
		{
			[][]string{
				{
					programName,
					"--body-file=testbody.txt",
					"https://somehost.somedomain",
				},
				{
					programName,
					"--body-file", "testbody.txt",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-f", "testbody.txt",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns:     defaultNumberOfConns,
				timeout:      defaultTimeout,
				headers:      new(headersList),
				method:       "GET",
				bodyFilePath: "testbody.txt",
				url:          "https://somehost.somedomain",
			},
		},
		{
			[][]string{
				{
					programName,
					"--stream",
					"https://somehost.somedomain",
				},
				{
					programName,
					"-s",
					"https://somehost.somedomain",
				},
			},
			config{
				numConns: defaultNumberOfConns,
				timeout:  defaultTimeout,
				headers:  new(headersList),
				method:   "GET",
				stream:   true,
				url:      "https://somehost.somedomain",
			},
		},
	}
	for _, e := range expectations {
		for _, args := range e.in {
			p := newKingpinParser()
			cfg, err := p.parse(args)
			if err != nil {
				t.Error(err)
				continue
			}
			if !reflect.DeepEqual(cfg, e.out) {
				t.Logf("Expected: %#v", e.out)
				t.Logf("Got: %#v", cfg)
				t.Fail()
			}
		}
	}
}
