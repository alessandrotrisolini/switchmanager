package config

import (
	"menteslibres.net/gosexy/yaml"
	"menteslibres.net/gosexy/to"
)

const AGENT_IP_ADDRESS 	string = "agent_ip_address"
const AGENT_PORT 		string = "agent_port"
const INTERFACE			string = "interface"
const REAUTH_TIMEOUT	string = "reauth_timeout"
const OPENVSWITCH		string = "openvswitch"

type Config struct {
	AgentIpAddress 	string
	AgentPort		string
	Interface 		string
	ReauthTimeout	string
	OpenvSwitch		string
}

func GetConfig(path string) (Config, error) {
	var config Config

	configfile, err := yaml.Open(path)
	if err != nil {
		return config, err
	}

	config.AgentIpAddress 	= to.String(configfile.Get(AGENT_IP_ADDRESS))
	config.AgentPort 		= to.String(configfile.Get(AGENT_PORT))
	//config.Interface	 	= to.String(configfile.Get(INTERFACE))
	//config.ReauthTimeout 	= to.String(configfile.Get(REAUTH_TIMEOUT))
	//config.OpenvSwitch	 	= to.String(configfile.Get(OPENVSWITCH))

	if CheckValidConfig(config) {
		return config, nil
	} else {
		return config, nil
	}
}

func CheckValidConfig(c Config) bool {
	return true
}
