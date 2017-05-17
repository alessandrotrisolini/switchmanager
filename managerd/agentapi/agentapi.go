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

// NewAgentAPI returns a new agentd
func NewAgentAPI(certPath, certKeyPath, caCertPath string) *AgentAPI {
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

// InitAgentAPI sets the base URL
func (a *AgentAPI) InitAgentAPI(baseURL string) {
	a.baseURL = baseURL
}
