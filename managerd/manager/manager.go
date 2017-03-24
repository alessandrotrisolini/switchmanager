package manager

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
	dm "switchmanager/datamodel"
)

// Port where the Manager server exposes the service
const port string = "5000"

// Manager contains data structures for managing
// the registration of agents
type Manager struct {
	// Agents are registered inside a map, where the key
	// is the IP address
	agents map[string]dm.AgentConfig
	router *mux.Router
}

// NewManager returns a new Manager
func NewManager() *Manager {
	agents := make(map[string]dm.AgentConfig)
	router := mux.NewRouter()

	m := &Manager{
		agents: agents,
		router: router,
	}

	return m
}

// SetHandleFunc adds an handler to the router
func (m *Manager) SetHandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	m.router.HandleFunc(url, f).Methods(method)
}

// Start starts the server
func (m *Manager) Start() {
	log.Fatal(http.ListenAndServe(":"+port, m.router))
}

// RegisterAgent registers a new agent to the manager
func (m *Manager) RegisterAgent(conf dm.AgentConfig) {
	m.agents[conf.AgentIPAddress] = conf
}
