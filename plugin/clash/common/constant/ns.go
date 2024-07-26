package constant

import (
	"context"
	"github.com/miekg/dns"
)

const (
	UDP NameserverType = iota
	TCP
	TLS
	HTTPS
	REJECT

	RANDOM
	ROUND_ROBIN
	FAST
	FALLBACK_NS
)

type Nameserver interface {
	Name() string
	Type() NameserverType
	Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error)
}

type NameserverType int

func (ns NameserverType) String() string {
	switch ns {
	case UDP:
		return "UDP"
	case TCP:
		return "TCP"
	case TLS:
		return "TLS"
	case HTTPS:
		return "HTTPS"
	case REJECT:
		return "REJECT"
	case ROUND_ROBIN:
		return "ROUND_ROBIN"
	case FAST:
		return "FAST"
	case FALLBACK_NS:
		return "FALLBACK"
	default:
		return "Unknown"
	}
}
