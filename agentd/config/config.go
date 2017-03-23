package config

import (
	"errors"

	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
)

const agentIPAddressConf string = "agent_ip_address"
const agentPortConf string = "agent_port"
const interfacesConf string = "interfaces"
const openvSwitchConf string = "openvswitch"

// Config is the model representing the yaml config file for the agent
type Config struct {
	AgentIPAddress string
	AgentPort      string
	Interfaces     []string
	OpenvSwitch    string
}

// GetConfig returns a Config struct when a path to a yaml file is passed
func GetConfig(path string) (Config, error) {
	var config Config

	configFile, err := yaml.Open(path)
	if err != nil {
		return config, err
	}

	agentIPAddress := configFile.Get(agentIPAddressConf)
	if agentIPAddress == nil {
		return config, errors.New("Agent IP address is not present")
	}

	agentPort := configFile.Get(agentPortConf)
	if agentPort == nil {
		return config, errors.New("Agent port is not present")
	}
	
	interfaces := configFile.Get(interfacesConf)
	if interfaces == nil {
		return config, errors.New("Interfaces are not present")
	}

	openvswitch := configFile.Get(openvSwitchConf)
	if openvswitch == nil {
		return config, errors.New("Open vSwitch name is not present")
	}

	config.AgentIPAddress = to.String(agentIPAddress)
	config.AgentPort = to.String(agentPort)
	ifc := to.List(interfaces)
	for _, i := range ifc {
		iMap := to.Map(i)
		name, found := iMap["name"]
		if !found {
			return config, errors.New("Can not find the name of an interface")
		}
		config.Interfaces = append(config.Interfaces, name.(string))
	}
	config.OpenvSwitch = to.String(openvswitch)

	return config, nil
}

func checkValidConfig(c Config) bool {
	return true
}
