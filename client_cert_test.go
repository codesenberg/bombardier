package main

import (
	"testing"
)

func TestGenerateTLSConfig(t *testing.T) {
	expectations := []struct {
		certPath string
		keyPath  string
		errIsNil bool
	}{
		{
			certPath: "testclient.cert",
			keyPath:  "testclient.key",
			errIsNil: true,
		},
		{
			certPath: "doesnotexist.pem",
			keyPath:  "doesnotexist.pem",
			errIsNil: false,
		},
		{
			certPath: "",
			keyPath:  "",
			errIsNil: true,
		},
	}
	for _, e := range expectations {
		_, r := generateTLSConfig(
			config{
				url:      ParseURLOrPanic("https://doesnt.exist.com"),
				certPath: e.certPath,
				keyPath:  e.keyPath,
			},
		)
		if (r == nil) != e.errIsNil {
			t.Error(e.certPath, e.keyPath, r)
		}
	}
}

func TestInvalidKeyLogPath(t *testing.T) {
	expectations := []struct {
		keyLogPath string
		errIsNil   bool
	}{
		{
			keyLogPath: "/path/to/invalid/log",
			errIsNil:   false,
		},
	}
	for _, e := range expectations {
		_, err := generateTLSConfig(
			config{
				url:        ParseURLOrPanic("https://doesnt.exist.com"),
				keyLogPath: e.keyLogPath,
			},
		)
		if (err == nil) != e.errIsNil {
			t.Error(e.keyLogPath, err)
		}
	}
}
