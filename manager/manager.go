package main

import (
	"bufio"
	"flag"
	"os"
	"fmt"

	"switchmanager/agentd/agentapi"
	"switchmanager/manager/managerutil"
	"switchmanager/manager/cli"
	"github.com/fatih/color"
)

var agentIPAddress string
var agentPort string
var configFile string

/*
 *	Called by flag in order to parse command line parameters
 */
func init() {
	flag.StringVar(&agentIPAddress, "address", "", "agentd IP address")
	flag.StringVar(&agentPort, "port", "", "agentd port")
}

/*
 *	Entry point for managercli
 */
func main() {
	if !parseCommandLine() {
	 	return 
	}

	/*
	 *	API initialization
	 */
	a := agentapi.NewAgentd()
	a.InitAgentd("http://" + agentIpAddress + ":" + agentPort)
	
	/*
	 *	CLI initialization
	 */
	c := color.New(color.FgYellow, color.Bold)
	r := bufio.NewReader(os.Stdin)

	cli.Start(a, c, r)
}

/*
 *	Parse managercli startup command line
 */
func parseCommandLine() bool {
	flag.Parse()
	if flag.NFlag() < 2 {
		flag.PrintDefaults()
		return false
	}
	return true
}
