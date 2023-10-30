/*
Package template documents the way user-defined output templates are
ment to be used.

User-defined templates use Go's text/template package, so you might
want to check its documentation first.
There are a bunch of helper methods available inside a template
besides those described in aforementioned documentation, namely:
  - URLString()
    Returns the URL string used for the load test.
  - WithLatencies()
    Tells whether --latencies flag were activated.
  - FormatBinary(numberOfBytes float64) string
    Converts bytes to kilo-, mega-, giga-, etc.- bytes, and
    appends appropriate suffix "KB", "MB", "GB", etc.
  - FormatTimeUs(us float64) string
    Converts microseconds to milliseconds, seconds, minutes or
    hours and appends appropriate suffix.
  - FormatTimeUsUint64(us uint64) string
    Same as above, but for uint64, since type conversions are
    not available in templates.
  - FloatsToArray(ps ...float64) []float64
    Converts a bunch of floats into array, since, again,
    type conversions are not available in templates.
  - Multiply(num, coeff float64) float64
    Arithmetics are not available inside of templates either.
  - StringToBytes(s string) []byte
    Convenience function to convert string to []byte.
  - UUIDV1() (UUID, error)
    Generates UUID Version 1, based on timestamp and
    MAC address (RFC 4122)
  - UUIDV2(domain byte) (UUID, error)
    Generates UUID Version 2, based on timestamp, MAC address
    and POSIX UID/GID (DCE 1.1)
  - UUIDV3(ns UUID, name string) UUID
    Generates UUID Version 3, based on MD5 hashing (RFC 4122)
  - UUIDV4() (UUID, error)
    Generates UUID Version 4, based on random numbers (RFC 4122)
  - UUIDV5(ns UUID, name string) UUID
    Generates UUID Version 5, based on SHA-1 hashing (RFC 4122)

The structure that gets passed to the template is documented in
the package github.com/codesenberg/bombardier/internal. The structure
of interest is TestInfo. It basically consists of Spec and Result
fields, the former contains various information about the test
(number of connections, URL, HTTP method, headers, body, rate, etc.)
performed, while the latter contains results obtained during the
execution of this test (bytes read/written, time taken, RPS, etc.).

Link to GoDoc for the structure used in template:
https://godoc.org/github.com/codesenberg/bombardier/internal#TestInfo

Examples of templates can be found in:
https://github.com/codesenberg/bombardier/blob/master/templates.go
*/
package template
