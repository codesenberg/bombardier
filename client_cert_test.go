package main

import (
	"bytes"
	"crypto/tls"
	"errors"
	"io/ioutil"
	"log"
	"testing"
)

var (
	testCertKey = []byte(`-----BEGIN PRIVATE KEY-----
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
EJMN4Ww1GVMYvIfdK7cSpDU04LWOKX0=
-----END PRIVATE KEY-----
`)
	testCert = []byte(`-----BEGIN CERTIFICATE-----
MIIFXTCCA0WgAwIBAgIJAMMAnp6cISAqMA0GCSqGSIb3DQEBCwUAMEUxCzAJBgNV
BAYTAkFVMRMwEQYDVQQIDApTb21lLVN0YXRlMSEwHwYDVQQKDBhJbnRlcm5ldCBX
aWRnaXRzIFB0eSBMdGQwHhcNMTYxMDI4MTIwMDUxWhcNMTcxMDI4MTIwMDUxWjBF
MQswCQYDVQQGEwJBVTETMBEGA1UECAwKU29tZS1TdGF0ZTEhMB8GA1UECgwYSW50
ZXJuZXQgV2lkZ2l0cyBQdHkgTHRkMIICIjANBgkqhkiG9w0BAQEFAAOCAg8AMIIC
CgKCAgEAx15O6ivZlsJjbadoKDbomEjrlt8s5RqobnIMOhJ6+TkIqgmPVdrD1g2V
HwUhWeiojT9k2TYjPsEHGRLUADfjVfHDp8ZIn6we/UQoYzwjpiIl4kpVeQlYK6OA
Ly9xfMKO1nB6kIKatB1Xr4C/K+pHoVrSXrtiKRJZLU9m2Ff4RhxpbjzV5aK3gPXU
EkKw/plzmRryKza2ZTiedgxikESupuvsHpOcJL/+QfExQp01BNcVOgUm1fH1qWls
FgBt1JRlGeKUU3RtPp4h1eKTXPBtqW+m6/f1ElH2ivfFWHw8RH6i1ideR7an3wJr
4h6S1BBPK0CL/Aicq0eC1Mx8QjEaQW6Kdgk691PNi7yXqWinvZnjg2Z/ZnB/rRST
CPCmNSFqb+FCazFHTs++9LQiQ2IJecjRA5IKicwKetl2L4aFek28LfzNwVxjlq9A
xiCvfytks4sXTkbXlNWh1DrndOEo85H4GS+Z13qbkr7l55EYh6KsdQH9RltbH/td
sH/xXdudlW+9X0HTGq9iEc2RX+wiii/GWuaOsKphgSHPA8UeA0HzghX9lLT9urGk
zWbltcNLAxs2pXOhPXpOqcPpISEdF+y/TmOV1lJErcQGGkZKQXNpwA+nu5Dwie2n
U6MqV+faurW+8buH2WO5uWDJbvEds0y62DWdd8XXx4YJ/zK6Z1kCAwEAAaNQME4w
HQYDVR0OBBYEFBB2lPL2uqy9z4xvVV6bla1AePuJMB8GA1UdIwQYMBaAFBB2lPL2
uqy9z4xvVV6bla1AePuJMAwGA1UdEwQFMAMBAf8wDQYJKoZIhvcNAQELBQADggIB
AG5acpIl15s3u9JpPRShwXGb2ti3o1jbx3riKJAY8jJZTSxQeFARO4GQL+C0Z8+j
arleCxDZoqq+xDgVITCBHZzkyp9wLbihi8EXMYPbrY0Ncv9IoE5A0Oh2c4iRlQAZ
KJy2GNgq8aOKzVo4E44aEPfuzgqaGFhX5SBCBanr+cddKQmdsj/XMd+3FGPTZPkZ
LhG/G7o+tkhewBcd7ZeTSmmYSPTnD6U+12N/UDwzcJoSfNx65j2S8BYhsXv4AnK+
k1HHqWZEDLW6F8oG4RZb9o1sKG3g65vho0HdZXapPVJ00wczXHAD1JyBm5xZwokM
MJUT7L+zBktMY4O0YsjiKc9WfeqcnSUF4wQ5p9+DM1eawYTdN1cpmOt7DYohUVcu
fbgHV2xUJJr3cmvgctEukyp1vbVA2XSQW+cRVZCOxfeC/OKdRR9G76cbM5175DMj
1CSs8M5pxMGl/arUKrEULA3/9QVbd2BBhg4uJ0DdvczqgYG9PgaoLvb2Ej9Rv/F5
kGf1sKdsBbTpVZWPLWXsBOamoSu1ufm6AYmlTdRqrvn9QeK19MNw12YFkFq3lwpN
UmsdOQnlJKUaLwGg3l6/qdZ+YM8S4lYH4ISR8PbknQ8MukUyEmRypjBckFuoMMhW
GOJqhwmTIXWxyg2JXftsm2geg5GHjFXjcZ5dhkrGegrR
-----END CERTIFICATE-----`)
)

func TestGenerateTLSConfigSuccess(t *testing.T) {
	// mocked ioutil readfile
	ioutilReadFile = func(filename string) ([]byte, error) {
		return append(testCertKey, testCert...), nil
	}
	defer func() { ioutilReadFile = ioutil.ReadFile }()

	expectations := []struct {
		in   config
		test func(tlsConfig *tls.Config) error
	}{
		{
			config{
				clientCert: "",
				insecure:   true,
			},
			func(tlsConfig *tls.Config) error {
				if tlsConfig.Certificates != nil {
					return errors.New("TLS Config Certificates should be nil when passed in an empty client certificate filename")
				}
				if !tlsConfig.InsecureSkipVerify {
					return errors.New("TLS Config InsecureSkipVerify should be true when passed in a true insecure config")
				}
				return nil
			},
		},
		{
			config{
				clientCert: "something",
				insecure:   true,
			},
			func(tlsConfig *tls.Config) error {
				// expected result
				shouldBe, err := tls.X509KeyPair(testCert, testCertKey)
				if err != nil {
					t.Error("failed to load X509KeyPair for validation of test")
				}

				for i, v := range tlsConfig.Certificates[0].Certificate {
					if bytes.Compare(v, shouldBe.Certificate[i]) != 0 {
						t.Error("Certificate in tls.Config is not what it should be")
					}
				}

				if !tlsConfig.InsecureSkipVerify {
					return errors.New("TLS Config InsecureSkipVerify should be true when passed in a true insecure config")
				}
				return nil
			},
		},
	}
	for _, e := range expectations {
		if err := e.test(generateTLSConfig(e.in)); err != nil {
			t.Error(err.Error())
		}
	}
}

func TestGenerateTLSConfigFailures(t *testing.T) {
	// mocked ioutil readfile
	ioutilReadFile = func(filename string) ([]byte, error) {
		return nil, errors.New("failure")
	}
	defer func() { ioutilReadFile = ioutil.ReadFile }()

	// mocked ioutil readfile
	tlsX509KeyPair = func(certPEMBlock, keyPEMBlock []byte) (tls.Certificate, error) {
		return tls.Certificate{}, errors.New("failure")
	}
	defer func() { tlsX509KeyPair = tls.X509KeyPair }()

	var (
		errCount        = 0
		fileErrCount    = 0
		keyPairErrCount = 0
	)

	// mocked fatalf, which will keep track of our seen fatalf
	logFatalf = func(format string, v ...interface{}) {
		errCount++
		if format == "failed to read client certificate file: %v" {
			fileErrCount++
		}
		if format == "unable to load client cert and key pair: %v" {
			keyPairErrCount++
		}

	}
	defer func() { logFatalf = log.Fatalf }()

	generateTLSConfig(config{clientCert: "something"})

	if errCount != 2 {
		t.Error("incorrect number of log.Fatalf's witnessed: ", errCount)
	}
	if fileErrCount != 1 {
		t.Error("incorrect number of log.Fatalf's witnessed for file read: ", fileErrCount)
	}
	if keyPairErrCount != 1 {
		t.Error("incorrect number of log.Fatalf's witnessed for loading keypair: ", keyPairErrCount)
	}
}
