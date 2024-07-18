package common

import (
	"strconv"
	"strings"
)

func TrimArr(arr []string) (r []string) {
	for _, e := range arr {
		r = append(r, strings.Trim(e, " "))
	}
	return
}

// convert udp://127.0.0.1:53 or udp://127.0.0.1 to 127.0.0.1:53
func CanonicalAddr(addr string, port int) string {
	addrWithoutProto := strings.Join(strings.Split(addr, "://")[1:], "")
	if strings.Contains(addrWithoutProto, ":") {
		return addrWithoutProto
	} else {
		return addrWithoutProto + ":" + strconv.Itoa(port)
	}
}
