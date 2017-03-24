package managerapi

import (
	"fmt"
	"bytes"
	"encoding/json"
	"net/http"

	dm"switchmanager/datamodel"
)

func (m *Manager) send(method string, url string, request interface{}, response interface{}) (error) {	
	_request, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error during json marshalling.")
		return err
	}

	var _response *http.Response

	switch method {
	case "POST":
		_response, err = m.client.Post(m.baseURL+url, "application/json", bytes.NewReader(_request))

	case "GET":
		_response, err = m.client.Get(m.baseURL+url)
	
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

// RegisterAgentPOST ...
func (m *Manager) RegisterAgentPOST(config dm.AgentConfig) {
	res := map[string]interface{}{}
	err := m.send("POST", "/do_register", config, &res)

	if err != nil {
		fmt.Println(err)			
	} else {
		fmt.Println("Agent registered")
	}
}

// UnregisterAgentPOST ...
func (m *Manager) UnregisterAgentPOST() {
	// TODO	
}
