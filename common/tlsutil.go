package common

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net/http"
)

func commonTLSConfig(caCertPath string) (*tls.Config, error) {

	ca, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return nil, err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(ca) {
		return nil, errors.New("Can not load CA certificate")
	}

	tlsConfig := &tls.Config{
		MinVersion:             tls.VersionTLS12,
		SessionTicketsDisabled: true,
		RootCAs:                caCertPool,
	}

	/*
		tlsConfig.CipherSuites = []uint16{tls.TLS_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_256_CBC_SHA,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256}
	*/

	return tlsConfig, nil
}

//SetupTLSClient initializes a TLS client
func SetupTLSClient(client *http.Client, certPath string, keyPath string, caCertPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	tlsConfig, err := commonTLSConfig(caCertPath)
	if err != nil {
		return err
	}

	tlsConfig.Certificates = []tls.Certificate{cert}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client.Transport = transport

	return nil
}

// SetupTLSServer initializes a TLS server which requires client authN
func SetupTLSServer(server *http.Server, caCertPath string) error {
	ca, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(ca) {
		return errors.New("Can not load CA certificate")
	}

	tlsConfig, err := commonTLSConfig(caCertPath)
	if err != nil {
		return err
	}

	tlsConfig.ClientCAs = tlsConfig.RootCAs
	tlsConfig.ClientAuth = tls.RequireAndVerifyClientCert

	tlsConfig.BuildNameToCertificate()

	server.TLSConfig = tlsConfig

	return nil
}
