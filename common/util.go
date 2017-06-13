package common

import (
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

const (
	digits     = "0123456789"
	uintbuflen = 20
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

//List converts an interface to a list - ref. menteslibres.net/gosexy/to
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

//Map converts an interface to a map - ref. menteslibres.net/gosexy/to
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

//String converts an interface to a string - ref. menteslibres.net/gosexy/to
func String(val interface{}) string {
	var buf []byte

	if val == nil {
		return ""
	}

	switch t := val.(type) {

	case int:
		buf = int64ToBytes(int64(t))
	case int8:
		buf = int64ToBytes(int64(t))
	case int16:
		buf = int64ToBytes(int64(t))
	case int32:
		buf = int64ToBytes(int64(t))
	case int64:
		buf = int64ToBytes(int64(t))

	case uint:
		buf = uint64ToBytes(uint64(t))
	case uint8:
		buf = uint64ToBytes(uint64(t))
	case uint16:
		buf = uint64ToBytes(uint64(t))
	case uint32:
		buf = uint64ToBytes(uint64(t))
	case uint64:
		buf = uint64ToBytes(uint64(t))

	case float32:
		buf = float32ToBytes(t)
	case float64:
		buf = float64ToBytes(t)

	case complex128:
		buf = complex128ToBytes(t)
	case complex64:
		buf = complex128ToBytes(complex128(t))

	case bool:
		if val.(bool) == true {
			return "true"
		}

		return "false"

	case string:
		return t

	case []byte:
		return string(t)

	default:
		return fmt.Sprintf("%v", val)
	}

	return string(buf)
}

func uint64ToBytes(v uint64) []byte {
	buf := make([]byte, uintbuflen)

	i := len(buf)

	for v >= 10 {
		i--
		buf[i] = digits[v%10]
		v = v / 10
	}

	i--
	buf[i] = digits[v%10]

	return buf[i:]
}

func int64ToBytes(v int64) []byte {
	negative := false

	if v < 0 {
		negative = true
		v = -v
	}

	uv := uint64(v)

	buf := uint64ToBytes(uv)

	if negative {
		buf2 := []byte{'-'}
		buf2 = append(buf2, buf...)
		return buf2
	}

	return buf
}

func float32ToBytes(v float32) []byte {
	slice := strconv.AppendFloat(nil, float64(v), 'g', -1, 32)
	return slice
}

func float64ToBytes(v float64) []byte {
	slice := strconv.AppendFloat(nil, v, 'g', -1, 64)
	return slice
}

func complex128ToBytes(v complex128) []byte {
	buf := []byte{'('}

	r := strconv.AppendFloat(buf, real(v), 'g', -1, 64)

	im := imag(v)
	if im >= 0 {
		buf = append(r, '+')
	} else {
		buf = r
	}

	i := strconv.AppendFloat(buf, im, 'g', -1, 64)

	buf = append(i, []byte{'i', ')'}...)

	return buf
}
