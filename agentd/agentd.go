package main

import (
	"flag"
	"os"

	as "switchmanager/agentd/agentserver"
	"switchmanager/agentd/config"
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

	conf, err := config.GetConfig(yamlPath)
	if err != nil {
		log.Error(err)
		return
	}

	log.Info("********************************")
	log.Info("*         AGENT DAEMON         *")
	log.Info("********************************")
	log.Info("Configuration:", conf)

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
