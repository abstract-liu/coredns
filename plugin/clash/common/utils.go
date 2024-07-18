package common

import (
	"github.com/miekg/dns"
	"net"
	"net/netip"
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

func IpToAddr(slice net.IP) netip.Addr {
	ip := slice
	if len(ip) != 4 {
		if ip = slice.To4(); ip == nil {
			ip = slice
		}
	}

	if addr, ok := netip.AddrFromSlice(ip); ok {
		return addr
	}
	return netip.Addr{}
}

func MsgToIP(msg *dns.Msg) []netip.Addr {
	ips := []netip.Addr{}

	for _, answer := range msg.Answer {
		switch ans := answer.(type) {
		case *dns.AAAA:
			ips = append(ips, IpToAddr(ans.AAAA))
		case *dns.A:
			ips = append(ips, IpToAddr(ans.A))
		}
	}

	return ips
}
