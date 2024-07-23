package host

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"net"
	"strings"
)

func ParseHostFile() {
}

func ParseHost(host string, ip string) constant.Host {
	addr := parseIP(ip)
	if addr == nil {
		return nil
	}
	isIPV4 := true
	if addr.To4() == nil {
		isIPV4 = false
	}
	name := plugin.Name(host).Normalize()

	return constant.NewDefaultHost(name, []net.IP{addr}, isIPV4)
}

func parseIP(addr string) net.IP {
	if i := strings.Index(addr, "%"); i >= 0 {
		// discard ipv6 zone
		addr = addr[0:i]
	}

	return net.ParseIP(addr)
}
