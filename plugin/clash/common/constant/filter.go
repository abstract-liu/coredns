package constant

import "net/netip"

type FilterType int

func (ct FilterType) String() string {
	switch ct {
	case IP_CIDR:
		return "IP-CIDR"
	case IP_ASN:
		return "IP-ASN"
	case GEOIP:
		return "GEOIP"
	default:
		return "UNKNOWN"
	}
}

const (
	IP_CIDR FilterType = iota
	IP_ASN
	GEOIP
)

type Filter interface {
	FilterType() FilterType
	Match(addr netip.Addr) bool
}
