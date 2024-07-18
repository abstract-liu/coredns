package outboundgroup

import (
	"github.com/coredns/coredns/plugin/clash/ns"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
)

type GroupBase struct {
	*outbound.Base
	nameservers map[string]ns.Nameserver
}

type GroupBaseOption struct {
	*outbound.BaseOption
	Name        string   `group:"name"`
	Type        string   `group:"type"`
	Nameservers []string `group:"nameservers,omitempty"`
	Use         []string `group:"use,omitempty"`
}

func NewGroupBase(opt *GroupBaseOption) *GroupBase {
	gb := &GroupBase{
		Base: outbound.NewBase(opt.BaseOption),
	}
	gb.nameservers = make(map[string]ns.Nameserver)

	return gb
}
