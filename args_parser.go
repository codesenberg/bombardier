package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alecthomas/kingpin"
)

type argsParser interface {
	parse([]string) (config, error)
}

type kingpinParser struct {
	app *kingpin.Application

	url string

	numReqs      *nullableUint64
	duration     *nullableDuration
	headers      *headersList
	numConns     uint64
	timeout      time.Duration
	latencies    bool
	insecure     bool
	method       string
	body         string
	bodyFilePath string
	stream       bool
	certPath     string
	keyPath      string
	rate         *nullableUint64
	clientType   clientTyp

	printSpec *nullableString
}

func newKingpinParser() argsParser {
	kparser := &kingpinParser{
		numReqs:      new(nullableUint64),
		duration:     new(nullableDuration),
		headers:      new(headersList),
		numConns:     defaultNumberOfConns,
		timeout:      defaultTimeout,
		latencies:    false,
		method:       "GET",
		body:         "",
		bodyFilePath: "",
		stream:       false,
		certPath:     "",
		keyPath:      "",
		insecure:     false,
		url:          "",
		rate:         new(nullableUint64),
		clientType:   fhttp,
		printSpec:    new(nullableString),
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
	app.Flag("body-file", "File to use as request body").
		Default("").
		Short('f').
		StringVar(&kparser.bodyFilePath)
	app.Flag("stream", "Specify whether to stream body using "+
		"chunked transfer encoding or to serve it from memory").
		Short('s').
		BoolVar(&kparser.stream)
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
		PlaceHolder("\"K: V\"").
		Short('H').
		SetValue(kparser.headers)
	app.Flag("requests", "Number of requests").
		PlaceHolder("[pos. int.]").
		Short('n').
		SetValue(kparser.numReqs)
	app.Flag("duration", "Duration of test").
		PlaceHolder(defaultTestDuration.String()).
		Short('d').
		SetValue(kparser.duration)

	app.Flag("rate", "Rate limit in requests per second").
		PlaceHolder("[pos. int.]").
		Short('r').
		SetValue(kparser.rate)

	app.Flag("fasthttp", "Use fasthttp client").
		Action(func(*kingpin.ParseContext) error {
			kparser.clientType = fhttp
			return nil
		}).
		Bool()
	app.Flag("http1", "Use net/http client with forced HTTP/1.x").
		Action(func(*kingpin.ParseContext) error {
			kparser.clientType = nhttp1
			return nil
		}).
		Bool()
	app.Flag("http2", "Use net/http client with enabled HTTP/2.0").
		Action(func(*kingpin.ParseContext) error {
			kparser.clientType = nhttp2
			return nil
		}).
		Bool()

	app.Flag(
		"print", "Specifies what to output. Comma-separated list of values"+
			" 'intro' (short: 'i'), 'progress' (short: 'p'),"+
			" 'result' (short: 'r'). Examples:"+
			"\n\t* i,p,r (prints everything)"+
			"\n\t* intro,result (intro & result)"+
			"\n\t* r (result only)"+
			"\n\t* result (same as above)").
		PlaceHolder("<spec>").
		Short('p').
		SetValue(kparser.printSpec)

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
	pi, pp, pr := true, true, true
	if k.printSpec.val != nil {
		pi, pp, pr, err = parsePrintSpec(*k.printSpec.val)
		if err != nil {
			return emptyConf, err
		}
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
		bodyFilePath:   k.bodyFilePath,
		stream:         k.stream,
		keyPath:        k.keyPath,
		certPath:       k.certPath,
		printLatencies: k.latencies,
		insecure:       k.insecure,
		rate:           k.rate.val,
		clientType:     k.clientType,
		printIntro:     pi,
		printProgress:  pp,
		printResult:    pr,
	}, nil
}

func parsePrintSpec(spec string) (bool, bool, bool, error) {
	pi, pp, pr := false, false, false
	if spec == "" {
		return false, false, false, errEmptyPrintSpec
	}
	parts := strings.Split(spec, ",")
	partsCount := 0
	for _, p := range parts {
		switch p {
		case "i", "intro":
			pi = true
		case "p", "progress":
			pp = true
		case "r", "result":
			pr = true
		default:
			return false, false, false,
				fmt.Errorf("%q is not a valid part of print spec", p)
		}
		partsCount++
	}
	if partsCount < 1 || partsCount > 3 {
		return false, false, false,
			fmt.Errorf(
				"Spec %q has too many parts, at most 3 are allowed", spec,
			)
	}
	return pi, pp, pr, nil
}
