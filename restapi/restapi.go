package restapi

import (
	"net/http"
)

type ProcessPid struct {
	Pid		int		`json:"pid, omitempty"`
}

type HostapdConfig struct {
	Interface		string	`json:"interface, omitempty"`
	ReathTimeout	int		`json:"reauth_timeout, omitempty"`
	OpenvSwitch		string	`json:"openvswitch, omitempty"`
}

type Agentd struct {
	client	*http.Client
	baseUrl	string
}

func NewAgentd() *Agentd {
	client := &http.Client{}
	d := &Agentd { client: client, }
	return d
}

func (a *Agentd) InitAgentd(baseUrl string) {
	a.baseUrl = baseUrl
}
