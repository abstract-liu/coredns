package single

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Domain struct {
	*Base
	domain string
	ns     constant.Nameserver
}

func (d *Domain) NS() constant.Nameserver {
	return d.ns
}

func (d *Domain) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	return msg.Question[0].Name == d.domain, d.ns, fmt.Sprintf("%s-%s", d.RuleType().String(), d.domain)
}

func NewDomain(domain string, ns constant.Nameserver) *Domain {
	return &Domain{
		Base: &Base{
			RT: constant.DOMAIN,
		},
		domain: dns.Fqdn(domain),
		ns:     ns,
	}
}
