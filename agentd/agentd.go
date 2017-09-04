package main

import (
	"flag"
	"fmt"
	"os"

	as "switchmanager/agentd/agentserver"
	"switchmanager/agentd/config"
	"switchmanager/agentd/managerapi"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"

	nl "github.com/vishvananda/netlink"
)

var yamlPath string
var log *l.Log

func init() {
	flag.StringVar(&yamlPath, "config", "", "yaml configuration file path")
}

// Entry point of the agentd
func main() {
	if os.Geteuid() != 0 {
        fmt.Println("You must run agentd as super user")
		return
	}

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

	//Check if all the network interfaces exist
	for _, ifc := range yamlconf.Interfaces {
		l, err := nl.LinkByName(ifc)
		if err != nil {
			log.Error(ifc, "network interface does not exists.")
			return
		}
		err = nl.LinkSetUp(l)
		if err != nil {
			log.Error("Can not turn on network interface", ifc)
			return
		}
	}

	// Check if Open vSwitch exists
	_, err = nl.LinkByName(yamlconf.OpenvSwitch)
	if err != nil {
		log.Error(yamlconf.OpenvSwitch, "switch does not exists.")
		return
	}

	// Initialize client for manager REST API
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

	// Register agentd to the manager
	err = m.RegisterAgentPOST(conf)
	if err != nil {
		log.Error("Can not register agent:", err)
		return
	}

	// Configure and start agentd server
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
