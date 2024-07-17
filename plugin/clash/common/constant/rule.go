package constant

type RuleType int

func (rt RuleType) String() string {
	switch rt {
	case DOMAIN:
		return "DOMAIN"
	default:
		return "UNKNOWN"
	}
}

const (
	DOMAIN RuleType = iota
)
