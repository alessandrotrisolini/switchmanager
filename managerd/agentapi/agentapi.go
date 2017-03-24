package agentapi

import (
	"net/http"

	l "switchmanager/logging"
)

const run string = "/do_run"
const kill string = "/do_kill"
const dump string = "/do_dump"

var log *l.Log

// Agentd ... 
type Agentd struct {
	client	*http.Client
	baseURL	string
}

// NewAgentd returns a new agentd
func NewAgentd() *Agentd {
	client := &http.Client{}
	d := &Agentd { client: client, }
	log = l.GetLogger()
	return d
}

// InitAgentd sets the base URL
func (a *Agentd) InitAgentd(baseURL string) {
	a.baseURL = baseURL
}
