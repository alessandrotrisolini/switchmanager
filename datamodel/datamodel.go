package datamodel

// ProcessPid wraps an integer representing the PID of a process
// that has been instantiated or killed
type ProcessPid struct {
	Pid int `json:"pid, omitempty"`
}

// HostapdConfig represents the configuration sent by the manager
// to a specific agent in order to start hostapd on a certain
// interface
type HostapdConfig struct {
	Interface     string `json:"interface, omitempty"`
	ReauthTimeout uint64 `json:"reauth_timeout, omitempty"`
	OpenvSwitch   string `json:"openvswitch, omitempty"`
}

// AgentConfig represents the configuration of an agent and it
// is sent to the manager during the registration phase
type AgentConfig struct {
	AgentIPAddress string   `json:"ipaddress, omitempty"`
	AgentPort      string   `json:"port, omitempty"`
	Interfaces     []string `json:"interfaces, omitempty"`
	OpenvSwitch    string   `json:"openvswitch, omitempty"`
}
