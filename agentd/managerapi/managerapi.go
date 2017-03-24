package managerapi

import (
	"net/http"
)

// Manager contains an HTTP endpoint and a URL, which
// is used as root for the REST calls
type Manager struct {
	client  *http.Client
	baseURL string
}

// NewManager returns a new instance of Manager
func NewManager() *Manager {
    client := &http.Client{}
	m := &Manager { client: client, }
	return m
}

// InitManager sets the base URL that is used as basis for
// REST calls
func (m *Manager) InitManager(baseURL string) {
	m.baseURL = baseURL
}
