package main

import (
	"crypto/tls"
)

// readClientCert - helper function to read client certificate
// from pem formatted certPath and keyPath files
func readClientCert(certPath, keyPath string) ([]tls.Certificate, error) {
	if certPath != "" && keyPath != "" {
		// load keypair
		cert, err := tls.LoadX509KeyPair(certPath, keyPath)
		if err != nil {
			return nil, err
		}

		return []tls.Certificate{cert}, nil
	}
	return nil, nil
}

// generateTLSConfig - helper function to generate a TLS configuration based on
// config
func generateTLSConfig(c config) (*tls.Config, error) {
	certs, err := readClientCert(c.certPath, c.keyPath)
	if err != nil {
		return nil, err
	}
	// Disable gas warning, because InsecureSkipVerify may be set to true
	// for the purpose of testing
	/* #nosec */
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.insecure,
		Certificates:       certs,
	}
	return tlsConfig, nil
}
