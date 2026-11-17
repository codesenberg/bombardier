package main

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/goware/urlx"
)

func TestShouldReturnNilIfNoHeadersWhereSet(t *testing.T) {
	h := new(headersList)
	if headersToFastHTTPHeaders(h) != nil {
		t.Fail()
	}
}

func TestShouldReturnEmptyHeadersIfNoHeaadersWhereSet(t *testing.T) {
	h := new(headersList)
	if len(headersToHTTPHeaders(h)) != 0 {
		t.Fail()
	}
}

func TestShouldProperlyConvertToHttpHeaders(t *testing.T) {
	h := new(headersList)
	for _, hs := range []string{
		"Content-Type: application/json", "Custom-Header: xxx42xxx",
	} {
		if err := h.Set(hs); err != nil {
			t.Error(err)
		}
	}
	fh := headersToFastHTTPHeaders(h)
	{
		e, a := []byte("application/json"), fh.Peek("Content-Type")
		if !bytes.Equal(e, a) {
			t.Errorf("Expected %v, but got %v", e, a)
		}
	}
	if e, a := []byte("xxx42xxx"), fh.Peek("Custom-Header"); !bytes.Equal(e, a) {
		t.Errorf("Expected %v, but got %v", e, a)
	}

	nh := headersToHTTPHeaders(h)
	{
		e, a := "application/json", nh.Get("Content-Type")
		if e != a {
			t.Errorf("Expected %v, but got %v", e, a)
		}
	}
	if e, a := "xxx42xxx", nh.Get("Custom-Header"); e != a {
		t.Errorf("Expected %v, but got %v", e, a)
	}
}

func TestHTTP2Client(t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewUnstartedServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !r.ProtoAtLeast(2, 0) {
			t.Errorf("invalid HTTP proto version: %v", r.Proto)
		}

		w.WriteHeader(http.StatusOK)
		_, err := w.Write(response)
		if err != nil {
			t.Error(err)
		}
	}))
	s.EnableHTTP2 = true
	s.TLS = &tls.Config{
		InsecureSkipVerify: true,
	}
	s.StartTLS()
	defer s.Close()

	bytesRead, bytesWritten := int64(0), int64(0)
	requestURL, err := urlx.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	c := newHTTPClient(&clientOpts{
		HTTP2: true,

		headers:    new(headersList),
		requestURL: requestURL,
		method:     "GET",
		tlsConfig: &tls.Config{
			InsecureSkipVerify: true,
		},

		body: new(string),

		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	})
	code, _, err := c.do()
	if err != nil {
		t.Error(err)
		return
	}
	if code != http.StatusOK {
		t.Errorf("invalid response code: %v", code)
	}
	if atomic.LoadInt64(&bytesRead) == 0 {
		t.Errorf("invalid response size: %v", bytesRead)
	}
	if atomic.LoadInt64(&bytesWritten) == 0 {
		t.Errorf("empty request of size: %v", bytesWritten)
	}
}

func TestHTTP1Clients(t *testing.T) {
	responseSize := 1024
	response := bytes.Repeat([]byte{'a'}, responseSize)
	s := httptest.NewServer(http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			if r.ProtoMajor != 1 {
				t.Errorf("invalid HTTP proto version: %v", r.Proto)
			}

			w.WriteHeader(http.StatusOK)
			_, err := w.Write(response)
			if err != nil {
				t.Error(err)
			}
		},
	))
	defer s.Close()

	bytesRead, bytesWritten := int64(0), int64(0)
	requestURL, err := urlx.Parse(s.URL)
	if err != nil {
		t.Fatal(err)
	}
	cc := &clientOpts{
		HTTP2: false,

		headers:    new(headersList),
		requestURL: requestURL,
		method:     "GET",

		body: new(string),

		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	}
	clients := []client{
		newHTTPClient(cc),
		newFastHTTPClient(cc),
	}
	for _, c := range clients {
		bytesRead, bytesWritten = 0, 0
		code, _, err := c.do()
		if err != nil {
			t.Error(err)
			return
		}
		if code != http.StatusOK {
			t.Errorf("invalid response code: %v", code)
		}
		if bytesRead == 0 {
			t.Errorf("invalid response size: %v", bytesRead)
		}
		if bytesWritten == 0 {
			t.Errorf("empty request of size: %v", bytesWritten)
		}
	}
}

func TestClientWithNilBodyAndNilBodProd(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	}))
	defer ts.Close()

	u, err := urlx.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	var bytesRead, bytesWritten int64
	cc := &clientOpts{
		headers:      new(headersList),
		requestURL:   u,
		method:       "GET",
		body:         nil,
		bodProd:      nil,
		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	}

	for _, c := range []client{newHTTPClient(cc), newFastHTTPClient(cc)} {
		code, _, err := c.do()
		if err != nil {
			t.Errorf("unexpected error with nil body: %v", err)
		}
		if code != http.StatusOK {
			t.Errorf("expected 200, got %d", code)
		}
	}
}

func TestConcurrentHTTPClientDo(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("response"))
	}))
	defer ts.Close()

	u, err := urlx.Parse(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	hl := new(headersList)
	_ = hl.Set("User-Agent: test-agent")
	_ = hl.Set("Accept: */*")

	var bytesRead, bytesWritten int64
	cc := &clientOpts{
		headers:      hl,
		requestURL:   u,
		method:       "GET",
		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	}

	c := newHTTPClient(cc)
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 5; j++ {
				code, _, err := c.do()
				if err != nil {
					t.Errorf("concurrent do error: %v", err)
				}
				if code != http.StatusOK {
					t.Errorf("expected 200, got %d", code)
				}
			}
		}()
	}
	wg.Wait()
}
