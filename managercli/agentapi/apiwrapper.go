package agentapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	dm "switchmanager/datamodel"
)

func (a *AgentAPI) send(method string, url string, request interface{}, response interface{}) error {
	b, err := json.Marshal(request)
	if err != nil {
		log.Error("Error during json marshalling")
		return err
	}

	var _response *http.Response

	switch method {
	case "POST":
		_response, err = a.client.Post(a.baseURL+url, "application/json", bytes.NewReader(b))

	case "GET":
		_response, err = a.client.Get(a.baseURL + url)

	case "DELETE":
		r, err := http.NewRequest("DELETE", a.baseURL+url, bytes.NewReader(b))
		if err != nil {
			log.Error(err)
		}
		_response, err = a.client.Do(r)

	default:
		log.Error("Unknown method")
	}

	if err != nil {
		return err
	}

	defer _response.Body.Close()

	if _response.StatusCode != http.StatusOK {
		return fmt.Errorf("Server returned: %s", _response.Status)
	}

	if response != nil {
		err = json.NewDecoder(_response.Body).Decode(response)
	}

	return err
}

// InstantiateProcessPOST allows a manager to start a new process
func (a *AgentAPI) InstantiateProcessPOST(hostapdConfig dm.HostapdConfig) {
	var pid dm.ProcessDescriptor
	err := a.send("POST", "/processes", hostapdConfig, &pid)
	if err != nil {
		log.Error(err)
	} else {
		log.Info("Created process with PID", pid.Pid)
	}
}

// KillProcessDELETE allows a manager to kill an active process
// that has been instantiated by calling InstantiateProcessPOST
func (a *AgentAPI) KillProcessDELETE(pid int) {
	req := map[string]interface{}{}
	err := a.send("DELETE", "/processes/"+strconv.Itoa(pid), req, nil)
	if err != nil {
		log.Error(err)
	} else if pid != 0 {
		log.Info("Killed process with PID", pid)
	}
}

// DumpProcessesGET allows a manager to dump all the active processes
// that have been instantiated by calling InstantiateProcessPOST
func (a *AgentAPI) DumpProcessesGET() {
	var res map[int]dm.ProcessDescriptor
	err := a.send("GET", "/processes", nil, &res)
	if err != nil {
		log.Error(err)
	} else {
		if len(res) == 0 {
			log.Info("No process is currently running")
		} else {
			fmt.Println("PID of instantiated processes @", strings.Split(a.baseURL, "/")[2], ":")
			for k, p := range res {
				fmt.Println(">> PID:", k, "- State:", p.State)
			}
		}
	}
}

// IsAliveGET allows a manager to check if an agent is available
func (a *AgentAPI) IsAliveGET() error {
	err := a.send("GET", "/alive", nil, nil)
	return err
}
