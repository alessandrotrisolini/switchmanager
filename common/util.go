package common

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"
)

// Check if a string is composed only by alphabetic characters
// and numbers
var Sanitize = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// CheckIPAndPort checks if IP is well-formed and port is a
// non-standard port
func CheckIPAndPort(ipAddress string, port string) bool {
	ip := net.ParseIP(ipAddress)
	if ip == nil {
		fmt.Println("IP address is not valid")
		return false
	}
	
	_port, err := strconv.Atoi(port)
	if err != nil {
		fmt.Println(err)
		return false
	}

	if _port < 1024 || _port > 65535 {
		fmt.Println("Port is not in range <1024,65535>")
		return false
	}
	return true
}

// CheckArgsPresence checks that the length of an array of strings
// is at least 2
func CheckArgsPresence(args []string) bool {
	return !(len(args) < 2)
}

// CheckPort checks if a port is non-standard
func CheckPort(port string) bool {
	numericPort, err := strconv.Atoi(port)
	return err == nil && 
		port != "" && 
		numericPort > 1023 &&
		numericPort < 65536
}

//TrimSuffix deletes a suffix from a string and returns it
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
