package constant

type TunnelMode int

const (
	GLOBAL TunnelMode = iota
	DIRECT
	RULE
)

func (tm TunnelMode) String() string {
	switch tm {
	case GLOBAL:
		return "GLOBAL"
	case DIRECT:
		return "DIRECT"
	case RULE:
		return "RULE"
	default:
		return "UNKNOWN"
	}
}
