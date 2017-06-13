package config

import (
	"errors"

	"github.com/spf13/viper"
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

	viper.SetConfigName("manager")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		return config, err
	}

	managerCertPath := viper.Get(managerCertPathConf)
	if managerCertPath == nil {
		return config, errors.New("Manager certificate path is not present")
	}

	managerKeyPath := viper.Get(managerKeyPathConf)
	if managerKeyPath == nil {
		return config, errors.New("Manager key path is not present")
	}

	caCertPath := viper.Get(caCertPathConf)
	if caCertPath == nil {
		return config, errors.New("CA certificate path is not present")
	}

	config.ManagerCertPath = managerCertPath.(string)
	config.ManagerKeyPath = managerKeyPath.(string)
	config.CACertPath = caCertPath.(string)

	return config, nil
}
