package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule"
	"github.com/miekg/dns"
)

type Domain struct {
	*rule.Base
	domain  string
	adapter string
}

func (d *Domain) RuleType() constant.RuleType {
	return constant.DOMAIN
}

func (d *Domain) Adapter() string {
	return d.adapter
}

func (d *Domain) Match(msg *dns.Msg) (bool, string) {
	return msg.Question[0].Name == d.domain, d.adapter
}

func NewDomain(domain string, adapter string) *Domain {
	return &Domain{
		Base:    &rule.Base{},
		domain:  dns.Fqdn(domain),
		adapter: adapter,
	}
}
