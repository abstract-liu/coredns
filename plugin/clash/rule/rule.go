package rule

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/logic"
	"github.com/coredns/coredns/plugin/clash/rule/single"
)

func ParseRule(ruleType, payload, target string, params []string, filters map[string][]constant.Filter) (rule constant.Rule, err error) {
	switch ruleType {
	case "DOMAIN":
		rule = single.NewDomain(payload, target)
	case "DOMAIN-SUFFIX":
		rule = single.NewDomainSuffix(payload, target)
	case "DOMAIN-KEYWORD":
		rule = single.NewDomainKeyword(payload, target)
	case "FINAL":
		rule = single.NewFinal(target)
	case "TYPE":
		rule = single.NewType(payload, target)
	case "FALLBACK":
		if filters == nil || filters[payload] == nil {
			return nil, fmt.Errorf("fallback rule[%s] error: filters not found", payload)
		}
		rule = logic.NewFallback(filters[payload], target)
	default:
		// TODO: ignore now
		return nil, nil
		// return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
