package agent

import (
	"log"
	"net/http"
	"os"
	"syscall"

	cmn "switchmanager/common"

	"github.com/gorilla/mux"
)

// Agent contains the processes that have been instantiated
// by a manager and an handler for HTTP requests
type Agent struct {
	processes map[int]*os.Process
	router    *mux.Router
	server    *http.Server

	certPath string
	keyPath  string
}

// NewAgent returns a new agent
func NewAgent(certPath string, keyPath string, caCertPath string) (*Agent, error) {
	processes := make(map[int]*os.Process)
	router := mux.NewRouter()
	server := &http.Server{}

	err := cmn.SetupTLSServer(server, caCertPath)
	if err != nil {
		return nil, err
	}

	server.Handler = router
	server.Addr = ":8080"

	a := &Agent{
		processes: processes,
		router:    router,
		server:    server,
		certPath:  certPath,
		keyPath:   keyPath,
	}

	return a, nil
}

// SetHandleFunc adds an handler to the router
func (a *Agent) SetHandleFunc(url string, f func(http.ResponseWriter, *http.Request), method string) {
	a.router.HandleFunc(url, f).Methods(method)
}

// Start starts the server
func (a *Agent) Start(port string) {
	a.server.Addr = ":" + port
	log.Fatal(a.server.ListenAndServeTLS(a.certPath, a.keyPath))
}

// AddProcess adds a process
func (a *Agent) AddProcess(pid int, p *os.Process) {
	a.processes[pid] = p
}

// DeleteProcess deletes a process with PID==pid, if it exists
func (a *Agent) DeleteProcess(pid int) error {
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
func (a *Agent) CheckPid(pid int) bool {
	return a.processes[pid] != nil
}

// DumpProcesses returns all the instantiated processes
func (a *Agent) DumpProcesses() map[int]*os.Process {
	for pid, p := range a.processes {
		err := p.Signal(syscall.Signal(0))
		if err != nil {
			log.Println(err)
			delete(a.processes, pid)
		}
	}
	return a.processes
}
