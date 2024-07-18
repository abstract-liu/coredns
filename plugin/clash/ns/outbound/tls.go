package outbound

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
)

type TlsNs struct {
	*Base
}

type TlsOption struct {
	BaseOption
}

func (ns *TlsNs) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	return nil, nil
}

func NewTlsNs(option TlsOption) (*TlsNs, error) {
	return &TlsNs{
		Base: &Base{
			name:   option.Name,
			addr:   option.Address,
			nsType: constant.Tls,
		},
	}, nil
}
