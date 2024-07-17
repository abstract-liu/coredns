package common

import (
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type DomainSuffix struct {
	*Base
	suffix  string
	adapter string
}

func (d *DomainSuffix) RuleType() constant.RuleType {
	return constant.DOMAIN
}

func (d *DomainSuffix) Adapter() string {
	return d.adapter
}

func (d *DomainSuffix) Match(msg *dns.Msg) (bool, string) {
	domain := msg.Question[0].Name
	return strings.HasSuffix(domain, "."+d.suffix) || domain == d.suffix, d.adapter
}

func NewDomainSuffix(suffix string, adapter string) *DomainSuffix {
	return &DomainSuffix{
		Base:    &Base{},
		suffix:  common.RenameToRootDomain(suffix),
		adapter: adapter,
	}
}
