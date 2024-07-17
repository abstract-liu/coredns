package rule

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule/common"
)

type Rule interface {
	RuleType() constant.RuleType
	Adapter() string
}

func ParseRule(ruleType, payload, target string, params []string) (rule Rule, err error) {
	switch ruleType {
	case "DOMAIN":
		rule = common.NewDomain(payload, target)
	default:
		return nil, fmt.Errorf("unknown rule type: %s", ruleType)
	}

	return rule, nil
}
