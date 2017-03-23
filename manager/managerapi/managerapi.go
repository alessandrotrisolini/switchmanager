package managerapi

import (
	"net/http"
)

// ManagerServer contains an HTTP endpoint and a URL, which
// is used as root for the REST calls
type ManagerServer struct {
	client  *http.Client
	baseURL string
}

// NewManager returns a new instance of Manager
func NewManager() *ManagerServer {
    client := &http.Client{}
    m = &Manager { client: client, }
	return m
}

// InitManager sets the base URL that is used as basis for
// REST calls
func (m *ManagerServer) InitManager(baseURL string) {
	m.baseURL = baseURL
}
