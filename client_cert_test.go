package main

import (
	"crypto/tls"
	"errors"
	"testing"
)

func TestReadClientCertNoFilePaths(t *testing.T) {
	if _, err := readClientCert("certPath", ""); err != nil {
		t.Errorf("got an error that was not expected: %v\n", err)
	}
}

func TestReadClientCertFailedTLSLoadX509KeyPair(t *testing.T) {
	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, errors.New("failure")
	}
	defer func() { tlsLoadX509KeyPair = tls.LoadX509KeyPair }()

	if _, err := readClientCert("certPath", "keyPath"); err == nil {
		t.Errorf("expected an error from tlsLoadX509KeyPair\n")
	}
}

func TestReadClientCertSuccess(t *testing.T) {
	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}
	defer func() { tlsLoadX509KeyPair = tls.LoadX509KeyPair }()

	if _, err := readClientCert("certPath", "keyPath"); err != nil {
		t.Errorf("unexpected an error from readClientCert: %v\n", err)
	}
}

func TestGenerateTLSConfigError(t *testing.T) {
	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, errors.New("failure")
	}
	defer func() { tlsLoadX509KeyPair = tls.LoadX509KeyPair }()

	if _, err := generateTLSConfig(config{certPath: "certPath", keyPath: "keyPath"}); err == nil {
		t.Errorf("expected an error from generateTLSConfig\n")
	}
}

func TestGenerateTLSConfigSuccess(t *testing.T) {
	tlsLoadX509KeyPair = func(certFile, keyFile string) (tls.Certificate, error) {
		return tls.Certificate{}, nil
	}
	defer func() { tlsLoadX509KeyPair = tls.LoadX509KeyPair }()

	if _, err := generateTLSConfig(config{certPath: "certPath", keyPath: "keyPath"}); err != nil {
		t.Errorf("unexpected an error from generateTLSConfig: %v\n", err)
	}
}
