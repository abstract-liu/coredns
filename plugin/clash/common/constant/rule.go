package constant

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case DOMAIN:
		return "DOMAIN"
	case DOMAIN_SUFFIX:
		return "DOMAIN_SUFFIX"
	case FINAL:
		return "FINAL"
	case TYPE:
		return "TYPE"
	default:
		return "UNKNOWN"
	}
}

const (
	DOMAIN RuleType = iota
	DOMAIN_SUFFIX
	FINAL
	TYPE
)
