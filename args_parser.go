package main

import (
	"runtime"
	"strconv"
	"time"

	"github.com/alecthomas/kingpin"
)

const (
	decBase = 10
)

var (
	emptyConf = config{}
	parser    = newKingpinParser()
)

type argsParser interface {
	parse([]string) (config, error)
}

type kingpinParser struct {
	app *kingpin.Application

	url string

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

func newKingpinParser() argsParser {
	kparser := &kingpinParser{
		numReqs:   new(nullableUint64),
		duration:  new(nullableDuration),
		headers:   new(headersList),
		numConns:  defaultNumberOfConns,
		timeout:   defaultTimeout,
		latencies: false,
		method:    "GET",
		body:      "",
		certPath:  "",
		keyPath:   "",
		insecure:  false,
		url:       "",
	}

	app := kingpin.New("", "Fast cross-platform HTTP benchmarking tool").
		Version("bombardier version " + version + " " + runtime.GOOS + "/" +
			runtime.GOARCH)
	app.Flag("connections", "Maximum number of concurrent connections").
		Short('c').
		PlaceHolder(strconv.FormatUint(defaultNumberOfConns, decBase)).
		Uint64Var(&kparser.numConns)
	app.Flag("timeout", "Socket/request timeout").
		PlaceHolder(defaultTimeout.String()).
		Short('t').
		DurationVar(&kparser.timeout)
	app.Flag("latencies", "Print latency statistics").
		Short('l').
		BoolVar(&kparser.latencies)
	app.Flag("method", "Request method").
		PlaceHolder("GET").
		Short('m').
		StringVar(&kparser.method)
	app.Flag("body", "Request body").
		Default("").
		Short('b').
		StringVar(&kparser.body)
	app.Flag("cert", "Path to the client's TLS Certificate").
		Default("").
		StringVar(&kparser.certPath)
	app.Flag("key", "Path to the client's TLS Certificate Private Key").
		Default("").
		StringVar(&kparser.keyPath)
	app.Flag("insecure",
		"Controls whether a client verifies the server's certificate"+
			" chain and host name").
		Short('k').
		BoolVar(&kparser.insecure)

	app.Flag("header", "HTTP headers to use(can be repeated)").
		PlaceHolder("[]").
		Short('H').
		SetValue(kparser.headers)
	app.Flag("requests", "Number of requests").
		PlaceHolder("[<pos. int.>]").
		Short('n').
		SetValue(kparser.numReqs)
	app.Flag("duration", "Duration of test").
		PlaceHolder(defaultTestDuration.String()).
		Short('d').
		SetValue(kparser.duration)

	app.Arg("url", "Target's URL").Required().
		StringVar(&kparser.url)

	kparser.app = app
	return argsParser(kparser)
}

func (k *kingpinParser) parse(args []string) (config, error) {
	k.app.Name = args[0]
	_, err := k.app.Parse(args[1:])
	if err != nil {
		return emptyConf, err
	}
	return config{
		numConns:       k.numConns,
		numReqs:        k.numReqs.val,
		duration:       k.duration.val,
		url:            k.url,
		headers:        k.headers,
		timeout:        k.timeout,
		method:         k.method,
		body:           k.body,
		keyPath:        k.keyPath,
		certPath:       k.certPath,
		printLatencies: k.latencies,
		insecure:       k.insecure,
	}, nil
}
