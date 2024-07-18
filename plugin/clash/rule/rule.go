package rule

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/common"
)

func ParseRule(ruleType, payload, target string, params []string) (rule constant.Rule, err error) {
	switch ruleType {
	case "DOMAIN":
		rule = common.NewDomain(payload, target)
	case "DOMAIN-SUFFIX":
		rule = common.NewDomainSuffix(payload, target)
	case "FINAL":
		rule = common.NewFinal(target)
	case "TYPE":
		rule = common.NewType(payload, target)
	default:
		// TODO: ignore now
		return nil, nil
		// return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
