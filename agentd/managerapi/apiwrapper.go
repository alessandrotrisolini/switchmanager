package managerapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	dm "switchmanager/datamodel"
)

func (m *Manager) send(method string, url string, request interface{}, response interface{}) error {
	_request, err := json.Marshal(request)
	if err != nil {
		log.Error("Error during json marshalling")
		return err
	}

	var _response *http.Response

	switch method {
	case "POST":
		_response, err = m.client.Post(m.baseURL+url, "application/json", bytes.NewReader(_request))

	case "GET":
		_response, err = m.client.Get(m.baseURL + url)

	default:
		log.Error("Unknown method")
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
func (m *Manager) RegisterAgentPOST(config dm.AgentConfig) error {
	res := map[string]interface{}{}
	err := m.send("POST", "/agents", config, &res)

	if err != nil {
		log.Error(err)
		return err
	}

	log.Info("Agent registered")

	return nil
}

// UnregisterAgentPOST ...
func (m *Manager) UnregisterAgentPOST() {
	// TODO
}
