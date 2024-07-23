package constant

import "net"

type HostType int

const (
	A HostType = iota
	AAAA
)

type HostTable struct {
	name4 map[string][]net.IP
	name6 map[string][]net.IP
}

func NewHostTable() *HostTable {
	return &HostTable{
		name4: make(map[string][]net.IP),
		name6: make(map[string][]net.IP),
	}
}

func (ht *HostTable) AddHost(name string, ip []net.IP, hostType HostType) {
	if hostType == A {
		ht.name4[name] = append(ht.name4[name], ip...)
	} else {
		ht.name6[name] = append(ht.name6[name], ip...)
	}
}

func (ht *HostTable) LookupHost(name string, hostType HostType) []net.IP {
	if hostType == A {
		return ht.name4[name]
	}
	return ht.name6[name]
}

func (ht *HostTable) Size() int {
	return len(ht.name4) + len(ht.name6)
}

type Host interface {
	Hostname() string
	Type() HostType
	IPs() []net.IP
}

type DefaultHost struct {
	hostname string
	ips      []net.IP
	hostType HostType
}

func NewDefaultHost(hostname string, ips []net.IP, isIPV4 bool) *DefaultHost {
	hostType := A
	if !isIPV4 {
		hostType = AAAA
	}
	return &DefaultHost{
		hostname: hostname,
		ips:      ips,
		hostType: hostType,
	}
}

func (h *DefaultHost) Hostname() string {
	return h.hostname
}

func (h *DefaultHost) Type() HostType {
	return h.hostType
}

func (h *DefaultHost) IPs() []net.IP {
	return h.ips
}
