package main

import (
	"bufio"
	"flag"
	"os"

	"switchmanager/agentd/agentapi"
	"switchmanager/managercli/managerutil"
	"switchmanager/managercli/cli"
	"github.com/fatih/color"
)

var agentIpAddress string
var agentPort string

/*
 *	Called by flag in order to parse command line parameters
 */
func init() {
	flag.StringVar(&agentIpAddress, "address", "", "agentd IP address")
	flag.StringVar(&agentPort, "port", "", "agentd port")
}

/*
 *	Entry point for managercli
 */
func main() {
	if !ParseCommandLine() || 
		!managerutil.CheckIpAndPort(agentIpAddress, agentPort) {
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
func ParseCommandLine() bool {
	flag.Parse()
	if flag.NFlag() < 2 {
		flag.PrintDefaults()
		return false
	}
	return true
}
