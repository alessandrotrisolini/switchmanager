package agent

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"os/exec"
	"os"

	"restapi"
	"github.com/gorilla/mux"
)

type Agent struct {
	/*
	 *	Map that maintains a list of the processes that has
	 *	been instantiated by calling DoRun and that are 
	 *	running.
	 */
	processes	map[int]*os.Process

	/*
	 *	Handler for HTTP requests
	 */
	handler		http.Handler	
}

/*
 *	Global variable representing the agent with its data
 *	structures and handler.
 */
var agent Agent


func DoRun(w http.ResponseWriter, req *http.Request) {

	cmd := exec.Command("./foo")

	err := cmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	pid := cmd.Process.Pid

	fmt.Println("Process started - PID:", pid)

	agent.processes[pid] = cmd.Process

	json.NewEncoder(w).Encode(restapi.ProcessPid{Pid: pid})
}

/*
 *	Kills a process that has been instantiated by DoRun.
 *	The PID must be specified in the POST request.
 */
func DoKill(w http.ResponseWriter, req *http.Request) {
	var kill restapi.ProcessPid

	_ = json.NewDecoder(req.Body).Decode(&kill)
	
	if kill.Pid != 0 {
		fmt.Println("Trying to kill process with PID:", kill.Pid)
		for pid, process := range agent.processes {	
			if pid == kill.Pid {
				
				err := process.Kill()
				if err != nil {
					log.Fatal("Cannot kill process with PID:", pid)
				}

				_, err = process.Wait()
				if err != nil {
					log.Fatal("Cannot kill process with PID:", pid)
				}
		
				delete(agent.processes, pid)
				fmt.Println("Process killed!")
				json.NewEncoder(w).Encode(kill)
			}
		}
	}
}

/*
 *	Returns a list containing all the PID
 */
func DoProcDump(w http.ResponseWriter, req *http.Request) {
	json.NewEncoder(w).Encode(agent.processes)
}

/*
 *	Initialization of the agent
 */
func AgentInit() {
	processes := make(map[int]*os.Process)
	handler := mux.NewRouter()

	handler.HandleFunc("/do_run", DoRun).Methods("POST")
	handler.HandleFunc("/do_kill", DoKill).Methods("POST")
	handler.HandleFunc("/do_dump", DoProcDump).Methods("GET")

	agent = Agent {
		processes	: processes,
		handler 	: handler,
	}
}

/*
 *	Start agent HTTP server
 */
func AgentStart(port string) {
	log.Fatal(http.ListenAndServe(":" + port, agent.handler))
}
