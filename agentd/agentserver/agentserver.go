package agentserver

import (
	"encoding/json"
	"net/http"
	"os/exec"

	"switchmanager/agentd/agent"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
)

// Global variable representing the agent with its data
// structures and handler
var _agent *agent.Agent

var log *l.Log

func doRun(w http.ResponseWriter, req *http.Request) {
	cmd := exec.Command("./foo")

	err := cmd.Start()
	if err != nil {
		log.Error(err)
	}

	pid := cmd.Process.Pid

	log.Info("Process started - PID:", pid)

	_agent.AddProcess(pid, cmd.Process)

	json.NewEncoder(w).Encode(dm.ProcessPid{Pid: pid})
}

// Kills a process that has been instantiated by DoRun
// The PID must be specified in the POST request
func doKill(w http.ResponseWriter, req *http.Request) {
	var kill dm.ProcessPid

	_ = json.NewDecoder(req.Body).Decode(&kill)

	if kill.Pid != 0 {
		log.Info("Trying to kill process with PID:", kill.Pid)

		if _agent.CheckPid(kill.Pid) {
			err := _agent.DeleteProcess(kill.Pid)
			if err != nil {
				log.Error("Cannot stop process with PID:", kill.Pid)
			}
			log.Info("Process killed!")
			json.NewEncoder(w).Encode(kill)

		} else {
			log.Error("Process with PID", kill.Pid, "does not exist")
			kill.Pid = 0
			json.NewEncoder(w).Encode(kill)
		}
	}
}

// Returns a list containing all the PID
func doDump(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(_agent.DumpProcesses())
}

// Init initializes the agent server
func Init() {
	_agent = agent.NewAgent()

	_agent.SetHandleFunc("/do_run", doRun, "POST")
	_agent.SetHandleFunc("/do_kill", doKill, "POST")
	_agent.SetHandleFunc("/do_dump", doDump, "GET")

	log = l.GetLogger()
}

// Start starts the agent server
func Start(port string) {
	_agent.Start(port)
}
