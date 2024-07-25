package rule

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/logic"
	"github.com/coredns/coredns/plugin/clash/rule/single"
)

func ParseRule(ruleType, payload, target string, params []string, nameservers map[string]constant.Nameserver, filters map[string][]constant.Filter) (rule constant.Rule, err error) {
	ns, ok := nameservers[target]
	if !ok {
		return nil, fmt.Errorf("nameserver[%s] not found", target)
	}

	switch ruleType {
	case "DOMAIN":
		rule = single.NewDomain(payload, ns)
	case "DOMAIN-SUFFIX":
		rule = single.NewDomainSuffix(payload, ns)
	case "DOMAIN-KEYWORD":
		rule = single.NewDomainKeyword(payload, ns)
	case "FINAL":
		rule = single.NewFinal(ns)
	case "TYPE":
		rule = single.NewType(payload, ns)
	case "FALLBACK":
		if filters == nil || filters[payload] == nil {
			return nil, fmt.Errorf("fallback rule[%s] error: filters not found", payload)
		}
		rule = logic.NewFallback(filters[payload], ns)
	default:
		// TODO: ignore now
		return nil, nil
		// return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
