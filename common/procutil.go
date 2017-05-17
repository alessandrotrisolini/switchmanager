package common

import (
	"bufio"
	"os"
	"regexp"
	"strconv"
	"strings"
)

// CheckPID checks if a PID is well-formed (integer > 1).
// Return value of 0 means that the PID is invalid
func CheckPID(pid string, npid *int) bool {
	_pid, err := strconv.Atoi(pid)
	if err != nil || _pid < 2 {
		return false
	}
	*npid = _pid
	return true
}

// GetProcessState returns the state of a process with a given pid.
// The state of a process can be RUNNING, SLEEPING, or ZOMBIE
func GetProcessState(pid int) (string, error) {
	var state string
	var isStateLine = regexp.MustCompile("^State:.*").MatchString

	status, err := os.Open("/proc/" + strconv.Itoa(pid) + "/status")
	if err != nil {
		return "", err
	}
	defer status.Close()

	scanner := bufio.NewScanner(status)
	for scanner.Scan() {
		if isStateLine(scanner.Text()) {
			t := strings.Split(scanner.Text(), "\t")
			switch string(t[1][0]) {
			case "R":
				state = "RUNNING"
			case "S":
				state = "SLEEPING"
			case "Z":
				state = "ZOMBIE"
			}
		}
	}

	return state, nil
}
