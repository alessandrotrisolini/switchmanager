package agentserver

import (
	"encoding/json"
	"net/http"
	"os"
	"os/exec"

	"github.com/gorilla/mux"

	a "switchmanager/agentd/agent"
	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
)

// AgentServer implements the server functionalities of the
// agent daemon
type AgentServer struct {
	agent  *a.Agent
	log    *l.Log
	router *mux.Router
	server *http.Server

	certPath string
	keyPath  string
}

// NewAgentServer initializes the agent server
func NewAgentServer(certPath string, keyPath string, caCertPath string) (*AgentServer, error) {
	agent := a.NewAgent()
	log := l.GetLogger()
	router := mux.NewRouter()
	server := &http.Server{Handler: router}

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	as := &AgentServer{
		agent:    agent,
		router:   router,
		server:   server,
		certPath: certPath,
		keyPath:  keyPath,
		log:      log,
	}

	router.Handle("/do_run", doRun(as)).Methods("POST")
	router.Handle("/do_kill", doKill(as)).Methods("POST")
	router.Handle("/do_dump", doDump(as)).Methods("GET")

	return as, nil
}

// Start starts the agent server
func (as *AgentServer) Start(port string) {
	as.server.Addr = ":" + port
	err := as.server.ListenAndServeTLS(as.certPath, as.keyPath)
	if err != nil {
		as.log.Error(err)
		os.Exit(1)
	}
}

func doRun(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		cmd := exec.Command("./foo")
		err := cmd.Start()
		if err != nil {
			as.log.Error(err)
		}
		var hostapdConfig dm.HostapdConfig
		_ = json.NewDecoder(req.Body).Decode(&hostapdConfig)
		as.log.Info(hostapdConfig)
		pid := cmd.Process.Pid
		as.log.Info("Process started - PID:", pid)
		as.agent.AddProcess(pid, cmd.Process)
		json.NewEncoder(w).Encode(dm.ProcessPid{Pid: pid})
	})
}

func doKill(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		var kill dm.ProcessPid
		_ = json.NewDecoder(req.Body).Decode(&kill)
		if kill.Pid != 0 {
			as.log.Info("Trying to kill process with PID:", kill.Pid)
			if as.agent.CheckPid(kill.Pid) {
				err := as.agent.DeleteProcess(kill.Pid)
				if err != nil {
					as.log.Error("Cannot stop process with PID:", kill.Pid)
				}
				as.log.Info("Process killed!")
				json.NewEncoder(w).Encode(kill)

			} else {
				as.log.Error("Process with PID", kill.Pid, "does not exist")
				kill.Pid = 0
				json.NewEncoder(w).Encode(kill)
			}
		}
	})
}

func doDump(as *AgentServer) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		json.NewEncoder(w).Encode(as.agent.DumpProcesses())
	})
}
