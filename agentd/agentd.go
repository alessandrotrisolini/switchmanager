package main

import (
	"flag"
	"os"

	as "switchmanager/agentd/agentserver"
	"switchmanager/agentd/config"
	"switchmanager/agentd/managerapi"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
)

var yamlPath string
var log *l.Log

func init() {
	flag.StringVar(&yamlPath, "config", "", "yaml configuration file path")
}

// Entry point of the agentd
func main() {
	if !parseCommandLine() { return }

	logInit()

	yamlconf, err := config.GetConfig(yamlPath)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("********************************")
	log.Info("*         AGENT DAEMON         *")
	log.Info("********************************")
	log.Info("Configuration:", yamlconf)

	m := managerapi.NewManager()
	m.InitManager("http://127.0.0.1:5000")

	conf := dm.AgentConfig {
		AgentIPAddress: yamlconf.AgentIPAddress,
		AgentPort: yamlconf.AgentPort,
		Interfaces: yamlconf.Interfaces,
		OpenvSwitch: yamlconf.OpenvSwitch,
	}

	m.RegisterAgentPOST(conf)

	as.Init()
	as.Start(conf.AgentPort)
}

func logInit() {
	l.LogInit(os.Stdout)
	log = l.GetLogger()
}

func parseCommandLine() bool {
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		return false
	}
	return true
}
