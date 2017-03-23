package agentapi

import (
	"net/http"
)

const RUN	string = "/do_run"
const KILL	string = "/do_kill"
const DUMP	string = "/do_dump"

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
