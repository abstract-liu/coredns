package constant

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case DOMAIN:
		return "DOMAIN"
	case DOMAIN_SUFFIX:
		return "DOMAIN-SUFFIX"
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
	FINAL
	TYPE
	FALLBACK
)
