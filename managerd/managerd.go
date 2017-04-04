package main

import (
	"bufio"
	"flag"
	"log/syslog"
	"os"

	l "switchmanager/logging"
	"switchmanager/managerd/cli"
	c "switchmanager/managerd/config"
	ms "switchmanager/managerd/managerserver"

	"github.com/fatih/color"
)

var yamlPath string
var log *l.Log

func init() {
	flag.StringVar(&yamlPath, "config", "", "yaml configuration file path")
}

// Entry point for managercli
func main() {
	// Parse and check command line
	if !parseCommandLine() {
		return
	}

	// Log initialization
	sl, _ := syslog.New(syslog.LOG_INFO, "")
	l.LogInit(sl)
	log = l.GetLogger()
	log.AddInfoOutput(os.Stdout)
	log.AddErrorOutput(os.Stdout)

	// Read configuration file
	conf, err := c.GetConfig(yamlPath)
	if err != nil {
		log.Error(err)
		return
	}

	// CLI initialization
	c := color.New(color.FgYellow, color.Bold)
	r := bufio.NewReader(os.Stdin)

	// Starting manager server
	err = ms.Init(conf.ManagerCertPath, conf.ManagerKeyPath, conf.CACertPath)
	if err != nil {
		log.Error("Manager server init failed:", err)
		return
	}
	go ms.Start()

	// Starting CLI
	cli.Start(c, r)
}

func parseCommandLine() bool {
	flag.Parse()
	if flag.NFlag() < 1 {
		flag.PrintDefaults()
		return false
	}
	return true
}
