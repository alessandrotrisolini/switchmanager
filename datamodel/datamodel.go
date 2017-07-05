package datamodel

// ProcessDescriptor wraps an integer representing the PID of a process
// that has been instantiated or killed and its state
type ProcessDescriptor struct {
	Pid   int    `json:"pid, omitempty"`
	State string `json:"state, omitempty"`
}

// HostapdConfig represents the configuration sent by the manager
// to a specific agent in order to start hostapd on a certain
// interface
type HostapdConfig struct {
	Interface        string `json:"interface, omitempty"`
	ReauthTimeout    uint64 `json:"reauth_timeout, omitempty"`
	OpenvSwitch      string `json:"openvswitch, omitempty"`
	RadiusAuthServer string `json:"radiusauthserver, omitempty"`
	RadiusAcctServer string `json:"radiusacctserver, omitempty"`
	RadiusSecret     string `json:"radiussecret, omitempty"`
}

// AgentConfig represents the configuration of an agent and it
// is sent to the manager during the registration phase
type AgentConfig struct {
	AgentDNSName string   `json:"dnsname, omitempty"`
	AgentPort    string   `json:"port, omitempty"`
	Interfaces   []string `json:"interfaces, omitempty"`
	OpenvSwitch  string   `json:"openvswitch, omitempty"`
}
