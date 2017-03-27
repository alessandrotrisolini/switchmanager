package main

import (
	"bufio"
	"log/syslog"
	"os"

	l "switchmanager/logging"
	"switchmanager/managerd/cli"
	ms "switchmanager/managerd/managerserver"
	"github.com/fatih/color"
)

// Entry point for managercli
func main() {
	// CLI initialization
	c := color.New(color.FgYellow, color.Bold)
	r := bufio.NewReader(os.Stdin)
	
	// Log initialization
	sl, _ := syslog.New(syslog.LOG_INFO, "")
	l.LogInit(sl)
	l.GetLogger().AddInfoOutput(os.Stdout)
	l.GetLogger().AddErrorOutput(os.Stdout)

	// Starting manager server
	ms.Init()
	go ms.Start()

	// Starting CLI
	cli.Start(c, r)
}
