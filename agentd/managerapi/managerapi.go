package managerapi

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"

	l "switchmanager/logging"
)

var log *l.Log

// Manager contains an HTTP endpoint and a URL, which
// is used as root for the REST calls
type Manager struct {
	client  *http.Client
	baseURL string
}

// NewManager returns a new instance of Manager
func NewManager(cert string, key string, caCert string) (*Manager, error) {
	client := &http.Client{}
	m := &Manager{client: client}
	log = l.GetLogger()
	return m, nil
}

// InitManager sets the base URL that is used as basis for
// REST calls
func (m *Manager) InitManager(baseURL string) {
	m.baseURL = baseURL
}

func setupTLSClient(client *http.Client, certPath string, keyPath string, caCertPath string) error {
	cert, err := tls.LoadX509KeyPair(certPath, keyPath)
	if err != nil {
		return err
	}

	ca, err := ioutil.ReadFile(caCertPath)
	if err != nil {
		return err
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(ca)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      caCertPool,
	}
	tlsConfig.BuildNameToCertificate()

	transport := &http.Transport{TLSClientConfig: tlsConfig}
	client.Transport = transport
}

func setupTLSServer() {

}
