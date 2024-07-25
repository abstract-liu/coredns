package single

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Type struct {
	*Base
	tp uint16
	ns constant.Nameserver
}

func (d *Type) NS() constant.Nameserver {
	return d.ns
}

func (d *Type) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	return d.tp == msg.Question[0].Qtype, d.ns,
		fmt.Sprintf("%s-%s", d.RuleType().String(), dns.TypeToString[d.tp])
}

func NewType(tpStr string, ns constant.Nameserver) *Type {
	tp, exist := dns.StringToType[tpStr]
	if !exist {
		return nil
	}

	return &Type{
		Base: &Base{
			RT: constant.TYPE,
		},
		tp: tp,
		ns: ns,
	}
}
