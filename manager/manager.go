package main

import (
	"bufio"
	"flag"
	"os"
	"fmt"

	"switchmanager/agentd/agentapi"
	"switchmanager/manager/managerutil"
	"switchmanager/manager/cli"
	"switchmanager/manager/config"
	"github.com/fatih/color"
)

var agentIpAddress string
var agentPort string
var configFile string

/*
 *	Called by flag in order to parse command line parameters
 */
func init() {
	flag.StringVar(&agentIpAddress, "address", "", "agentd IP address")
	flag.StringVar(&configFile, "config", "", "managercli configuration file")
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

	config, err := config.GetConfig(configFile)
	if err != nil {
		fmt.Println("Error while reading configuration file:", err)
	}

	fmt.Println(config)

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
