package outboundgroup

import (
	"context"
	"errors"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	"github.com/miekg/dns"
)

type Fallback struct {
	*GroupBase
	DefaultNS  constant.Nameserver
	FallbackNS constant.Nameserver
}

type FallbackOption struct {
	GroupBaseOption
	DefaultNS  string `group:"default-nameserver"`
	FallbackNS string `group:"fallback-nameserver"`
}

func (fb *Fallback) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	return nil, errors.New("fallback ns bad used")
}

func (fb *Fallback) DefaultQuery(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	return fb.DefaultNS.Query(ctx, msg)
}

func (fb *Fallback) FallbackQuery(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	return fb.FallbackNS.Query(ctx, msg)
}

func NewFallback(option *FallbackOption, nameservers map[string]constant.Nameserver) (*Fallback, error) {
	if nameservers[option.DefaultNS] == nil || nameservers[option.FallbackNS] == nil {
		return nil, errors.New("default or fallback nameserver not found")
	}

	return &Fallback{
		GroupBase: &GroupBase{
			Base: outbound.NewBase(&outbound.BaseOption{
				Name:   option.Name,
				NSType: constant.FALLBACK_NS,
			}),
		},
		DefaultNS:  nameservers[option.DefaultNS],
		FallbackNS: nameservers[option.FallbackNS],
	}, nil
}
