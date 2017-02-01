package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"time"
)

var (
	errNoURL       = errors.New("No URL supplied")
	errTooManyArgs = errors.New("Too many arguments are supplied")

	emptyConf = config{}
	parser    = newDefaultParser()
)

type defaultParser struct {
	fs *flag.FlagSet

	numReqs   *nullableUint64
	duration  *nullableDuration
	headers   *headersList
	numConns  uint64
	timeout   time.Duration
	latencies bool
	insecure  bool
	method    string
	body      string
	certPath  string
	keyPath   string
}

func newDefaultParser() *defaultParser {
	dp := new(defaultParser)
	dp.fs = flag.NewFlagSet(programName, flag.ContinueOnError)
	dp.numReqs = new(nullableUint64)
	dp.duration = new(nullableDuration)
	dp.headers = new(headersList)
	dp.fs.Uint64Var(&dp.numConns, "c", defaultNumberOfConns,
		"Maximum number of concurrent connections")
	dp.fs.DurationVar(&dp.timeout, "timeout", defaultTimeout,
		"Socket/request timeout")
	dp.fs.BoolVar(&dp.latencies, "latencies", false,
		"Print latency statistics")
	dp.fs.StringVar(&dp.method, "m", "GET",
		"Request method")
	dp.fs.StringVar(&dp.body, "data", "",
		"Request body")
	dp.fs.StringVar(&dp.certPath, "cert", "",
		"Path to the client's TLS Certificate")
	dp.fs.StringVar(&dp.keyPath, "key", "",
		"Path to the client's TLS Certificate Private Key")
	dp.fs.BoolVar(&dp.insecure, "insecure", false,
		"Controls whether a client verifies the server's certificate"+
			" chain and host name")
	dp.fs.Var(dp.headers, "H", "HTTP headers to use(can be repeated)")
	dp.fs.Var(dp.numReqs, "n", "Number of requests")
	dp.fs.Var(dp.duration, "d", "Duration of test")
	return dp
}

func (p *defaultParser) usage(out io.Writer) {
	fmt.Fprintln(out, programName, "[options] <url>")
	p.fs.SetOutput(out)
	p.fs.PrintDefaults()
}

func (p *defaultParser) parse(args []string) (config, error) {
	err := p.parseArgs(args)
	if err != nil {
		return emptyConf, err
	}
	if p.fs.NArg() == 0 {
		return emptyConf, errNoURL
	}
	if p.fs.NArg() > 1 {
		return emptyConf, errTooManyArgs
	}
	return config{
		numConns:       p.numConns,
		numReqs:        p.numReqs.val,
		duration:       p.duration.val,
		url:            p.fs.Arg(0),
		headers:        p.headers,
		timeout:        p.timeout,
		method:         p.method,
		body:           p.body,
		keyPath:        p.keyPath,
		certPath:       p.certPath,
		printLatencies: p.latencies,
		insecure:       p.insecure,
	}, nil
}

func (p *defaultParser) parseArgs(args []string) error {
	p.fs.SetOutput(ioutil.Discard)
	err := p.fs.Parse(args[1:])
	p.fs.SetOutput(os.Stdout)
	return err
}
