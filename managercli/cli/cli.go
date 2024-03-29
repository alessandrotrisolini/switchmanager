package cli

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
	"time"

	cmn "switchmanager/common"
	dm "switchmanager/datamodel"
	l "switchmanager/logging"
	"switchmanager/managercli/agentapi"
	c "switchmanager/managercli/config"
	m "switchmanager/managercli/manager"

	rl "github.com/chzyer/readline"
)

const shellString string = "\x1b[33m\x1b[1mmanager$\x1b[0m "

// Cli is the data type that maps the command line interface
type Cli struct {
	rlinstance *rl.Instance
	manager    *m.Manager
	conf       *c.Config
	reader     *bufio.Reader
	log        *l.Log
}

// NewCli returns a reference to a new command line interface
func NewCli(r *bufio.Reader, mc *c.Config, manager *m.Manager) *Cli {
	cli := &Cli{
		manager: manager,
		conf:   mc,
		reader: r,
	}
	cli.log = l.GetLogger()

	var completer = rl.NewPrefixCompleter(
		rl.PcItem("list"),
		rl.PcItem("run", rl.PcItem("-hostname")),
		rl.PcItem("kill", rl.PcItem("-hostname"), rl.PcItem("-pid")),
		rl.PcItem("dump", rl.PcItem("-hostname")),
	)

	l, _ := rl.NewEx(&rl.Config{
		Prompt:          shellString,
		HistoryFile:     "/tmp/readline.tmp",
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	})
	cli.rlinstance = l

	return cli
}

// Start starts the main cli loop
func (cli *Cli) Start() {
	cli.startPolling()
	defer cli.rlinstance.Close()

	for {
		args, err := cli.rlinstance.Readline()
		if err == rl.ErrInterrupt {
			if len(args) == 0 {
				break
			} else {
				continue
			}
		} else if err == io.EOF {
			break
		}
		// Input validation and related actions
		if len(args) > 0 {
			doCmd(strings.Split(args, " "), cli)
		}
	}
}

func (cli *Cli) startPolling() {
	ticker := time.NewTicker(20 * time.Second).C
	go func() {
		for {
			<-ticker
			agents := cli.manager.GetRegisteredAgents()
			for dnsName := range agents {
				a := createAgentAPI(cli, cli.manager.GetAgentURL(dnsName))
				err := a.IsAliveGET()
				if err != nil {
					cli.manager.DeleteAgent(dnsName)
				}
			}
		}
	}()
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
		if cli.manager.IsAgentRegistered(args[2]) {
			var hostapdConfig dm.HostapdConfig
			for i := 0; i < 5; i++ {
				switch i {
				case 0:
					fmt.Print(">> OpenvSwitch name: ")
					s := readLine(cli.reader)
					if cli.manager.CheckOvsName(args[2], s[0]) {
						hostapdConfig.OpenvSwitch = s[0]
					} else {
						cli.log.Error("OvS does not exists")
						return true
					}
				case 1:
					fmt.Print(">> Interface name: ")
					s := readLine(cli.reader)
					if cli.manager.CheckInterfaceName(args[2], s[0]) {
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
				case 3:
					fmt.Print(">> RADIUS server IP address: ")
					s := readLine(cli.reader)
					ip := net.ParseIP(s[0])
					if ip == nil {
						cli.log.Error("Wrong IP address format")
						return true
					}
					hostapdConfig.RadiusAuthServer = s[0]
					hostapdConfig.RadiusAcctServer = s[0]
				case 4:
					fmt.Print(">> RADIUS server shared secret: ")
					s := readLine(cli.reader)
					if len(s[0]) == 0 {
						cli.log.Error("RADIUS server shared secret can not be empty")
						return true
					}
					hostapdConfig.RadiusSecret = s[0]
				}
			}

			a := createAgentAPI(cli, cli.manager.GetAgentURL(args[2]))
			a.InstantiateProcessPOST(hostapdConfig)
		} else {
			cli.log.Error("Agent @", args[2], "in not registered")
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
		if cli.manager.IsAgentRegistered(args[2]) {
			a := createAgentAPI(cli, cli.manager.GetAgentURL(args[2]))
			a.KillProcessDELETE(pid)
		} else {
			cli.log.Error("Agent @", args[2], "in not registered")
		}
		return true
	}
	return false
}

func dump(args []string, cli *Cli) bool {
	if len(args) == 3 &&
		args[1] == "-hostname" {
		if cli.manager.IsAgentRegistered(args[2]) {
			a := createAgentAPI(cli, cli.manager.GetAgentURL(args[2]))
			a.DumpProcessesGET()
		} else {
			cli.log.Error("Agent @", args[2], "in not registered")
		}
		return true
	}
	return false
}

func list(cli *Cli) {
	agents := cli.manager.GetRegisteredAgents()
	if len(agents) == 0 {
		cli.log.Info("No agents have been registered")
	} else {
		fmt.Println(strings.Repeat("-", 48))
		fmt.Println("|               REGISTERED AGENTS              |")
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
