package main

import "io"

type proxyReader struct {
	io.Reader
}
