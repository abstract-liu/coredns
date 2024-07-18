package filter

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"net/netip"
)

type Filter interface {
	FilterType() constant.FilterType
	Match(addr netip.Addr) bool
}

func ParseFilter(filterType, payload string) (filter Filter, err error) {
	switch filterType {
	case "IP-CIDR":
	case "IP-ASN":
	case "GEOIP":
	default:
		return nil, nil
	}

	return filter, nil
}
