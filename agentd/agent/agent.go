package agent

import (
	"net/http"
	"os"
	"log"

	"github.com/gorilla/mux"
)

// Agent contains the processes that have been instantiated
// by a manager and an handler for HTTP requests
type Agent struct {
	processes map[int]*os.Process
	router    *mux.Router	
}

// NewAgent returns a new agent
func NewAgent() *Agent {
	processes := make(map[int]*os.Process)
	router := mux.NewRouter()

	a := &Agent{
		processes: processes,
		router: router,
	}

	return a
}

// SetHandleFunc adds an handler to the router
func (a *Agent) SetHandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	a.router.HandleFunc(url, f).Methods(method) 
}

// Start starts the server
func (a *Agent) Start(port string) {
	log.Fatal(http.ListenAndServe(":" + port, a.router))
}

// AddProcess adds a process
func (a* Agent) AddProcess(pid int, p *os.Process) {
	a.processes[pid] = p
}

// DeleteProcess deletes a process with PID==pid, if it exists
func (a *Agent) DeleteProcess(pid int) error{
	err := a.processes[pid].Kill()
	if err != nil { 
		return err 
	}

	_, err = a.processes[pid].Wait()
	if err != nil {
		return err
	}

	delete(a.processes, pid)

	return nil
}

// CheckPid checks if a process with PID==pid is in the map
func (a* Agent) CheckPid(pid int) bool {
	return a.processes[pid] != nil
}

// DumpProcesses returns all the instantiated processes
func (a* Agent) DumpProcesses() map[int]*os.Process {
	return a.processes
}
