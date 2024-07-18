package ip

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"net/netip"
)

type IPCIDR struct {
	*Base
	ipnet netip.Prefix
}

func (i *IPCIDR) FilterType() constant.FilterType {
	return constant.IP_CIDR
}

func (i *IPCIDR) Match(addr netip.Addr) bool {
	return addr.IsValid() && i.ipnet.Contains(addr)
}

func NewIPCIDR(s string) (*IPCIDR, error) {
	ipnet, err := netip.ParsePrefix(s)
	if err != nil {
		return nil, err
	}

	ipcidr := &IPCIDR{
		Base:  &Base{},
		ipnet: ipnet,
	}

	return ipcidr, nil
}
