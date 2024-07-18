package logic

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/filter"
	"github.com/coredns/coredns/plugin/clash/rule"
	"github.com/miekg/dns"
)

type Fallback struct {
	*rule.Base
	filters         []*filter.Filter
	fallbackAdapter string
}

func (i *Fallback) RuleType() constant.RuleType {
	return constant.FALLBACK
}

func (i *Fallback) Adapter() string {
	return i.fallbackAdapter
}

func (i *Fallback) Match(msg *dns.Msg) (bool, string) {
	// cause it's if else, we should always return true
	return true, i.defaultAdapter
}
