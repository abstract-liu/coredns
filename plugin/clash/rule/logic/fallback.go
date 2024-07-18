package logic

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/filter"
	"github.com/coredns/coredns/plugin/clash/rule/common"
	"github.com/miekg/dns"
)

type Fallback struct {
	*common.Base
	filters    []*filter.Filter
	fallbackNS string
}

func (i *Fallback) RuleType() constant.RuleType {
	return constant.FALLBACK
}

func (i *Fallback) NS() string {
	return i.fallbackNS
}

func (i *Fallback) Match(msg *dns.Msg) (bool, string) {
	// cause it's fallback, we should always return true
	return true, i.fallbackNS
}
