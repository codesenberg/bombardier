package main

import "testing"

func TestGenerateTLSConfig(t *testing.T) {
	expectations := []struct {
		certPath string
		keyPath  string
		errIsNil bool
	}{
		{
			certPath: "test_cert.pem",
			keyPath:  "test_key.pem",
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
		_, r := generateTLSConfig(config{url: "https://doesnt.exist.com", certPath: e.certPath, keyPath: e.keyPath})
		if (r == nil) != e.errIsNil {
			t.Log(e.certPath, e.keyPath, r)
			t.Fail()
		}
	}
}
