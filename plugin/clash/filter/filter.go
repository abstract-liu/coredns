package filter

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/filter/ip"
)

func ParseFilter(filterType, payload string) (filter constant.Filter, err error) {
	switch filterType {
	case "IP-CIDR":
		filter, err = ip.NewIPCIDR(payload)
		if err != nil {
			return nil, err
		}
	case "IP-ASN":
	case "GEOIP":
		filter, err = ip.NewGEOIP(payload)
		if err != nil {
			return nil, err
		}
	default:
		return nil, nil
	}

	return filter, nil
}
