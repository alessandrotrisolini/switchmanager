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
	if !parseCommandLine() {
		return
	}

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

	m, err := managerapi.NewManager(yamlconf.AgentCertPath, yamlconf.AgentKeyPath, yamlconf.CACertPath)
	if err != nil {
		log.Error("Can not initialize manager API:", err)
		return
	}
	m.InitManager("https://" + yamlconf.ManagerDNSName + ":" + yamlconf.ManagerPort)

	conf := dm.AgentConfig{
		AgentDNSName: yamlconf.AgentDNSName,
		AgentPort:    yamlconf.AgentPort,
		Interfaces:   yamlconf.Interfaces,
		OpenvSwitch:  yamlconf.OpenvSwitch,
	}

	err = m.RegisterAgentPOST(conf)
	if err != nil {
		log.Error("Can not register agent:", err)
		return
	}

	agentServer, err := as.NewAgentServer(yamlconf.AgentCertPath, yamlconf.AgentKeyPath, yamlconf.CACertPath)
	if err != nil {
		log.Error("Can not initalize TLS server:", err)
		return
	}

	agentServer.Start(conf.AgentPort)
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
