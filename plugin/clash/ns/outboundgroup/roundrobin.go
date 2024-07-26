package outboundgroup

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	"github.com/miekg/dns"
	"sync"
)

type RoundRobin struct {
	*GroupBase
	lock sync.Mutex
	idx  int
}

type RoundRobinOption struct {
	GroupBaseOption
}

func (r *RoundRobin) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	r.lock.Lock()
	currentNS := r.nameservers[r.idx]
	r.idx = (r.idx + 1) % len(r.nameservers)
	r.lock.Unlock()

	log.Debugf("RoundRobin query: [%s], use ns: [%s]", msg.Question[0].Name, currentNS.Name())
	return currentNS.Query(ctx, msg)
}

func NewRoundRobin(option *RoundRobinOption, nameservers []constant.Nameserver) *RoundRobin {
	return &RoundRobin{
		GroupBase: &GroupBase{
			Base: outbound.NewBase(&outbound.BaseOption{
				Name:   option.Name,
				NSType: constant.ROUND_ROBIN,
			}),
			nameservers: nameservers,
		},
		idx: 0,
	}
}
