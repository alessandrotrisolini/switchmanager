package managerapi

import (
	"net/http"

	cmn "switchmanager/common"
	l "switchmanager/logging"
)

var log *l.Log

// ManagerAPI contains an HTTP endpoint and a URL, which
// is used as root for the REST calls
type ManagerAPI struct {
	client  *http.Client
	baseURL string
}

// NewManager returns a new instance of Manager
func NewManager(certPath string, keyPath string, caCertPath string) (*ManagerAPI, error) {
	client := &http.Client{}
	m := &ManagerAPI{client: client}
	err := cmn.SetupTLSClient(m.client, certPath, keyPath, caCertPath)
	if err != nil {
		return nil, err
	}
	log = l.GetLogger()
	return m, nil
}

// InitManager sets the base URL that is used as basis for
// REST calls
func (m *ManagerAPI) InitManager(baseURL string) {
	m.baseURL = baseURL
}
