package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Type struct {
	*Base
	tp uint16
	ns string
}

func (d *Type) RuleType() constant.RuleType {
	return constant.TYPE
}

func (d *Type) NS() string {
	return d.ns
}

func (d *Type) Match(msg *dns.Msg) (bool, string) {
	return d.tp == msg.Question[0].Qtype, d.ns
}

func NewType(tpStr string, ns string) *Type {
	tp, exist := dns.StringToType[tpStr]
	if !exist {
		return nil
	}

	return &Type{
		Base: &Base{},
		tp:   tp,
		ns:   ns,
	}
}
