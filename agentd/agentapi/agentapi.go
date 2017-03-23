package agentapi

import (
	"net/http"
)

const Run string = "/do_run"
const Kill string = "/do_kill"
const Dump string = "/do_dump"

type Agentd struct {
	client	*http.Client
	baseURL	string
}

func NewAgentd() *Agentd {
	client := &http.Client{}
	d := &Agentd { client: client, }
	return d
}

func (a *Agentd) InitAgentd(baseURL string) {
	a.baseURL = baseURL
}
