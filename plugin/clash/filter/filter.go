package filter

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
)

func ParseFilter(filterType, payload string) (filter constant.Filter, err error) {
	switch filterType {
	case "IP-CIDR":
	case "IP-ASN":
	case "GEOIP":
	default:
		return nil, nil
	}

	return filter, nil
}
