package common

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

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

//TrimSuffix deletes a suffix from a string and returns it
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
