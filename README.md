# bombardier [![Build Status](https://travis-ci.org/codesenberg/bombardier.svg?branch=master)](https://travis-ci.org/codesenberg/bombardier) [![Go Report Card](https://goreportcard.com/badge/github.com/codesenberg/bombardier)](https://goreportcard.com/report/github.com/codesenberg/bombardier) [![GoDoc](https://godoc.org/github.com/codesenberg/bombardier?status.svg)](http://godoc.org/github.com/codesenberg/bombardier) [![Coverage](https://gocover.io/_badge/github.com/codesenberg/bombardier)](https://gocover.io/github.com/codesenberg/bombardier)
bombardier is a HTTP(S) benchmarking tool. It is written in Go programming language and uses excellent [fasthttp](https://github.com/valyala/fasthttp) instead of Go's default http library, because of its lightning fast performance.

Tested on go1.6 and higher. Use go1.7+ for best performance.

## Installation
You can grab binaries in the [releases](https://github.com/codesenberg/bombardier/releases) section.
Alternatively, to get latest and greatest run:

`go get -u github.com/codesenberg/bombardier`

## Usage
```
bombardier [<flags>] <url>
```

For a more detailed information about flags consult [GoDoc](http://godoc.org/github.com/codesenberg/bombardier).

## Examples
Example of running `bombardier` against [this server](https://godoc.org/github.com/codesenberg/bombardier/cmd/utils/simplebenchserver):
```
> bombardier -c 125 -n 10000000 http://localhost:8080
Bombarding http://localhost:8080 with 10000000 requests using 125 connections
 10000000 / 10000000 [============================================] 100.00% 37s Done!
Statistics        Avg      Stdev        Max
  Reqs/sec    264560.00   10733.06     268434
  Latency      471.00us   522.34us    51.00ms
  HTTP codes:
    1xx - 0, 2xx - 10000000, 3xx - 0, 4xx - 0, 5xx - 0
    others - 0
  Throughput:   292.92MB/s
```
Or, against a realworld server(with latency distribution):
```
> bombardier -c 200 -d 10s -l http://ya.ru
Bombarding http://ya.ru for 10s using 200 connections
[=========================================================================] 10s Done!
Statistics        Avg      Stdev        Max
  Reqs/sec      6607.00     524.56       7109
  Latency       29.86ms     5.36ms   305.02ms
  Latency Distribution
     50%    28.00ms
     75%    32.00ms
     90%    34.00ms
     99%    48.00ms
  HTTP codes:
    1xx - 0, 2xx - 0, 3xx - 66561, 4xx - 0, 5xx - 0
    others - 5
  Errors:
    dialing to the given TCP address timed out - 5
  Throughput:     3.06MB/s
```
