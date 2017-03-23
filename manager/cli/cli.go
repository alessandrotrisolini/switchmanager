package cli

import (
	"bufio"
	"strings"
	"strconv"
	"fmt"

	"switchmanager/agentd/agentapi"
	"switchmanager/manager/managerutil"
	"github.com/fatih/color"
)

const shellString string = "manager$ "

/*
 *	Main cli loop
 */
func Start(a *agentapi.Agentd, c *color.Color, r *bufio.Reader) {
	for {
		args := NewLine(c, r)
		/*
		 *	Input validation and related actions
		 */
		if len(args) > 0 {
			DoCmd(a, args)
		}
	}
}

/*
 * 	Read new line
 */
func NewLine(c *color.Color, r *bufio.Reader) []string {
	c.Print(shellString)
	line, _ := r.ReadString('\n')
	line = managerutil.TrimSuffix(line, "\n")
	args := strings.Split(line, " ")
	return args
}

/*
 *	Parse and execute commands fed to managercli (when running)
 */
func DoCmd(a *agentapi.Agentd, args []string) {
	switch args[0] {
	case "run":
		a.InstantiateProcessPOST()
	case "kill":
		if len(args) > 1 {
			pid, err := strconv.Atoi(args[1])
			if err != nil || pid < 1 {
				fmt.Println("PID must be a positive number")
			} else {
				a.KillProcessPOST(pid)
			}
		} else {
			fmt.Println("PID is missing")
		}
	case "dump":
		a.DumpProcessesGET()
	case "":
	default:
		fmt.Println("Invalid command")
	}
}
