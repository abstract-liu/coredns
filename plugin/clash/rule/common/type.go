package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule"
	"github.com/miekg/dns"
)

type Type struct {
	*rule.Base
	tp      uint16
	adapter string
}

func (d *Type) RuleType() constant.RuleType {
	return constant.TYPE
}

func (d *Type) Adapter() string {
	return d.adapter
}

func (d *Type) Match(msg *dns.Msg) (bool, string) {
	return d.tp == msg.Question[0].Qtype, d.adapter
}

func NewType(tpStr string, adapter string) *Type {
	tp, exist := dns.StringToType[tpStr]
	if !exist {
		return nil
	}

	return &Type{
		Base:    &rule.Base{},
		tp:      tp,
		adapter: adapter,
	}
}
