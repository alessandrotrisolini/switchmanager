package config

import (
	"errors"

	"menteslibres.net/gosexy/yaml"
)

const managerCertPathConf string = "manager_cert"
const managerKeyPathConf string = "manager_key"
const caCertPathConf string = "ca_cert"

// Config is the model representing the yaml config file for the manager
type Config struct {
	ManagerCertPath string
	ManagerKeyPath  string
	CACertPath      string
}

// GetConfig returns a Config struct when a path to a yaml file is passed
func GetConfig(path string) (Config, error) {
	var config Config

	configFile, err := yaml.Open(path)
	if err != nil {
		return config, err
	}

	managerCertPath := configFile.Get(managerCertPathConf)
	if managerCertPath == nil {
		return config, errors.New("Manager certificate path is not present")
	}

	managerKeyPath := configFile.Get(managerKeyPathConf)
	if agentKeyPath == nil {
		return config, errors.New("Manager key path is not present")
	}

	caCertPath := configFile.Get(caCertPathConf)
	if caCertPath == nil {
		return config, errors.New("CA certificate path is not present")
	}

	return config, nil
}
