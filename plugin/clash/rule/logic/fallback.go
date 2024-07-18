package logic

import (
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/single"
	"github.com/miekg/dns"
	"github.com/samber/lo"
	"net/netip"
)

type Fallback struct {
	*single.Base
	filters    []constant.Filter
	fallbackNS string
}

func (i *Fallback) RuleType() constant.RuleType {
	return constant.FALLBACK
}

func (i *Fallback) NS() string {
	return i.fallbackNS
}

func (i *Fallback) Match(msg *dns.Msg) (bool, string) {
	return true, i.fallbackNS
}

func (i *Fallback) ShouldFallback(msg *dns.Msg) bool {
	if ips := common.MsgToIP(msg); len(ips) != 0 {
		return lo.EveryBy(ips, func(ip netip.Addr) bool {
			return i.shouldFallback(ip)
		})
	} else {
		return true
	}
}

func (i *Fallback) shouldFallback(ip netip.Addr) bool {
	for _, f := range i.filters {
		if f.Match(ip) {
			return false
		}
	}
	return true
}

func NewFallback(filters []constant.Filter, fallbackNS string) *Fallback {
	return &Fallback{
		Base:       &single.Base{},
		filters:    filters,
		fallbackNS: fallbackNS,
	}
}
