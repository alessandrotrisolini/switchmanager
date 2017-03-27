package cli

import (
	"bufio"
	"strings"

	"github.com/fatih/color"
	cmn "switchmanager/common"
	"switchmanager/managerd/agentapi"
	ms "switchmanager/managerd/managerserver"
	l "switchmanager/logging"
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
				a := createAgentd(args[2])
				a.InstantiateProcessPOST()
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
				a := createAgentd(args[2])
				a.KillProcessPOST(pid)
		} else {
			log.Error("Syntax: kill -address <ip:port> -pid <PID>")
		}
	case "dump":
		//a.DumpProcessesGET()
	case "list":
		agents, err := ms.RegistredAgents()
		if err != nil {
			log.Error(err)
		} else {
			if len(agents) == 0 {
				log.Info("No agents have been registred")
			} else {
				for k, v := range agents {
					log.Info("---------------------------------------------")
					log.Info("IP ADDRESS:", k)
					log.Info("PORT      :", v.AgentPort)
					log.Info("OvS       :", v.OpenvSwitch)
					log.Info("INTERFACES:", v.Interfaces)	
				}
				log.Info("---------------------------------------------")
			}
		}
	case "":
	default:
		log.Error("Invalid command")
	}
}

func createAgentd(IPAndPort string) *agentapi.Agentd {
	a := agentapi.NewAgentd()
	a.InitAgentd("http://"+IPAndPort)
	return a
}
