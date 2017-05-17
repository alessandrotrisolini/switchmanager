package agentapi

import (
	"net/http"

	cmn "switchmanager/common"
	l "switchmanager/logging"
)

var log *l.Log

// AgentAPI ...
type AgentAPI struct {
	client  *http.Client
	baseURL string
}

// NewAgentd returns a new agentd
func NewAgentd(certPath, certKeyPath, caCertPath string) *Agentd {
	client := &http.Client{}
	d := &AgentAPI{client: client}
	log = l.GetLogger()
	err := cmn.SetupTLSClient(d.client, certPath, certKeyPath, caCertPath)
	if err != nil {
		log.Error(err)
	}
	log = l.GetLogger()
	return d
}

// InitAgentd sets the base URL
func (a *AgentAPI) InitAgentd(baseURL string) {
	a.baseURL = baseURL
}
