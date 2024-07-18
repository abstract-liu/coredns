package rule

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/logic"
	"github.com/coredns/coredns/plugin/clash/rule/single"
)

func ParseRule(ruleType, payload, target string, params []string) (rule constant.Rule, err error) {
	switch ruleType {
	case "DOMAIN":
		rule = single.NewDomain(payload, target)
	case "DOMAIN-SUFFIX":
		rule = single.NewDomainSuffix(payload, target)
	case "FINAL":
		rule = single.NewFinal(target)
	case "TYPE":
		rule = single.NewType(payload, target)
	case "FALLBACK":
		rule = logic.NewFallback(nil, target)
	default:
		// TODO: ignore now
		return nil, nil
		// return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
