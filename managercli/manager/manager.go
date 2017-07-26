package manager

import (
	dm "switchmanager/datamodel"
)

// Manager contains data structures for managing
// the registration of agents
type Manager struct {
	// Agents are registered inside a map, where the key
	// is the agent DNS name
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

// DeleteAgent deletes an agent given its DNS name
func (m *Manager) DeleteAgent(dnsName string) {
	delete(m.agents, dnsName)
}

// GetRegisteredAgents returns the registered agents
func (m *Manager) GetRegisteredAgents() map[string]dm.AgentConfig {
	return m.agents
}

// GetRegisteredAgent returns the configuration of an agent
func (m *Manager) GetRegisteredAgent(dnsName string) dm.AgentConfig {
	return m.agents[dnsName]
}
