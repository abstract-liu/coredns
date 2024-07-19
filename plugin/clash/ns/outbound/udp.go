package outbound

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

const _udpPort = 53

type UdpNs struct {
	*Base
	canonicalAddr string
	client        *dns.Client
}

type UdpOption struct {
	BaseOption
}

func (ns *UdpNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	resp, rtt, err := ns.client.Exchange(msg, ns.canonicalAddr)
	clog.Infof("ns: [%s], query: [%s], rtt: %s", ns.Name(), msg.Question[0].Name, rtt)
	return resp, err
}

func NewUdpNs(option UdpOption) (*UdpNs, error) {
	return &UdpNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.UDP,
		},
		client:        new(dns.Client),
		canonicalAddr: common.CanonicalAddr(option.Address, _udpPort),
	}, nil
}
