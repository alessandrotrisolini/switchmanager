package agent

import (
	"os"
)

// Process is a wrapper for os.Process which includes the actual state of the process
type Process struct {
	Process *os.Process
	State   string
}

// Agent contains the processes that have been instantiated
// by a manager and an handler for HTTP requests
type Agent struct {
	processes map[int]*Process
}

// NewAgent returns a new agent
func NewAgent() *Agent {
	processes := make(map[int]*Process)
	a := &Agent{
		processes: processes,
	}
	return a
}

// AddProcess adds a process
func (a *Agent) AddProcess(pid int, p *Process) {
	a.processes[pid] = p
}

// DeleteProcess deletes a process with PID==pid, if it exists
func (a *Agent) DeleteProcess(pid int) error {
	err := a.processes[pid].Process.Kill()
	if err != nil {
		return err
	}
	_, err = a.processes[pid].Process.Wait()
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
func (a *Agent) DumpProcesses() map[int]*Process {
	return a.processes
}
