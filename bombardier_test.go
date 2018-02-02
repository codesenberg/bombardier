package main

import (
	"bytes"
	"container/ring"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/valyala/fasthttp"
)

func TestBombardierShouldFireSpecifiedNumberOfRequests(t *testing.T) {
	testAllClients(t, testBombardierShouldFireSpecifiedNumberOfRequests)
}

func testBombardierShouldFireSpecifiedNumberOfRequests(
	clientType clientTyp, t *testing.T,
) {
	reqsReceived := uint64(0)
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			atomic.AddUint64(&reqsReceived, 1)
		}),
	)
	defer s.Close()
	numReqs := uint64(100)
	noHeaders := new(headersList)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &numReqs,
		url:        s.URL,
		headers:    noHeaders,
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	if reqsReceived != numReqs {
		t.Fail()
	}
}

func TestBombardierShouldFinish(t *testing.T) {
	testAllClients(t, testBombardierShouldFinish)
}

func testBombardierShouldFinish(clientType clientTyp, t *testing.T) {
	reqsReceived := uint64(0)
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			atomic.AddUint64(&reqsReceived, 1)
		}),
	)
	defer s.Close()
	noHeaders := new(headersList)
	desiredTestDuration := 1 * time.Second
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		duration:   &desiredTestDuration,
		url:        s.URL,
		headers:    noHeaders,
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	waitCh := make(chan struct{})
	go func() {
		b.bombard()
		waitCh <- struct{}{}
	}()
	select {
	case <-waitCh:
	// Do nothing here
	case <-time.After(desiredTestDuration + 5*time.Second):
		t.Fail()
	}
	if reqsReceived == 0 {
		t.Fail()
	}
}

func TestBombardierShouldSendHeaders(t *testing.T) {
	testAllClients(t, testBombardierShouldSendHeaders)
}

func testBombardierShouldSendHeaders(clientType clientTyp, t *testing.T) {
	requestHeaders := headersList([]header{
		{"Header1", "Value1"},
		{"Header-Two", "value-two"},
	})

	// It's a bit hacky, but FastHTTP can't send Host header correctly
	// as of now
	if clientType != fhttp {
		requestHeaders = append(requestHeaders, header{"Host", "web"})
	}

	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			for _, h := range requestHeaders {
				av := r.Header.Get(h.key)
				if h.key == "Host" {
					av = r.Host
				}
				if av != h.value {
					t.Logf("%q <-> %q", av, h.value)
					t.Fail()
				}
			}
		}),
	)
	defer s.Close()
	numReqs := uint64(1)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &numReqs,
		url:        s.URL,
		headers:    &requestHeaders,
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
}

func TestBombardierHTTPCodeRecording(t *testing.T) {
	testAllClients(t, testBombardierHTTPCodeRecording)
}

func testBombardierHTTPCodeRecording(clientType clientTyp, t *testing.T) {
	cs := []int{1, 102, 200, 302, 404, 505, 606}
	codes := ring.New(len(cs))
	for _, v := range cs {
		codes.Value = v
		codes = codes.Next()
	}
	codes = codes.Next()
	var m sync.Mutex
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			m.Lock()
			nextCode := codes.Value.(int)
			codes = codes.Next()
			m.Unlock()
			if nextCode/100 == 3 {
				rw.Header().Set("Location", "http://localhost:666")
			}
			rw.WriteHeader(nextCode)
		}),
	)
	defer s.Close()
	eachCodeCount := uint64(10)
	numReqs := uint64(len(cs)) * eachCodeCount
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &numReqs,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	expectation := []struct {
		name     string
		reqsGot  uint64
		expected uint64
	}{
		{"errored", b.others, eachCodeCount * 2},
		{"1xx", b.req1xx, eachCodeCount},
		{"2xx", b.req2xx, eachCodeCount},
		{"3xx", b.req3xx, eachCodeCount},
		{"4xx", b.req4xx, eachCodeCount},
		{"5xx", b.req5xx, eachCodeCount},
	}
	for _, e := range expectation {
		if e.reqsGot != e.expected {
			t.Error(e.name, e.reqsGot, e.expected)
		}
	}
	t.Logf("%+v", b.errors.byFrequency())
}

func TestBombardierTimeoutRecoding(t *testing.T) {
	testAllClients(t, testBombardierTimeoutRecoding)
}

func testBombardierTimeoutRecoding(clientType clientTyp, t *testing.T) {
	shortTimeout := 10 * time.Millisecond
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			time.Sleep(shortTimeout * 2)
		}),
	)
	defer s.Close()
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &numReqs,
		duration:   nil,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    shortTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	if b.errors.sum() != numReqs {
		t.Fail()
	}
}

func TestBombardierThroughputRecording(t *testing.T) {
	testAllClients(t, testBombardierThroughputRecording)
}

func testBombardierThroughputRecording(clientType clientTyp, t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			_, err := rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &numReqs,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
	}
	b.disableOutput()
	b.bombard()
	if b.bytesRead == 0 || b.bytesWritten == 0 {
		t.Error(b.bytesRead, b.bytesWritten)
	}
}

func TestBombardierStatsPrinting(t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			_, err := rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	numReqs := uint64(10)
	b, e := newBombardier(config{
		numConns:       defaultNumberOfConns,
		numReqs:        &numReqs,
		url:            s.URL,
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: true,
		printIntro:     true,
		printProgress:  true,
		printResult:    true,
	})
	if e != nil {
		t.Error(e)
		return
	}
	dummy := errors.New("dummy error")
	b.errors.add(dummy)

	out := new(bytes.Buffer)
	b.redirectOutputTo(out)
	b.bombard()

	b.printStats()
	l := out.Len()
	// Here we only test if anything is written
	if l == 0 {
		t.Fail()
	}
}

func TestBombardierErrorIfFailToReadClientCert(t *testing.T) {
	numReqs := uint64(10)
	_, e := newBombardier(config{
		numConns:       defaultNumberOfConns,
		numReqs:        &numReqs,
		url:            "http://localhost",
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: true,
		certPath:       "certPath",
		keyPath:        "keyPath",
	})
	if e == nil {
		t.Fail()
	}
}

func TestBombardierClientCerts(t *testing.T) {
	testAllClients(t, testBombardierClientCerts)
}

func testBombardierClientCerts(clientType clientTyp, t *testing.T) {
	clientCert, err := tls.LoadX509KeyPair("testclient.cert", "testclient.key")
	if err != nil {
		t.Error(err)
		return
	}

	x509Cert, err := x509.ParseCertificate(clientCert.Certificate[0])
	if err != nil {
		t.Error(err)
		return
	}

	server := fasthttp.Server{
		DisableKeepalive: true,
		Handler: func(ctx *fasthttp.RequestCtx) {
			certs := ctx.TLSConnectionState().PeerCertificates
			if numCerts := len(certs); numCerts != 1 {
				t.Errorf("expected 1 cert, but got %v", numCerts)
				ctx.Error("invalid number of certs", http.StatusBadRequest)
				return
			}

			cert := certs[0]
			if !cert.Equal(x509Cert) {
				t.Error("certificates don't match")
				ctx.Error("wrong cert", http.StatusBadRequest)
				return
			}

			ctx.Success("text/plain; charset=utf-8", []byte("OK"))
		},
	}

	ln, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		t.Error(err)
		return
	}

	go func() {
		serr := server.ServeTLS(ln, "testserver.cert", "testserver.key")
		if serr != nil {
			t.Error(err)
		}
	}()

	numReqs := uint64(1)
	b, e := newBombardier(config{
		numConns:       defaultNumberOfConns,
		numReqs:        &numReqs,
		url:            "https://localhost:8080/",
		headers:        new(headersList),
		timeout:        defaultTimeout,
		method:         "GET",
		body:           "",
		printLatencies: true,
		certPath:       "testclient.cert",
		keyPath:        "testclient.key",
		insecure:       true,
		clientType:     clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()

	b.bombard()
	if b.req2xx != 1 {
		t.Error("no requests succeeded")
	}

	err = ln.Close()
	if err != nil {
		t.Error(err)
	}
	// TODO(codesenberg): remove. Another hacky attempt to fix Travis CI's
	// slowness
	time.Sleep(100 * time.Millisecond)
}

func TestBombardierRateLimiting(t *testing.T) {
	testAllClients(t, testBombardierRateLimiting)
}

func testBombardierRateLimiting(clientType clientTyp, t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			_, err := rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	rate := uint64(5000)
	testDuration := 1 * time.Second
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		duration:   &testDuration,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    defaultTimeout,
		method:     "GET",
		body:       "",
		rate:       &rate,
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()
	b.bombard()
	if float64(b.req2xx) < float64(rate)*0.75 ||
		float64(b.req2xx) > float64(rate)*1.25 {
		t.Error(rate, b.req2xx)
	}
}

func testAllClients(parent *testing.T, testFun func(clientTyp, *testing.T)) {
	clients := []clientTyp{fhttp, nhttp1, nhttp2}
	for _, ct := range clients {
		parent.Run(ct.String(), func(t *testing.T) {
			testFun(ct, t)
		})
	}
}

func TestBombardierSendsBody(t *testing.T) {
	testAllClients(t, testBombardierSendsBody)
}

func testBombardierSendsBody(clientType clientTyp, t *testing.T) {
	response := []byte("OK")
	requestBody := "abracadabra"
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}
			if string(body) != requestBody {
				t.Errorf("Expected %v, but got %v", requestBody, string(body))
			}
			_, err = rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	one := uint64(1)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &one,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    defaultTimeout,
		method:     "POST",
		body:       requestBody,
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()
	b.bombard()
}

func TestBombardierSendsBodyFromFile(t *testing.T) {
	testAllClients(t, testBombardierSendsBodyFromFile)
}

func testBombardierSendsBodyFromFile(clientType clientTyp, t *testing.T) {
	response := []byte("OK")
	bodyPath := "testbody.txt"
	requestBody, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		t.Error(err)
		return
	}
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}
			if string(body) != string(requestBody) {
				t.Errorf("Expected %v, but got %v", string(requestBody), string(body))
			}
			_, err = rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	one := uint64(1)
	b, e := newBombardier(config{
		numConns:     defaultNumberOfConns,
		numReqs:      &one,
		url:          s.URL,
		headers:      new(headersList),
		timeout:      defaultTimeout,
		method:       "POST",
		bodyFilePath: bodyPath,
		clientType:   clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()
	b.bombard()
}

func TestBombardierFileDoesntExist(t *testing.T) {
	bodyPath := "/does/not/exist.forreal"
	_, e := newBombardier(config{
		numConns:     defaultNumberOfConns,
		url:          "http://example.com",
		headers:      new(headersList),
		timeout:      defaultTimeout,
		method:       "POST",
		bodyFilePath: bodyPath,
	})
	_, ok := e.(*os.PathError)
	if !ok {
		t.Errorf("Expected to get PathError, but got %v", e)
	}
}

func TestBombardierStreamsBody(t *testing.T) {
	testAllClients(t, testBombardierStreamsBody)
}

func testBombardierStreamsBody(clientType clientTyp, t *testing.T) {
	response := []byte("OK")
	requestBody := "abracadabra"
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if te := r.TransferEncoding; !reflect.DeepEqual(te, []string{"chunked"}) {
				t.Errorf("Expected chunked transfer encoding, but got %v", te)
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}
			if string(body) != requestBody {
				t.Errorf("Expected %v, but got %v", requestBody, string(body))
			}
			_, err = rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	one := uint64(1)
	b, e := newBombardier(config{
		numConns:   defaultNumberOfConns,
		numReqs:    &one,
		url:        s.URL,
		headers:    new(headersList),
		timeout:    defaultTimeout,
		method:     "POST",
		body:       requestBody,
		stream:     true,
		clientType: clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()
	b.bombard()
}

func TestBombardierStreamsBodyFromFile(t *testing.T) {
	testAllClients(t, testBombardierStreamsBodyFromFile)
}

func testBombardierStreamsBodyFromFile(clientType clientTyp, t *testing.T) {
	response := []byte("OK")
	bodyPath := "testbody.txt"
	requestBody, err := ioutil.ReadFile(bodyPath)
	if err != nil {
		t.Error(err)
		return
	}
	s := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if te := r.TransferEncoding; !reflect.DeepEqual(te, []string{"chunked"}) {
				t.Errorf("Expected chunked transfer encoding, but got %v", te)
			}
			body, err := ioutil.ReadAll(r.Body)
			if err != nil {
				t.Error(err)
				return
			}
			if string(body) != string(requestBody) {
				t.Errorf("Expected %v, but got %v", string(requestBody), string(body))
			}
			_, err = rw.Write(response)
			if err != nil {
				t.Error(err)
			}
		}),
	)
	defer s.Close()
	one := uint64(1)
	b, e := newBombardier(config{
		numConns:     defaultNumberOfConns,
		numReqs:      &one,
		url:          s.URL,
		headers:      new(headersList),
		timeout:      defaultTimeout,
		method:       "POST",
		bodyFilePath: bodyPath,
		stream:       true,
		clientType:   clientType,
	})
	if e != nil {
		t.Error(e)
		return
	}
	b.disableOutput()
	b.bombard()
}
