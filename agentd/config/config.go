package config

import (
	"errors"

	"switchmanager/common"

	"github.com/spf13/viper"
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

	viper.SetConfigName("agent")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	agentCertPath := viper.Get(agentCertPathConf)
	if agentCertPath == nil {
		return config, errors.New("Agent certificate path is not present")
	}

	agentKeyPath := viper.Get(agentKeyPathConf)
	if agentKeyPath == nil {
		return config, errors.New("Agent key path is not present")
	}

	caCertPath := viper.Get(caCertPathConf)
	if caCertPath == nil {
		return config, errors.New("CA certificate path is not present")
	}

	managerDNSName := viper.Get(managerDNSNameConf)
	if managerDNSName == nil {
		return config, errors.New("Manager DNS name is not present")
	}

	managerPort := viper.Get(managerPortConf)
	if managerPort == nil {
		return config, errors.New("Manager port is not present")
	}

	agentDNSName := viper.Get(agentDNSNameConf)
	if agentDNSName == nil {
		return config, errors.New("Agent DNS name is not present")
	}

	agentPort := viper.Get(agentPortConf)
	if agentPort == nil {
		return config, errors.New("Agent port is not present")
	}

	interfaces := viper.Get(interfacesConf)
	if interfaces == nil {
		return config, errors.New("Interfaces are not present")
	}

	openvswitch := viper.Get(openvSwitchConf)
	if openvswitch == nil {
		return config, errors.New("Open vSwitch name is not present")
	}

	config.AgentCertPath = agentCertPath.(string)
	config.AgentKeyPath = agentKeyPath.(string)
	config.CACertPath = caCertPath.(string)
	config.ManagerDNSName = managerDNSName.(string)
	config.ManagerPort = managerPort.(string)
	config.AgentDNSName = agentDNSName.(string)
	config.AgentPort = agentPort.(string)

	ifc := common.List(interfaces)
	for _, i := range ifc {
		iMap := common.Map(i)
		name, found := iMap["name"]
		if !found {
			return config, errors.New("Can not find the name of an interface")
		}
		config.Interfaces = append(config.Interfaces, name.(string))
	}
	config.OpenvSwitch = openvswitch.(string)

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
