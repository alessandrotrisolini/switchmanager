package config

import (
	"menteslibres.net/gosexy/to"
	"menteslibres.net/gosexy/yaml"
)

const agentIPAddressConf string = "agent_ip_address"
const agentPortConf string = "agent_port"
const interfaceConf string = "interface"
const reauthTimeoutConf string = "reauth_timeout"
const openvSwitchConfstring = "openvswitch"

// Config is the model representing the yaml config file for the agent
type Config struct {
	AgentIPAddress string
	AgentPort      string
	Interface      []string
	ReauthTimeout  string
	OpenvSwitch    string
}

// GetConfig returns a Config struct when a path to a yaml file is passed
func GetConfig(path string) (Config, error) {
	var config Config

	configfile, err := yaml.Open(path)
	if err != nil {
		return config, err
	}

	config.AgentIpAddress = to.String(configfile.Get(agentIPAddressConf))
	config.AgentPort = to.String(configfile.Get(agentPortConf))
	//config.Interface	 	= to.String(configfile.Get(INTERFACE))
	//config.ReauthTimeout 	= to.String(configfile.Get(REAUTH_TIMEOUT))
	//config.OpenvSwitch	= to.String(configfile.Get(OPENVSWITCH))
}

func checkValidConfig(c Config) bool {
	return true
}
