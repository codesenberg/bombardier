package main

import (
	"crypto/tls"
	"io"
	"os"
)

// readClientCert - helper function to read client certificate
// from pem formatted certPath and keyPath files
func readClientCert(certPath, keyPath string) ([]tls.Certificate, error) {
	// load keypair
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	return []tls.Certificate{cert}, err
}

// generateTLSConfig - helper function to generate a TLS configuration based on
// config
func generateTLSConfig(c config) (*tls.Config, error) {
	var (
		certs        []tls.Certificate
		keyLogWriter io.Writer
		err          error
	)
	// This assumes that the caller has validated that either both or none of
	// the c.certPath and c.keyPath are set.
	if c.certPath != "" && c.keyPath != "" {
		certs, err = readClientCert(c.certPath, c.keyPath)
		if err != nil {
			return nil, err
		}
	}

	if c.keyLogPath != "" {
		f, err := os.OpenFile(c.keyLogPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)
		if err != nil {
			return nil, err
		}

		keyLogWriter = f
	}
	// Disable gas warning, because InsecureSkipVerify may be set to true
	// for the purpose of testing
	/* #nosec */
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.insecure,
		Certificates:       certs,
		KeyLogWriter:       keyLogWriter,
	}
	return tlsConfig, nil
}
