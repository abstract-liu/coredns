package single

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainSuffix struct {
	*Base
	suffix string
	ns     constant.Nameserver
}

func (d *DomainSuffix) NS() constant.Nameserver {
	return d.ns
}

func (d *DomainSuffix) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	domain := msg.Question[0].Name
	return strings.HasSuffix(domain, "."+d.suffix) || domain == d.suffix,
		d.ns,
		fmt.Sprintf("%s-%s", d.RuleType().String(), d.suffix)
}

func NewDomainSuffix(suffix string, ns constant.Nameserver) *DomainSuffix {
	return &DomainSuffix{
		Base: &Base{
			RT: constant.DOMAIN,
		},
		suffix: dns.Fqdn(suffix),
		ns:     ns,
	}
}
