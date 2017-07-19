package cli

import (
	"bufio"
	"fmt"
	"strconv"
	"strings"
	"time"

	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
	"switchmanager/managercli/agentapi"
	c "switchmanager/managercli/config"
	ms "switchmanager/managercli/managerserver"

	clr "github.com/fatih/color"
)

const shellString string = "manager$ "

// Cli is the data type that maps the command line interface
type Cli struct {
	server *ms.ManagerServer
	conf   *c.Config
	color  *clr.Color
	reader *bufio.Reader
	log    *l.Log
}

// NewCli returns a reference to a new command line interface
func NewCli(c *clr.Color, r *bufio.Reader, mc *c.Config, server *ms.ManagerServer) *Cli {
	cli := &Cli{
		server: server,
		conf:   mc,
		color:  c,
		reader: r,
	}
	cli.log = l.GetLogger()
	return cli
}

// Start starts the main cli loop
func (cli *Cli) Start() {
	cli.startPolling()
	for {
		args := newLine(cli.color, cli.reader)
		// Input validation and related actions
		if len(args) > 0 {
			doCmd(args, cli)
		}
	}
}

func (cli *Cli) startPolling() {
	ticker := time.NewTicker(20 * time.Second).C
	go func() {
		for {
			<-ticker
			agents, _ := cli.server.RegistredAgents()
			for dnsName := range agents {
				a := createAgentAPI(cli, cli.server.GetAgentURL(dnsName))
				err := a.IsAliveGET()
				if err != nil {
					cli.server.DeleteAgent(dnsName)
				}
			}
		}
	}()
}

// Read new line
func newLine(c *clr.Color, r *bufio.Reader) []string {
	c.Print(shellString)
	return readLine(r)
}

func readLine(r *bufio.Reader) []string {
	line, _ := r.ReadString('\n')
	line = cmn.TrimSuffix(line, "\n")
	args := strings.Split(line, " ")
	return args
}

func run(args []string, cli *Cli) bool {
	if len(args) == 3 &&
		args[1] == "-hostname" {
		if cli.server.IsAgentRegistred(args[2]) {
			var hostapdConfig dm.HostapdConfig
			for i := 0; i < 3; i++ {
				switch i {
				case 0:
					fmt.Print(">> OpenvSwitch name: ")
					s := readLine(cli.reader)
					if cli.server.CheckOvsName(args[2], s[0]) {
						hostapdConfig.OpenvSwitch = s[0]
					} else {
						cli.log.Error("OvS does not exists")
						return true
					}
				case 1:
					fmt.Print(">> Interface name: ")
					s := readLine(cli.reader)
					if cli.server.CheckInterfaceName(args[2], s[0]) {
						hostapdConfig.Interface = s[0]
					} else {
						cli.log.Error("Interface does not exists")
						return true
					}
				case 2:
					fmt.Print(">> Reauthentication timeout: ")
					s := readLine(cli.reader)
					t, err := strconv.ParseUint(s[0], 10, 32)
					if err != nil {
						cli.log.Error("Reauthentication timeout must be a positive integer value")
						return true
					}
					hostapdConfig.ReauthTimeout = t
				}
			}

			hostapdConfig.RadiusAuthServer = "127.0.0.1"
			hostapdConfig.RadiusAcctServer = "127.0.0.1"
			hostapdConfig.RadiusSecret = "testing123"

			a := createAgentAPI(cli, cli.server.GetAgentURL(args[2]))
			a.InstantiateProcessPOST(hostapdConfig)
		} else {
			cli.log.Error("Agent @", args[2], "in not registred")
		}
		return true
	}
	return false
}

func kill(args []string, cli *Cli) bool {
	var pid int
	if len(args) == 5 &&
		args[1] == "-hostname" &&
		args[3] == "-pid" &&
		cmn.CheckPID(args[4], &pid) {
		if cli.server.IsAgentRegistred(args[2]) {
			a := createAgentAPI(cli, cli.server.GetAgentURL(args[2]))
			a.KillProcessDELETE(pid)
		} else {
			cli.log.Error("Agent @", args[2], "in not registred")
		}
		return true
	}
	return false
}

func dump(args []string, cli *Cli) bool {
	if len(args) == 3 &&
		args[1] == "-hostname" {
		if cli.server.IsAgentRegistred(args[2]) {
			a := createAgentAPI(cli, cli.server.GetAgentURL(args[2]))
			a.DumpProcessesGET()
		} else {
			cli.log.Error("Agent @", args[2], "in not registred")
		}
		return true
	}
	return false
}

func list(cli *Cli) {
	agents, err := cli.server.RegistredAgents()
	if err != nil {
		cli.log.Error(err)
	} else {
		if len(agents) == 0 {
			cli.log.Info("No agents have been registred")
		} else {
			fmt.Println(strings.Repeat("-", 48))
			fmt.Println("|               REGISTRED AGENTS               |")
			for k, v := range agents {
				fmt.Println(strings.Repeat("-", 48))
				fmt.Println("| HOSTNAME  :", k, strings.Repeat(" ", 48-(13+len(k)+4)), "|")
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
}

// Parse and execute commands fed to managercli (when running)
func doCmd(args []string, cli *Cli) {
	switch args[0] {
	case "run":
		if !run(args, cli) {
			cli.log.Error("Syntax: run -hostname <hostname>")
		}
	case "kill":
		if !kill(args, cli) {
			cli.log.Error("Syntax: kill -hostname <hostname> -pid <PID>")
		}
	case "dump":
		if !dump(args, cli) {
			cli.log.Error("Syntax: dump -hostname <hostname>")
		}
	case "list":
		list(cli)
	case "":
	default:
		cli.log.Error("Invalid command")
	}
}

func createAgentAPI(cli *Cli, IPAndPort string) *agentapi.AgentAPI {
	a := agentapi.NewAgentAPI(cli.conf.ManagerCertPath, cli.conf.ManagerKeyPath, cli.conf.CACertPath)
	a.InitAgentAPI("https://" + IPAndPort)
	return a
}
