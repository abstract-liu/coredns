package common

import "github.com/coredns/coredns/plugin/clash/common/constant"

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

func NewDomain(domain string, adapter string) *Domain {
	return &Domain{
		Base:    &Base{},
		domain:  domain,
		adapter: adapter,
	}
}
