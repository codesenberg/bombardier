package main

import (
	"bytes"
	"crypto/tls"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
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
	c := newHTTPClient(&clientOpts{
		HTTP2: true,

		headers: new(headersList),
		url:     s.URL,
		method:  "GET",
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
	cc := &clientOpts{
		HTTP2: false,

		headers: new(headersList),
		url:     s.URL,
		method:  "GET",

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

func TestHTTP1ClientsWithHeaders(t *testing.T) {
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
	cc := &clientOpts{
		HTTP2: false,

		headers: &headersList{{"One", "Value one"}},
		url:     s.URL,
		method:  "GET",

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
