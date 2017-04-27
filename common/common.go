package common

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Check if a string is composed only by alphabetic characters
// and numbers
var Sanitize = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// CheckIPAndPort checks if IP is well-formed and port is a
// non-standard port
func CheckIPAndPort(s ...string) bool {
	var ipAddress, port string
	if len(s) == 1 {
		ipAddress, port = ParseIPAndPort(s[0])
	} else {
		ipAddress = s[0]
		port = s[1]
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}

	_port, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if _port < 1024 || _port > 65535 {
		return false
	}
	return true
}

// ParseIPAndPort returns IP and Port from a single string.
func ParseIPAndPort(s string) (string, string) {
	ipport := strings.Split(s, ":")
	if len(ipport) == 2 {
		return ipport[0], ipport[1]
	}

	return "", ""
}

// CheckArgsPresence checks that the length of an array of strings
// is at least 2
func CheckArgsPresence(args []string) bool {
	return !(len(args) < 2)
}

// CheckPort checks if a port is non-standard
func CheckPort(port string) bool {
	numericPort, err := strconv.Atoi(port)
	return err == nil &&
		port != "" &&
		numericPort > 1023 &&
		numericPort < 65536
}

// CheckPID checks if a PID is well-formed (integer > 1).
// Return value of 0 means that the PID is invalid
func CheckPID(pid string, npid *int) bool {
	_pid, err := strconv.Atoi(pid)
	if err != nil || _pid < 2 {
		return false
	}
	*npid = _pid
	return true
}

//Contains checks if an element is contained in a slice
func Contains(slice interface{}, elem interface{}) bool {
	arr := reflect.ValueOf(slice)
	if arr.Kind() == reflect.Slice {
		for i := 0; i < arr.Len(); i++ {
			if arr.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

//TrimSuffix deletes a suffix from a string and returns it
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

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
