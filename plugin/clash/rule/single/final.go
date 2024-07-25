package single

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Final struct {
	*Base
	ns constant.Nameserver
}

func (f *Final) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	return true, f.ns, fmt.Sprintf("%s-%s", f.RuleType().String(), f.ns)
}

func (f *Final) NS() constant.Nameserver {
	return f.ns
}

func NewFinal(ns constant.Nameserver) *Final {
	return &Final{
		Base: &Base{
			RT: constant.FINAL,
		},
		ns: ns,
	}
}
