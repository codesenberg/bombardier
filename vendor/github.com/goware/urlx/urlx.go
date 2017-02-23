// Package urlx parses and normalizes URLs. It can also resolve hostname to an IP address.
package urlx

import (
	"errors"
	"net"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/purell"
	"golang.org/x/net/idna"
)

// Parse parses raw URL string into the net/url URL struct.
// It uses the url.Parse() internally, but it slightly changes
// its behavior:
// 1. It forces the default scheme and port.
// 2. It favors absolute paths over relative ones, thus "example.com"
//    is parsed into url.Host instead of url.Path.
// 4. It lowercases the Host (not only the Scheme).
func Parse(rawURL string) (*url.URL, error) {
	// Force default http scheme, so net/url.Parse() doesn't
	// put both host and path into the (relative) path.
	if strings.Index(rawURL, "//") == 0 {
		// Leading double slashes (any scheme). Force http.
		rawURL = "http:" + rawURL
	}
	if strings.Index(rawURL, "://") == -1 {
		// Missing scheme. Force http.
		rawURL = "http://" + rawURL
	}

	// Use net/url.Parse() now.
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, err
	}

	host, _, err := SplitHostPort(u)
	if err != nil {
		return nil, err
	}
	if err := checkHost(host); err != nil {
		return nil, err
	}

	u.Host = strings.ToLower(u.Host)
	u.Scheme = strings.ToLower(u.Scheme)

	return u, nil
}

var (
	// RFC 1035.
	domainRegexp = regexp.MustCompile(`^([a-zA-Z0-9-]{1,63}\.)+[a-zA-Z0-9][a-zA-Z0-9-]{0,61}[a-zA-Z0-9]$`)
	ipv4Regexp   = regexp.MustCompile(`^[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}$`)
	ipv6Regexp   = regexp.MustCompile(`^\[[a-fA-F0-9:]+\]$`)
)

func checkHost(host string) error {
	if host == "" {
		return &url.Error{"host", host, errors.New("empty host")}
	}

	host = strings.ToLower(host)
	if domainRegexp.MatchString(host) || host == "localhost" {
		return nil
	}

	if punycode, err := idna.ToASCII(host); err != nil {
		return err
	} else if domainRegexp.MatchString(punycode) {
		return nil
	}

	// IPv4 and IPv6.
	if ipv4Regexp.MatchString(host) || ipv6Regexp.MatchString(host) {
		return nil
	}

	return &url.Error{"host", host, errors.New("invalid host")}
}

// SplitHostPort splits network address of the form "host:port" into
// host and port. Unlike net.SplitHostPort(), it doesn't remove brackets
// from [IPv6] host and it accepts net/url.URL struct instead of a string.
func SplitHostPort(u *url.URL) (host, port string, err error) {
	if u == nil {
		return "", "", &url.Error{"host", host, errors.New("empty url")}
	}
	host = u.Host

	// Find last colon.
	if i := strings.LastIndex(host, ":"); i != -1 {
		// If we're not inside [IPv6] brackets, split host:port.
		if len(host) > i && strings.Index(host[i:], "]") == -1 {
			port = host[i+1:]
			host = host[:i]
		}
	}

	// Port is optional. But if it's set, is it a number?
	if port != "" {
		if _, err := strconv.Atoi(port); err != nil {
			return "", "", &url.Error{"port", host, err}
		}
	}

	return host, port, nil
}

const normalizeFlags purell.NormalizationFlags = purell.FlagRemoveDefaultPort |
	purell.FlagDecodeDWORDHost | purell.FlagDecodeOctalHost | purell.FlagDecodeHexHost |
	purell.FlagRemoveUnnecessaryHostDots | purell.FlagRemoveDotSegments | purell.FlagRemoveDuplicateSlashes |
	purell.FlagUppercaseEscapes | purell.FlagDecodeUnnecessaryEscapes | purell.FlagEncodeNecessaryEscapes |
	purell.FlagSortQuery

// Normalize returns normalized URL string.
// Behavior:
// 1. Remove unnecessary host dots.
// 2. Remove default port (http://localhost:80 becomes http://localhost).
// 3. Remove duplicate slashes.
// 4. Remove unnecessary dots from path.
// 5. Sort query parameters.
// 6. Decode host IP into decimal numbers.
// 7. Handle escape values.
// 8. Decode Punycode domains into UTF8 representation.
func Normalize(u *url.URL) (string, error) {
	host, port, err := SplitHostPort(u)
	if err != nil {
		return "", err
	}
	if err := checkHost(host); err != nil {
		return "", err
	}

	// Decode Punycode.
	host, err = idna.ToUnicode(host)
	if err != nil {
		return "", err
	}

	u.Host = strings.ToLower(host)
	if port != "" {
		u.Host += ":" + port
	}
	u.Scheme = strings.ToLower(u.Scheme)

	return purell.NormalizeURL(u, normalizeFlags), nil
}

// NormalizeString returns normalized URL string.
// It's a shortcut for Parse() and Normalize() funcs.
func NormalizeString(rawURL string) (string, error) {
	u, err := Parse(rawURL)
	if err != nil {
		return "", err
	}

	return Normalize(u)
}

// Resolve resolves the URL host to its IP address.
func Resolve(u *url.URL) (*net.IPAddr, error) {
	host, _, err := SplitHostPort(u)
	if err != nil {
		return nil, err
	}

	addr, err := net.ResolveIPAddr("ip", host)
	if err != nil {
		return nil, err
	}

	return addr, nil
}

// Resolve resolves the URL host to its IP address.
// It's a shortcut for Parse() and Resolve() funcs.
func ResolveString(rawURL string) (*net.IPAddr, error) {
	u, err := Parse(rawURL)
	if err != nil {
		return nil, err
	}
	return Resolve(u)
}

func URIEncode(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
