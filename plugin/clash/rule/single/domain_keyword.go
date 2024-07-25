package single

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainKeyword struct {
	*Base
	keyword string
	ns      constant.Nameserver
}

func (d *DomainKeyword) NS() constant.Nameserver {
	return d.ns
}

func (d *DomainKeyword) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	domain := msg.Question[0].Name
	return strings.Contains(domain, d.keyword), d.ns, fmt.Sprintf("%s-%s", d.RuleType().String(), d.keyword)
}

func NewDomainKeyword(keyword string, ns constant.Nameserver) *DomainKeyword {
	return &DomainKeyword{
		Base: &Base{
			RT: constant.DOMAIN_KEYWORD,
		},
		keyword: dns.Fqdn(keyword),
		ns:      ns,
	}
}
