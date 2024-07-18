package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule"
	"github.com/miekg/dns"
)

type Final struct {
	*rule.Base
	adapter string
}

func (f *Final) RuleType() constant.RuleType {
	return constant.FINAL
}

func (f *Final) Match(msg *dns.Msg) (bool, string) {
	return true, f.adapter
}

func (f *Final) Adapter() string {
	return f.adapter
}

func NewFinal(adapter string) *Final {
	return &Final{
		Base:    &rule.Base{},
		adapter: adapter,
	}
}
