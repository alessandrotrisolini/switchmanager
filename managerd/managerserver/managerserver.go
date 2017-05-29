package managerserver

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"

	"github.com/gorilla/mux"

	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
	m "switchmanager/managerd/manager"
)

// Port where the Manager server exposes the service
const port string = ":5000"

//ManagerServer is the data type that models the server
type ManagerServer struct {
	manager  *m.Manager
	router   *mux.Router
	server   *http.Server
	certPath string
	keyPath  string
	log      *l.Log
}

func doRegister(ms *ManagerServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var conf dm.AgentConfig
		_ = json.NewDecoder(req.Body).Decode(&conf)
		ms.manager.RegisterAgent(conf)
		ms.log.Trace("Registered agent with config:", conf)
		//json.NewEncoder(w).Encode(dm.ProcessDescriptor{Pid: 0})
		w.WriteHeader(http.StatusOK)
	})
}

// NewManagerServer initializes the manager server
func NewManagerServer(certPath, keyPath, caCertPath string) (*ManagerServer, error) {
	manager := m.NewManager()
	log := l.GetLogger()
	router := mux.NewRouter()
	server := &http.Server{Addr: port}
	server.Handler = router

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	ms := &ManagerServer{
		manager:  manager,
		router:   router,
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
		log:      log,
	}

	router.Handle("/agents", doRegister(ms)).Methods("POST")

	return ms, nil
}

// Start starts the server
func (ms *ManagerServer) Start() {
	ms.log.Trace("Starting manager server...")
	err := ms.server.ListenAndServeTLS(ms.certPath, ms.keyPath)
	if err != nil {
		ms.log.Error(err)
		os.Exit(1)
	}
}

// RegistredAgents returns the list of the registered agents
func (ms *ManagerServer) RegistredAgents() (map[string]dm.AgentConfig, error) {
	var agents map[string]dm.AgentConfig
	if ms.manager == nil {
		return agents, errors.New("Manager server has not been initialized")
	}
	return ms.manager.GetRegistredAgents(), nil
}

// GetAgentURL returns the complete URL where the agent exposes its API
func (ms *ManagerServer) GetAgentURL(dnsName string) string {
	a := ms.manager.GetRegistredAgent(dnsName)
	return dnsName + ":" + a.AgentPort
}

// IsAgentRegistred checks if an agent has been registred
func (ms *ManagerServer) IsAgentRegistred(dnsName string) bool {
	a := ms.manager.GetRegistredAgent(dnsName)
	return a.AgentDNSName == dnsName
}

// CheckOvsName checks if a specific agent includes a switch named with ovsName
func (ms *ManagerServer) CheckOvsName(dnsName string, ovsName string) bool {
	a := ms.manager.GetRegistredAgent(dnsName)
	return a.AgentDNSName == dnsName &&
		a.OpenvSwitch == ovsName
}

// CheckInterfaceName checks if a specific agent includes an interface named
// with ifcName
func (ms *ManagerServer) CheckInterfaceName(dnsName string, ifcName string) bool {
	a := ms.manager.GetRegistredAgent(dnsName)
	return a.AgentDNSName == dnsName &&
		cmn.Contains(a.Interfaces, ifcName)
}
