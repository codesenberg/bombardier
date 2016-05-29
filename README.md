# bombardier
bombardier is a HTTP(S) benchmarking tool. It's written in Go programming language and uses excellent [fasthttp](https://github.com/valyala/fasthttp) instead of Go's default http library, because of it's lightning fast performance.

##Installation
You are encourages to grab the latest version of the tool in the [releases](https://github.com/bugsenberg/bombardier/releases) section and test the tool by yourself.

#### If you can't find your OS/ARCH combo(aka build from the source)
This one is actually pretty straightforward. Just run `go get github.com/bugsenberg/bombardier`(but you may need to get the deps first).

##Usage
Run it like:
```
bombardier <options> <url>
```
Also, you can supply these options:
```
  -H value
        HTTP headers to use (default [])
  -c int
        Maximum number of concurrent connections (default 200)
  -latencies
        Print latency statistics
  -n int
        Number of requests (default 10000)
  -timeout duration
        Socket/request timeout (default 2s)
```
You should see something like this if you done everything correctly:
```
> bombardier -c 200 -n 10000000 http://localhost:8080
Bombarding http://localhost:8080 with 10000000 requests using 200 connections
10000000 / 10000000 [============================================] 100.00 % 55s 
Statistics        Avg      Stdev        Max
  Reqs/sec    181631.00   13494.01     197924
  Latency        1.10ms   319.69us    82.51ms
  HTTP codes:
    1xx - 0, 2xx - 10000000, 3xx - 0, 4xx - 0, 5xx - 0
    errored - 0
  Throughput:   201.11MB/s
```
Or, on a realworld server(with latency distribution):
```
> bombardier -c 200 -n 10000 --latencies http://google.com
Bombarding http://google.com with 10000 requests using 200 connections
10000 / 10000 [===================================================] 100.00 % 2s 
Statistics        Avg      Stdev        Max
  Reqs/sec      4165.00    1382.95       4939
  Latency       43.14ms    26.01ms   394.05ms
  Latency Distribution
     50%    38.50ms
     75%    44.01ms
     90%    47.01ms
     99%   113.01ms
  HTTP codes:
    1xx - 0, 2xx - 0, 3xx - 9994, 4xx - 0, 5xx - 0
    errored - 6
  Throughput:     1.95MB/s
```