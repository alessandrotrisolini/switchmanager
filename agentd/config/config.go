package config

import (
	"errors"

	"switchmanager/common"

	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
)

const agentCertPathConf string = "agent_cert"
const agentKeyPathConf string = "agent_key"
const caCertPathConf string = "ca_cert"
const managerDNSNameConf string = "manager_dns_name"
const managerPortConf string = "manager_port"
const agentDNSNameConf string = "agent_dns_name"
const agentPortConf string = "agent_port"
const interfacesConf string = "interfaces"
const openvSwitchConf string = "openvswitch"

// Config is the model representing the yaml config file for the agent
type Config struct {
	AgentCertPath  string
	AgentKeyPath   string
	CACertPath     string
	ManagerDNSName string
	ManagerPort    string
	AgentDNSName   string
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

	agentCertPath := configFile.Get(agentCertPathConf)
	if agentCertPath == nil {
		return config, errors.New("Agent certificate path is not present")
	}

	agentKeyPath := configFile.Get(agentKeyPathConf)
	if agentKeyPath == nil {
		return config, errors.New("Agent key path is not present")
	}

	caCertPath := configFile.Get(caCertPathConf)
	if caCertPath == nil {
		return config, errors.New("CA certificate path is not present")
	}

	managerIPAddress := configFile.Get(managerDNSNameConf)
	if managerDNSName == nil {
		return config, errors.New("Manager DNS name is not present")
	}

	managerPort := configFile.Get(managerPortConf)
	if managerPort == nil {
		return config, errors.New("Manager port is not present")
	}

	agentDNSName := configFile.Get(agentDNSNameConf)
	if agentIPAddress == nil {
		return config, errors.New("Agent DNS name is not present")
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

	config.AgentCertPath = to.String(agentCertPath)
	config.AgentKeyPath = to.String(agentKeyPath)
	config.CACertPath = to.String(caCertPath)
	config.ManagerDNSName = to.String(managerDNSName)
	config.ManagerPort = to.String(managerPort)
	config.AgentDNSName = to.String(agentDNSName)
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

	if checkValidConfig(config) {
		return config, nil
	}

	return config, errors.New("Configuration is not valid: check " +
		"the format of the IP address, port must be in a range " +
		"between 1024 and 65535 and strings must not contain " +
		"special characters (only letters and numbers)")
}

func checkValidConfig(c Config) bool {
	/*
		if !common.CheckIPAndPort(c.ManagerIPAddress, c.ManagerPort) {
			return false
		}
		if !common.CheckIPAndPort(c.AgentIPAddress, c.AgentPort) {
			return false
		}
	*/
	if !common.Sanitize(c.OpenvSwitch) {
		return false
	}
	for _, i := range c.Interfaces {
		if !common.Sanitize(i) {
			return false
		}
	}
	return true
}
