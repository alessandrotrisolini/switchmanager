package manager

import (
	dm "switchmanager/datamodel"
)

// Manager contains data structures for managing
// the registration of agents
type Manager struct {
	// Agents are registered inside a map, where the key
	// is the IP address
	agents map[string]dm.AgentConfig
}

// NewManager returns a new Manager
func NewManager() *Manager {
	agents := make(map[string]dm.AgentConfig)
	m := &Manager{
		agents: agents,
	}
	return m
}

// RegisterAgent registers a new agent to the manager
func (m *Manager) RegisterAgent(conf dm.AgentConfig) {
	m.agents[conf.AgentDNSName] = conf
}

// GetRegistredAgents returns the registred agents
func (m *Manager) GetRegistredAgents() map[string]dm.AgentConfig {
	return m.agents
}

// GetRegistredAgent returns the configuration of an agent
func (m *Manager) GetRegistredAgent(dnsName string) dm.AgentConfig {
	return m.agents[dnsName]
}
