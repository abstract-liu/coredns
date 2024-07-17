package outbound

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type UdpNs struct {
	*Base
	pureAddr string
	client   *dns.Client
}

type UdpOption struct {
	BasicOption
}

func (ns *UdpNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	resp, rtt, err := ns.client.Exchange(msg, ns.pureAddr)
	clog.Debug("%d", rtt)
	return resp, err
}

func NewUdpNs(option UdpOption) (*UdpNs, error) {
	// get pureAddress from addr
	addrs := strings.Split(option.Address, "//")
	pureAddr := strings.Join(addrs[1:], "")
	return &UdpNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.Udp,
		},
		client:   new(dns.Client),
		pureAddr: pureAddr,
	}, nil
}
