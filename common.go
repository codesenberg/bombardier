package main

import (
	"errors"
	"sort"
	"time"
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
	cantHaveBody = []string{"GET", "HEAD"}

	errInvalidURL = errors.New(
		"No hostname or invalid scheme")
	errInvalidNumberOfConns = errors.New(
		"Invalid number of connections(must be > 0)")
	errInvalidNumberOfRequests = errors.New(
		"Invalid number of requests(must be > 0)")
	errInvalidTestDuration = errors.New(
		"Invalid test duration(must be >= 1s)")
	errNegativeTimeout = errors.New(
		"Timeout can't be negative")
	errBodyNotAllowed = errors.New(
		"GET and HEAD requests cannot have body")
	errNoPathToCert = errors.New(
		"No Path to TLS Client Certificate")
	errNoPathToKey = errors.New(
		"No Path to TLS Client Certificate Private Key")
	errZeroRate = errors.New(
		"Rate can't be less than 1")
	errBodyProvidedTwice = errors.New("Use either --body or --body-file")

	errInvalidHeaderFormat = errors.New("Invalid header format")
	errEmptyPrintSpec      = errors.New(
		"Empty print spec is not a valid print spec")
)

func init() {
	sort.Strings(httpMethods)
	sort.Strings(cantHaveBody)
}
