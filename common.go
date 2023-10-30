package main

import (
	"errors"
	"net/url"
	"sort"
	"time"

	"github.com/goware/urlx"
)

const (
	decBase = 10

	rateLimitInterval = 10 * time.Millisecond
	oneSecond         = 1 * time.Second

	exitFailure = 1
)

var (
	version = "unspecified"

	emptyConf = config{}
	parser    = newKingpinParser()

	defaultTestDuration  = 10 * time.Second
	defaultNumberOfConns = uint64(125)
	defaultTimeout       = 2 * time.Second

	httpMethods = []string{
		"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS",
		"PATCH",
	}
	cantHaveBody = []string{"HEAD"}

	errUnsupportedScheme    = errors.New("unsupported scheme")
	errInvalidNumberOfConns = errors.New(
		"invalid number of connections(must be > 0)")
	errInvalidNumberOfRequests = errors.New(
		"invalid number of requests(must be > 0)")
	errInvalidTestDuration = errors.New(
		"invalid test duration(must be >= 1s)")
	errNegativeTimeout = errors.New(
		"timeout can't be negative")
	errBodyNotAllowed = errors.New(
		"HEAD requests cannot have body")
	errNoPathToCert = errors.New(
		"no Path to TLS Client Certificate")
	errNoPathToKey = errors.New(
		"no Path to TLS Client Certificate Private Key")
	errZeroRate = errors.New(
		"rate can't be less than 1")
	errBodyProvidedTwice = errors.New("use either --body or --body-file")

	errInvalidHeaderFormat = errors.New("invalid header format")
	errEmptyPrintSpec      = errors.New(
		"empty print spec is not a valid print spec")
)

func ParseURLOrPanic(s string) *url.URL {
	u, err := urlx.Parse(s)
	if err != nil {
		panic(err)
	}
	return u
}

func init() {
	sort.Strings(httpMethods)
	sort.Strings(cantHaveBody)
}
