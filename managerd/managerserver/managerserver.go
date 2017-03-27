package managerserver

import (
	"encoding/json"
	"errors"
	"net/http"

	m "switchmanager/managerd/manager"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
)

// Global variable representing the agent with its data
// structures and handler
var _manager *m.Manager

var log *l.Log

func doRegister(w http.ResponseWriter, req *http.Request) {
	var conf dm.AgentConfig

	_ = json.NewDecoder(req.Body).Decode(&conf)

	_manager.RegisterAgent(conf)

	log.Info("Registered agent with config:", conf)
	json.NewEncoder(w).Encode(dm.ProcessPid{Pid: 0})
}

// Init initializes the manager server
func Init() {
	_manager = m.NewManager()

	_manager.SetHandleFunc("/do_register", doRegister, "POST")

	log = l.GetLogger()
}

// Start starts the agent server
func Start() {
	log.Info("Starting manager server...")
	_manager.Start()
}

//RegistredAgents returns the list of the registered agents
func RegistredAgents() (map[string]dm.AgentConfig, error) {
	var agents map[string]dm.AgentConfig

	if _manager == nil {
		return agents, errors.New("Manager server has not been initialized")
	}

	return _manager.GetRegistredAgents(), nil
}
