/*
Command line utility bombardier is a fast crossplatform HTTP
benchmarking tool written in Go.

Installation:
  go get -u github.com/codesenberg/bombardier

Usage:
  bombardier [options] <url>

Options available:
  -H value
        HTTP headers to use(can be repeated)
  -c uint
        Maximum number of concurrent connections (default 125)
  -cert string
        Path to the client's TLS Certificate
  -d value
        Duration of test
  -data string
        Request body
  -insecure
        Controls whether a client verifies the server's certificate chain
        and host name (default true)
  -key string
        Path to the client's TLS Certificate Private Key
  -latencies
        Print latency statistics
  -m string
        Request method (default "GET")
  -n value
        Number of requests
  -timeout duration
        Socket/request timeout (default 2s)
*/
package main
