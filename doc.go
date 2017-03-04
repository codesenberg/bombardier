/*
Command line utility bombardier is a fast cross-platform HTTP
benchmarking tool written in Go.

Installation:
  go get -u github.com/codesenberg/bombardier

Usage:
  bombardier [<flags>] <url>

Flags:
      --help                  Show context-sensitive help (also try --help-long
                              and --help-man).
      --version               Show application version.
  -c, --connections=125       Maximum number of concurrent connections
  -t, --timeout=2s            Socket/request timeout
  -l, --latencies             Print latency statistics
  -m, --method=GET            Request method
  -b, --body=""               Request body
      --cert=""               Path to the client's TLS Certificate
      --key=""                Path to the client's TLS Certificate Private Key
  -k, --insecure              Controls whether a client verifies the server's
                              certificate chain and host name
  -H, --header="K: V" ...     HTTP headers to use(can be repeated)
  -n, --requests=[pos. int.]  Number of requests
  -d, --duration=10s          Duration of test
Args:
  <url>  Target's URL
*/
package main
