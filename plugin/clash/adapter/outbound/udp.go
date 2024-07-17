package outbound

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
)

type UdpNs struct {
	*Base
}

type UdpOption struct {
	BasicOption
}

func NewUdpNs(option UdpOption) (*UdpNs, error) {
	return &UdpNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.Udp,
		},
	}, nil
}
