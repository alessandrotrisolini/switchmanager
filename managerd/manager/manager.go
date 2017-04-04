package manager

import (
	"log"
	"net/http"

	cmn "switchmanager/common"
	dm "switchmanager/datamodel"

	"github.com/gorilla/mux"
)

// Port where the Manager server exposes the service
const port string = ":5000"

// Manager contains data structures for managing
// the registration of agents
type Manager struct {
	// Agents are registered inside a map, where the key
	// is the IP address
	agents map[string]dm.AgentConfig
	router *mux.Router
	server *http.Server

	certPath string
	keyPath  string
}

// NewManager returns a new Manager
func NewManager(certPath, keyPath, caCertPath string) (*Manager, error) {
	agents := make(map[string]dm.AgentConfig)
	router := mux.NewRouter()
	server := &http.Server{Addr: port}

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	server.Handler = router

	m := &Manager{
		agents:   agents,
		router:   router,
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
	}

	return m, nil
}

// SetHandleFunc adds an handler to the router
func (m *Manager) SetHandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	m.router.HandleFunc(url, f).Methods(method)
}

// Start starts the server
func (m *Manager) Start() {
	log.Fatal(m.server.ListenAndServeTLS(m.certPath, m.keyPath))
}

// RegisterAgent registers a new agent to the manager
func (m *Manager) RegisterAgent(conf dm.AgentConfig) {
	m.agents[conf.AgentIPAddress] = conf
}

// GetRegistredAgents returns the registred agents
func (m *Manager) GetRegistredAgents() map[string]dm.AgentConfig {
	return m.agents
}

// GetRegistredAgent returns the configuration of an agent
func (m *Manager) GetRegistredAgent(ip string) dm.AgentConfig {
	return m.agents[ip]
}
