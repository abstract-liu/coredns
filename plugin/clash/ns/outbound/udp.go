package outbound

import "C"
import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"time"
)

const _udpPort = 53
const _timeout = 5 * time.Second
const _udpSize = 4096

type UdpNs struct {
	*Base
	canonicalAddr string
	client        *dns.Client
}

type UdpOption struct {
	BaseOption
}

func (ns *UdpNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	type result struct {
		msg *dns.Msg
		rtt time.Duration
		err error
	}
	ch := make(chan result, 1)
	go func() {
		resp, rtt, err := ns.client.ExchangeContext(ctx, msg, ns.canonicalAddr)
		ch <- result{resp, rtt, err}
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case ret := <-ch:
		if ret.err != nil {
			log.Errorf("udp-ns: [%s], query: [%s], error: %v", ns.Name(), msg.Question[0].Name, ret.err)
		} else {
			log.Debugf("udp-ns: [%s], query: [%s], rtt: %s", ns.Name(), msg.Question[0].Name, ret.rtt)
		}
		return ret.msg, ret.err
	}
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
