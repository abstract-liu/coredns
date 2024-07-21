package single

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainKeyword struct {
	*Base
	keyword string
	ns      string
}

func (d *DomainKeyword) RuleType() constant.RuleType {
	return constant.DOMAIN_KEYWORD
}

func (d *DomainKeyword) NS() string {
	return d.ns
}

func (d *DomainKeyword) Match(msg *dns.Msg) (bool, string) {
	domain := msg.Question[0].Name
	return strings.Contains(domain, d.keyword), d.ns
}

func NewDomainKeyword(keyword string, ns string) *DomainKeyword {
	return &DomainKeyword{
		Base:    &Base{},
		keyword: dns.Fqdn(keyword),
		ns:      ns,
	}
}
