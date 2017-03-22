package agentutil

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"

	"agentd/agent"
	"agentd/agentapi"
)

/*
 *	Global variable representing the agent with its data
 *	structures and handler.
 */
var _agent *agent.Agent


func DoRun(w http.ResponseWriter, req *http.Request) {

	cmd := exec.Command("./foo")

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	pid := cmd.Process.Pid

	fmt.Println("Process started - PID:", pid)

	_agent.AddProcess(pid, cmd.Process)

	json.NewEncoder(w).Encode(agentapi.ProcessPid{Pid: pid})
}

/*
 *	Kills a process that has been instantiated by DoRun.
 *	The PID must be specified in the POST request.
 */
func DoKill(w http.ResponseWriter, req *http.Request) {
	var kill agentapi.ProcessPid

	_ = json.NewDecoder(req.Body).Decode(&kill)
	
	if kill.Pid != 0 {
		fmt.Println("Trying to kill process with PID:", kill.Pid)

		if _agent.CheckPid(kill.Pid) {
			err := _agent.DeleteProcess(kill.Pid)
			if err != nil {
				log.Fatal("Cannot stop process with PID:", kill.Pid)
			}
			fmt.Println("Process killed!")
			json.NewEncoder(w).Encode(kill)

		}
	}
}

/*
 *	Returns a list containing all the PID
 */
func DoDump(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(_agent.DumpProcesses())
}

/*
 *	Initialization of the agent
 */
func AgentInit() {
	_agent = agent.NewAgent()

	_agent.SetHandleFunc("/do_run", DoRun, "POST")
	_agent.SetHandleFunc("/do_kill", DoKill, "POST")
	_agent.SetHandleFunc("/do_dump", DoDump, "GET")
}

/*
 *	Start agent HTTP server
 */
func AgentStart(port string) {
	_agent.Start(port)
}
