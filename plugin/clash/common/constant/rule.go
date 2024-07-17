package constant

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case DOMAIN:
		return "DOMAIN"
	case DOMAIN_SUFFIX:
		return "DOMAIN_SUFFIX"
	default:
		return "UNKNOWN"
	}
}

const (
	DOMAIN RuleType = iota
	DOMAIN_SUFFIX
)
