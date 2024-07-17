package constant

const (
	Udp NameserverType = iota
	Tcp
	Tls
)

type NameserverType int

func (ns NameserverType) String() string {
	switch ns {
	case Udp:
		return "Udp"
	case Tcp:
		return "Tcp"
	case Tls:
		return "Tls"
	default:
		return "Unknown"
	}
}
