package constant

const (
	Udp NameserverType = iota
	Tcp
)

type NameserverType int

func (ns NameserverType) String() string {
	switch ns {
	case Udp:
		return "Udp"
	case Tcp:
		return "Tcp"
	default:
		return "Unknown"
	}
}
