package common

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"strings"
)

type Domain struct {
	*Base
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
	var domainWithRoot string
	if strings.HasSuffix(domain, ".") {
		domainWithRoot = domain
	} else {
		domainWithRoot = domain + "."
	}
	return &Domain{
		Base:    &Base{},
		domain:  domainWithRoot,
		adapter: adapter,
	}
}
