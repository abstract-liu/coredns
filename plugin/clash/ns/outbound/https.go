package outbound

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type HttpsNs struct {
	*Base
}

type HttpsOption struct {
	BaseOption
}

func (ns *HttpsNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	return nil, nil
}

func NewHttpsNs(option HttpsOption) (*HttpsNs, error) {
	return &HttpsNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.HTTPS,
		},
	}, nil
}
