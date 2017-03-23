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

const SHELL_STRING string = "manager$ "

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
	c.Print(SHELL_STRING)
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
		fmt.Println("Invalid command")
	}
}
