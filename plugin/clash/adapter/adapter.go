package adapter

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin/clash/adapter/outbound"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/common/structure"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"strings"
)

type Nameserver interface {
	Name() string
	Type() constant.NameserverType
	Query(ctx context.Context, msg *dns.Msg) (*dns.Msg, error)
}

var log = clog.NewWithPlugin(constant.PluginName)

func ParseNameserver(mapping map[string]any) (Nameserver, error) {
	decoder := structure.NewDecoder(structure.Option{TagName: "ns", WeaklyTypedInput: true, KeyReplacer: structure.DefaultKeyReplacer})

	address, existAddr := mapping["address"].(string)
	if !existAddr {
		return nil, fmt.Errorf("missing address")
	}
	// extract first dns type: udp://127.0.0.1
	expr := strings.Split(address, "://")
	if len(expr) != 2 {
		return nil, fmt.Errorf("invalid address")
	}
	nsType := expr[0]

	var (
		ns  Nameserver
		err error
	)
	switch nsType {
	case "udp":
		udpOption := &outbound.UdpOption{}
		err = decoder.Decode(mapping, udpOption)
		if err != nil {
			break
		}
		ns, err = outbound.NewUdpNs(*udpOption)
	case "tls":
		tlsOption := &outbound.TlsOption{}
		err = decoder.Decode(mapping, tlsOption)
		if err != nil {
			break
		}
		ns, err = outbound.NewTlsNs(*tlsOption)
	case "reject":
		rejectOption := &outbound.RejectOption{}
		err = decoder.Decode(mapping, rejectOption)
		if err != nil {
			break
		}
		ns = outbound.NewRejectNsWithOption(*rejectOption)
	default:
		return nil, fmt.Errorf("unsupport nameserver type: %s", nsType)
	}

	return ns, nil
}
