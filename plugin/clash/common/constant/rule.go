package constant

import "github.com/miekg/dns"

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case DOMAIN:
		return "DOMAIN"
	case DOMAIN_SUFFIX:
		return "DOMAIN-SUFFIX"
	case DOMAIN_KEYWORD:
		return "DOMAIN-KEYWORD"
	case FINAL:
		return "FINAL"
	case TYPE:
		return "TYPE"
	case FALLBACK:
		return "FALLBACK"
	default:
		return "UNKNOWN"
	}
}

const (
	DOMAIN RuleType = iota
	DOMAIN_SUFFIX
	DOMAIN_KEYWORD
	FINAL
	TYPE
	FALLBACK
)

type Rule interface {
	RuleType() RuleType
	NS() string
	Match(msg *dns.Msg) (bool, string)
}
