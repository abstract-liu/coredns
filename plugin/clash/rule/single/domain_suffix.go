package single

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainSuffix struct {
	*Base
	suffix string
	ns     string
}

func (d *DomainSuffix) RuleType() constant.RuleType {
	return constant.DOMAIN
}

func (d *DomainSuffix) NS() string {
	return d.ns
}

func (d *DomainSuffix) Match(msg *dns.Msg) (bool, string) {
	domain := msg.Question[0].Name
	return strings.HasSuffix(domain, "."+d.suffix) || domain == d.suffix, d.ns
}

func NewDomainSuffix(suffix string, ns string) *DomainSuffix {
	return &DomainSuffix{
		Base:   &Base{},
		suffix: dns.Fqdn(suffix),
		ns:     ns,
	}
}
