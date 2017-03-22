package agent

import (
	"net/http"
	"os"
	"log"

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
	router		*mux.Router	
}

func NewAgent() *Agent {
	processes := make(map[int]*os.Process)
	router := mux.NewRouter()

	a := &Agent{
		processes: processes,
		router: router,
	}

	return a
}

func (a *Agent) SetHandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	a.router.HandleFunc(url, f).Methods(method) 
}

func (a *Agent) Start(port string) {
	log.Fatal(http.ListenAndServe(":" + port, a.router))
}

func (a* Agent) AddProcess(pid int, p *os.Process) {
	a.processes[pid] = p
}

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

func (a* Agent) CheckPid(pid int) bool {
	return a.processes[pid] != nil
}

func (a* Agent) DumpProcesses() map[int]*os.Process {
	return a.processes
}
