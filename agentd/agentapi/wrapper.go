package agentapi

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"

	dm"switchmanager/datamodel"
)

func (a *Agentd) send(method string, url string, request interface{}, response interface{}) (error) {	
	_request, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error during json marshalling.")
		return err
	}

	var _response *http.Response

	switch method {
	case "POST":
		_response, err = a.client.Post(a.baseUrl+url, "application/json", bytes.NewReader(_request))

	case "GET":
		_response, err = a.client.Get(a.baseUrl+url)
	
	default:
		fmt.Println("Unknown method.")
	}

	if err != nil {
		return err
	}

	if _response.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned: %s", _response.Status)
	}
	
	if response != nil {
		err = json.NewDecoder(_response.Body).Decode(response)
	}
	
	return err
}

// InstantiateProcessPOST allows a manager to start a new process
func (a *Agentd) InstantiateProcessPOST() {
	req := map[string]interface{}{}
	var pid dm.ProcessPid
	err := a.send("POST", Run, req, &pid)

	if err != nil {
		fmt.Println(err)			
	} else {
		fmt.Println("Created process with PID", pid.Pid)
	}
}

// KillProcessPOST allows a manager to kill an active process
// that has been instantiated by calling InstantiateProcessPOST
func (a *Agentd) KillProcessPOST(pid int) {
	var req dm.ProcessPid
	var res dm.ProcessPid
	req.Pid = pid

	err := a.send("POST", Kill, req, &res)

	if err != nil {
		fmt.Println(err)
	} else if res.Pid != 0 {
		fmt.Println("Killed process with PID", pid)
	} else {
		fmt.Println("Process with PID", pid, "does not exist")
	}
}

// DumpProcessesGET allows a manager to dump all the active processes
// that have been instantiated by calling InstantiateProcessPOST
func (a *Agentd) DumpProcessesGET() {
	req := map[string]interface{}{}
	res := map[int]interface{}{}

	err := a.send("GET", Dump, req, &res)
	
	if err != nil {
		fmt.Println(err)
	} else {
		if len(res) == 0 {
			fmt.Println("No process is currently running")
		} else {
			fmt.Println("PID of instantiated processes:")
			for k := range res {
				fmt.Println(">> ", k)
			}
		}
	}
}
