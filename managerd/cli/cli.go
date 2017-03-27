package cli

import (
	"bufio"
	"strconv"
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
func Start(a *agentapi.Agentd, c *color.Color, r *bufio.Reader) {
	log = l.GetLogger()
	for {
		args := newLine(c, r)
		// Input validation and related actions
		if len(args) > 0 {
			doCmd(a, args)
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
func doCmd(a *agentapi.Agentd, args []string) {
	switch args[0] {
	case "run":
		a.InstantiateProcessPOST()
	case "kill":
		if len(args) > 1 {
			pid, err := strconv.Atoi(args[1])
			if err != nil || pid < 1 {
				log.Error("PID must be a positive number")
			} else {
				a.KillProcessPOST(pid)
			}
		} else {
			log.Error("PID is missing")
		}
	case "dump":
		a.DumpProcessesGET()
	case "list":
		agents, err := ms.RegistredAgents()
		if err != nil {
			log.Error(err)
		} else {
			log.Info(agents)
		}
	case "":
	default:
		log.Error("Invalid command")
	}
}
