package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule"
	"github.com/miekg/dns"
)

type Final struct {
	*rule.Base
	ns string
}

func (f *Final) RuleType() constant.RuleType {
	return constant.FINAL
}

func (f *Final) Match(msg *dns.Msg) (bool, string) {
	return true, f.ns
}

func (f *Final) NS() string {
	return f.ns
}

func NewFinal(ns string) *Final {
	return &Final{
		Base: &rule.Base{},
		ns:   ns,
	}
}
