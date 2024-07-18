package outbound

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type RejectNs struct {
	*Base
}

type RejectOption struct {
	BaseOption
}

func (ns *RejectNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	resp := new(dns.Msg)
	resp.SetReply(msg)
	resp.Answer = []dns.RR{}
	resp.RecursionAvailable = true
	return resp, nil
}

func NewRejectNs() *RejectNs {
	return &RejectNs{
		Base: &Base{
			name:   "REJECT",
			addr:   "reject://127.0.0.1",
			nsType: constant.REJECT,
		},
	}
}

func NewRejectNsWithOption(option RejectOption) *RejectNs {
	return &RejectNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.REJECT,
		},
	}
}
