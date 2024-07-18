package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Domain struct {
	*Base
	domain string
	ns     string
}

func (d *Domain) RuleType() constant.RuleType {
	return constant.DOMAIN
}

func (d *Domain) NS() string {
	return d.ns
}

func (d *Domain) Match(msg *dns.Msg) (bool, string) {
	return msg.Question[0].Name == d.domain, d.ns
}

func NewDomain(domain string, ns string) *Domain {
	return &Domain{
		Base:   &Base{},
		domain: dns.Fqdn(domain),
		ns:     ns,
	}
}
