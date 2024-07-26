package outboundgroup

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/common/picker"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	"github.com/miekg/dns"
	"time"
)

const (
	_defaultDNSTimeout = 5 * time.Second
)

type FastGroup struct {
	*GroupBase
}

type FastGroupOption struct {
	GroupBaseOption
}

func (f *FastGroup) Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	fast, ctx := picker.WithTimeout[*dns.Msg](ctx, _defaultDNSTimeout)
	defer fast.Close()

	startTime := time.Now()
	for _, client := range f.nameservers {
		client := client // shadow define client to ensure the value captured by the closure will not be changed in the next loop
		fast.Go(func() (*dns.Msg, error) {
			if m, err := client.Query(ctx, msg); err != nil {
				return nil, err
			} else {
				return m, nil
			}
		})
	}

	msg = fast.Wait()
	log.Debugf("Fast query: [%s], rtt: %s", msg.Question[0].Name, time.Since(startTime))
	if msg == nil {
		return nil, fmt.Errorf("fast group query failed, %v", fast.Error())
	} else {
		return msg, nil
	}
}

func NewFastGroup(option *FastGroupOption, nameservers []constant.Nameserver) *FastGroup {
	return &FastGroup{
		GroupBase: &GroupBase{
			Base: outbound.NewBase(&outbound.BaseOption{
				Name:   option.Name,
				NSType: constant.FAST,
			}),
			nameservers: nameservers,
		},
	}
}
