package outboundgroup

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type Random struct {
	*GroupBase
	selected string
}

func (r *Random) Name() string {

}

func (r *Random) Type() constant.NameserverType {

}

func (r *Random) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {

}

func NewRandom(option *GroupCommonOption) *Random {
	return &Random{
		GroupBase: NewGroupBase(option),
		selected:  "COMPATIBLE",
	}
}
