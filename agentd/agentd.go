package main

import (
	"flag"
	"os"

	"switchmanager/agentd/config"
	l"switchmanager/logging"
)

var yamlPath string
var log *l.Log

func init() {
	flag.StringVar(&yamlPath, "config", "", "yaml configuration file path")
}

// Entry point of the agentd
func main() {

	if !parseCommandLine() { return }

	l.LogInit(os.Stdout)
	log = l.GetLogger()
	
	conf, err := config.GetConfig(yamlPath)
	if err != nil {
		log.Error(err)
	}

	log.Info("********************************")
	log.Info("*          CTRL AGENT          *")
	log.Info("********************************")
	log.Info(conf)

//	agentutil.AgentInit()
//	agentutil.AgentStart(args[1])
}

func parseCommandLine() bool {
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		return false
	}
	return true
}
