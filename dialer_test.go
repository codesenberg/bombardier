package main

import (
	"errors"
	"io"
	"net"
	"sync/atomic"
	"testing"
	"time"
)

type mockConn struct {
	net.Conn
	readData []byte
	readErr  error
	writeErr error
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	n = copy(b, m.readData)
	return n, m.readErr
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return len(b), m.writeErr
}

func (m *mockConn) Close() error                       { return nil }
func (m *mockConn) LocalAddr() net.Addr                { return nil }
func (m *mockConn) RemoteAddr() net.Addr               { return nil }
func (m *mockConn) SetDeadline(t time.Time) error      { return nil }
func (m *mockConn) SetReadDeadline(t time.Time) error  { return nil }
func (m *mockConn) SetWriteDeadline(t time.Time) error { return nil }

func TestCountingConnReadEOF(t *testing.T) {
	var bytesRead, bytesWritten int64
	mock := &mockConn{
		readData: []byte("hello world"),
		readErr:  io.EOF,
	}

	cc := &countingConn{
		Conn:         mock,
		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	}

	buf := make([]byte, 32)
	n, err := cc.Read(buf)
	if err != io.EOF {
		t.Fatalf("expected EOF, got %v", err)
	}
	if n != 11 {
		t.Fatalf("expected 11 bytes read, got %d", n)
	}
	if read := atomic.LoadInt64(&bytesRead); read != 11 {
		t.Fatalf("expected bytesRead to be 11 even on EOF, got %d", read)
	}
}

func TestCountingConnWriteError(t *testing.T) {
	var bytesRead, bytesWritten int64
	testErr := errors.New("write failure")
	mock := &mockConn{
		writeErr: testErr,
	}

	cc := &countingConn{
		Conn:         mock,
		bytesRead:    &bytesRead,
		bytesWritten: &bytesWritten,
	}

	n, err := cc.Write([]byte("data"))
	if err != testErr {
		t.Fatalf("expected writeErr, got %v", err)
	}
	if n != 4 {
		t.Fatalf("expected 4 bytes written, got %d", n)
	}
	if written := atomic.LoadInt64(&bytesWritten); written != 4 {
		t.Fatalf("expected bytesWritten to be 4 even on error, got %d", written)
	}
}
