package main

import "crypto/tls"

var (
	tlsLoadX509KeyPair = tls.LoadX509KeyPair
)

// readClientCert - helper function to read client certificate
// from pem formatted certPath and keyPath files
func readClientCert(certPath, keyPath string) ([]tls.Certificate, error) {
	if certPath == "" || keyPath == "" {
		return nil, nil
	}

	// load keypair
	cert, err := tlsLoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return nil, err
	}

	return []tls.Certificate{cert}, nil
}

// generateTLSConfig - helper function to generate a tls configuration based on
// config
func generateTLSConfig(c config) (*tls.Config, error) {
	certs, err := readClientCert(c.certPath, c.keyPath)
	if err != nil {
		return nil, err
	}
	return &tls.Config{
		InsecureSkipVerify: c.insecure,
		Certificates:       certs,
	}, nil
}
