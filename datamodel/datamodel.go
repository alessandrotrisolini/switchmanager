package datamodel

type ProcessPid struct {
	Pid		int		`json:"pid, omitempty"`
}

type HostapdConfig struct {
	Interface		string	`json:"interface, omitempty"`
	ReauthTimeout	int		`json:"reauth_timeout, omitempty"`
	OpenvSwitch		string	`json:"openvswitch, omitempty"`
}
