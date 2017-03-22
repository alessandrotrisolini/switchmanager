package main

import (
	"fmt"
	"os"
	"strconv"

	"agentd/agentutil"
)

var EMPTY_STRING string = ""

/*
 *	Enrty point of the agentd
 */
func main() {
	args := os.Args

	if !CheckArgsPresence(args) {
		PrintUsage()
		return
	}

	if !CheckPort(args[1]) {
		PrintUsage()
		return
	}

	fmt.Println("********************************")
	fmt.Println("*          CTRL AGENT          *")
	fmt.Println("********************************")

	agentutil.AgentInit()
	agentutil.AgentStart(args[1])
}

/*
 * Utilities with self-explaining function name
 */
func CheckArgsPresence(args []string) bool {
	return !(len(args) < 2)
}

func CheckPort(port string) bool {
	numeric_port, err := strconv.Atoi(port)
	return err == nil && port != EMPTY_STRING && numeric_port < 65536
}

func PrintUsage() {
	fmt.Println("Usage: agentd <port number>")
}
