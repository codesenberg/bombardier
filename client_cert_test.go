package main

import "testing"

func TestGenerateTLSConfig(t *testing.T) {
	expectations := []struct {
		certPath string
		keyPath  string
		errTest  func(error) bool
	}{
		{
			certPath: "test_cert.pem",
			keyPath:  "test_key.pem",
			errTest: func(err error) bool {
				if err == nil {
					return true
				}
				return false
			},
		},
		{
			certPath: "doesnotexist.pem",
			keyPath:  "doesnotexist.pem",
			errTest: func(err error) bool {
				if err != nil {
					return true
				}
				return false
			},
		},
		{
			certPath: "",
			keyPath:  "",
			errTest: func(err error) bool {
				if err == nil {
					return true
				}
				return false
			},
		},
	}
	for _, e := range expectations {
		if _, r := generateTLSConfig(config{certPath: e.certPath, keyPath: e.keyPath}); !e.errTest(r) {
			t.Log(e.certPath, e.keyPath, r)
			t.Fail()
		}
	}
}
