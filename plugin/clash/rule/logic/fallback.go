package logic

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns/outboundgroup"
	"github.com/coredns/coredns/plugin/clash/rule/single"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"github.com/samber/lo"
	"net/netip"
)

type Fallback struct {
	*single.Base
	filters    []constant.Filter
	fallbackNS constant.Nameserver
}

func (i *Fallback) NS() constant.Nameserver {
	return i.fallbackNS
}

func (i *Fallback) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	fallbackNS := i.fallbackNS.(*outboundgroup.Fallback)
	defaultAnswer, err := fallbackNS.DefaultQuery(context.Background(), msg)
	if err == nil {
		if !i.shouldFallback(defaultAnswer) {
			return true, fallbackNS.DefaultNS, "fallback filter hit, use default ns"
		}
	}
	return true, fallbackNS.FallbackNS, "fallback filter miss, use fallback ns"
}

func (i *Fallback) shouldFallback(msg *dns.Msg) bool {
	if ips := common.MsgToIP(msg); len(ips) != 0 {
		return lo.EveryBy(ips, func(ip netip.Addr) bool {
			for _, f := range i.filters {
				if f.Match(ip) {
					return false
				}
			}
			return true
		})
	} else {
		return true
	}
}

func NewFallback(filters []constant.Filter, fallbackNS constant.Nameserver) *Fallback {
	// check if fallbackNS is a fallback group
	if _, ok := fallbackNS.(*outboundgroup.Fallback); !ok {
		clog.Errorf("fallback ns[%s] is not a fallback group", fallbackNS.Name())
		return nil
	}

	return &Fallback{
		Base: &single.Base{
			RT: constant.FALLBACK,
		},
		filters:    filters,
		fallbackNS: fallbackNS,
	}
}
