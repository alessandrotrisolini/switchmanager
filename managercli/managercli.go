package main

import (
	"bufio"
	"os"
	"strings"
	"fmt"
	"flag"
	"net"
	"strconv"

	"agentd/agentapi"
	"github.com/fatih/color"
)

var agentIpAddress string
var agentPort string

func init() {
	flag.StringVar(&agentIpAddress, "address", "", "agentd IP address")
	flag.StringVar(&agentPort, "port", "", "agentd port")
}

func CheckIpAndPort() bool {
	ip := net.ParseIP(agentIpAddress)
	if ip == nil {
		fmt.Println("IP address is not valid")
		return false
	}
	
	port, err := strconv.Atoi(agentPort)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if port < 1024 || port > 65535 {
		fmt.Println("Port is not in range <1024,65535>")
		return false
	}
	return true
}


func main() {
	
	flag.Parse()
	if !CheckIpAndPort() {
		return
	}
	
	a := agentapi.NewAgentd()
	a.InitAgentd("http://" + agentIpAddress + ":" + agentPort)
	
	c := color.New(color.FgYellow, color.Bold)

	for {
		reader := bufio.NewReader(os.Stdin)
		c.Print("manager$ ")
		line, _ := reader.ReadString('\n')

		line = TrimSuffix(line, "\n")
		args := strings.Split(line, " ")
		
		if len(args) > 0 {
			switch args[0] {
			case "run":
				agentapi.InstantiateProcessPOST(a)
			case "kill":
				if len(args) > 1 {
					pid, err := strconv.Atoi(args[1])
					if err != nil || pid < 1 {
						fmt.Println("PID must be a positive number")
					} else {
						agentapi.KillProcessPOST(a, pid)
					}
				} else {
					fmt.Println("PID is missing")
				}
			case "dump":
				agentapi.DumpProcessesGET(a)
			case "":
			default:
				fmt.Println("Unknown command")
			}
		}
	}
}

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
