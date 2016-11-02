package main

import (
	"crypto/tls"
	"net"
	"net/url"
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
	uri, err := url.ParseRequestURI(c.url)
	if err != nil {
		return nil, err
	}
	tlsConfig := &tls.Config{
		InsecureSkipVerify: c.insecure,
		Certificates:       certs,
	}
	if uri.Scheme == "https" {
		host, _, err := net.SplitHostPort(uri.Host)
		if err != nil {
			host = uri.Host
		}
		tlsConfig.ServerName = host
	}
	return tlsConfig, nil
}
