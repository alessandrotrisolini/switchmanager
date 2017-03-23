package managerutil

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

/*
 *	Collection of functions used by the CLI for input validation, modification
 *	and stuff like that.
 */

func CheckIpAndPort(ipAddress string, port string) bool {
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

func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}
