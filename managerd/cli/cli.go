package cli

import (
	"bufio"
	"fmt"
	"strings"

	cmn "switchmanager/common"
	l "switchmanager/logging"
	"switchmanager/managerd/agentapi"
	ms "switchmanager/managerd/managerserver"

	"github.com/fatih/color"
)

const shellString string = "manager$ "

var log *l.Log

// Start starts the main cli loop
func Start(c *color.Color, r *bufio.Reader) {
	log = l.GetLogger()
	for {
		args := newLine(c, r)
		// Input validation and related actions
		if len(args) > 0 {
			doCmd(args)
		}
	}
}

// Read new line
func newLine(c *color.Color, r *bufio.Reader) []string {
	c.Print(shellString)
	line, _ := r.ReadString('\n')
	line = cmn.TrimSuffix(line, "\n")
	args := strings.Split(line, " ")
	return args
}

// Parse and execute commands fed to managercli (when running)
func doCmd(args []string) {
	switch args[0] {
	case "run":
		if len(args) == 3 &&
			args[1] == "-address" &&
			cmn.CheckIPAndPort(args[2]) {
			if ms.IsAgentRegistred(args[2]) {
				a := createAgentd(args[2])
				a.InstantiateProcessPOST()
			} else {
				log.Error("Agent @", args[2], "in not registred")
			}
		} else {
			log.Error("Syntax: run -address <ip:port>")
		}
	case "kill":
		var pid int
		if len(args) == 5 &&
			args[1] == "-address" &&
			cmn.CheckIPAndPort(args[2]) &&
			args[3] == "-pid" &&
			cmn.CheckPID(args[4], &pid) {
			if ms.IsAgentRegistred(args[2]) {
				a := createAgentd(args[2])
				a.KillProcessPOST(pid)
			} else {
				log.Error("Agent @", args[2], "in not registred")
			}
		} else {
			log.Error("Syntax: kill -address <ip:port> -pid <PID>")
		}
	case "dump":
		if len(args) == 3 &&
			args[1] == "-address" &&
			cmn.CheckIPAndPort(args[2]) {
			if ms.IsAgentRegistred(args[2]) {
				a := createAgentd(args[2])
				a.DumpProcessesGET()
			} else {
				log.Error("Agent @", args[2], "in not registred")
			}
		} else {
			log.Error("Syntax: dump -address <ip:port>")
		}
	case "list":
		agents, err := ms.RegistredAgents()
		if err != nil {
			log.Error(err)
		} else {
			if len(agents) == 0 {
				log.Info("No agents have been registred")
			} else {
				fmt.Println(strings.Repeat("-", 48))
				fmt.Println("|               REGISTRED AGENTS               |")
				for k, v := range agents {
					fmt.Println(strings.Repeat("-", 48))
					fmt.Println("| IP ADDRESS:", k, strings.Repeat(" ", 48-(13+len(k)+4)), "|")
					fmt.Println("| PORT      :", v.AgentPort,
						strings.Repeat(" ", 48-(13+len(v.AgentPort)+4)), "|")
					fmt.Println("| OvS       :", v.OpenvSwitch,
						strings.Repeat(" ", 48-(13+len(v.OpenvSwitch)+4)), "|")
					fmt.Println("| INTERFACES:", v.Interfaces[0],
						strings.Repeat(" ", 48-(13+len(v.Interfaces[0])+4)), "|")
					for _, ifc := range v.Interfaces[1:] {
						fmt.Println("|", strings.Repeat(" ", 11), ifc,
							strings.Repeat(" ", 48-(13+len(ifc)+4)), "|")
					}
				}
				fmt.Println(strings.Repeat("-", 48))
			}
		}
	case "":
	default:
		log.Error("Invalid command")
	}
}

func createAgentd(IPAndPort string) *agentapi.Agentd {
	var certPath = "/home/alessandro/go/test/manager.pem"
	var keyPath = "/home/alessandro/go/test/manager.key"
	var caCertPath = "/home/alessandro/go/test/ca.pem"

	a := agentapi.NewAgentd(certPath, keyPath, caCertPath)
	a.InitAgentd("https://" + IPAndPort)
	return a
}
