package common

import (
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

// Check if a string is composed only by alphabetic characters
// and numbers
var Sanitize = regexp.MustCompile(`^[a-zA-Z0-9]+$`).MatchString

// CheckIPAndPort checks if IP is well-formed and port is a
// non-standard port
func CheckIPAndPort(s ...string) bool {
	var ipAddress, port string
	if len(s) == 1 {
		ipAddress, port = ParseIPAndPort(s[0])
	} else {
		ipAddress = s[0]
		port = s[1]
	}

	ip := net.ParseIP(ipAddress)
	if ip == nil {
		return false
	}

	_port, err := strconv.Atoi(port)
	if err != nil {
		return false
	}

	if _port < 1024 || _port > 65535 {
		return false
	}
	return true
}

// ParseIPAndPort returns IP and Port from a single string.
func ParseIPAndPort(s string) (string, string) {
	ipport := strings.Split(s, ":")
	if len(ipport) == 2 {
		return ipport[0], ipport[1]
	}

	return "", ""
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

//Contains checks if an element is contained in a slice
func Contains(slice interface{}, elem interface{}) bool {
	arr := reflect.ValueOf(slice)
	if arr.Kind() == reflect.Slice {
		for i := 0; i < arr.Len(); i++ {
			if arr.Index(i).Interface() == elem {
				return true
			}
		}
	}
	return false
}

//TrimSuffix deletes a suffix from a string and returns it
func TrimSuffix(s, suffix string) string {
	if strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

//List converts an interface to a list
func List(val interface{}) []interface{} {
	list := []interface{}{}

	if val == nil {
		return list
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Slice:
		vval := reflect.ValueOf(val)

		size := vval.Len()
		list := make([]interface{}, size)
		vlist := reflect.ValueOf(list)

		for i := 0; i < size; i++ {
			vlist.Index(i).Set(vval.Index(i))
		}

		return list
	}

	return list
}

//Map converts an interface to a map
func Map(val interface{}) map[string]interface{} {
	list := map[string]interface{}{}

	if val == nil {
		return list
	}

	switch reflect.TypeOf(val).Kind() {
	case reflect.Map:
		vval := reflect.ValueOf(val)
		vlist := reflect.ValueOf(list)

		for _, vkey := range vval.MapKeys() {
			key := String(vkey.Interface())
			vlist.SetMapIndex(reflect.ValueOf(key), vval.MapIndex(vkey))
		}

		return list
	}

	return list
}
