package rule

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/common"
	"github.com/miekg/dns"
)

type Rule interface {
	RuleType() constant.RuleType
	Adapter() string
	Match(msg *dns.Msg) (bool, string)
}

func ParseRule(ruleType, payload, target string, params []string) (rule Rule, err error) {
	switch ruleType {
	case "DOMAIN":
		rule = common.NewDomain(payload, target)
	case "DOMAIN-SUFFIX":
		rule = common.NewDomainSuffix(payload, target)
	default:
		// TODO: ignore now
		return nil, nil
		// return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
