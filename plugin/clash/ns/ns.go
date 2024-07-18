package ns

import (
	"context"
	"errors"
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/common/structure"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	"github.com/coredns/coredns/plugin/clash/ns/outboundgroup"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"strings"
)

var (
	errFormat            = errors.New("format error")
	errType              = errors.New("unsupported type")
	errMissProxy         = errors.New("`use` or `proxies` missing")
	errDuplicateProvider = errors.New("duplicate provider name")
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

func ParseNSGroup(config map[string]any, nameservers map[string]Nameserver) (Nameserver, error) {
	decoder := structure.NewDecoder(structure.Option{TagName: "group", WeaklyTypedInput: true})

	groupOption := &outboundgroup.GroupBaseOption{}
	if err := decoder.Decode(config, groupOption); err != nil {
		return nil, errFormat
	}

	if groupOption.Type == "" || groupOption.Name == "" {
		return nil, errFormat
	}

	groupName := groupOption.Name
	if len(groupOption.Nameservers) == 0 {
		return nil, fmt.Errorf("%s: %w", groupName, errMissProxy)
	}

	var group Nameserver
	switch groupOption.Type {
	case "select":
	case "fallback":
	case "load-balance":
	case "random":
		group = outboundgroup.NewRandom(groupOption)
	default:
		return nil, fmt.Errorf("%w: %s", errType, groupOption.Type)
	}

	return group, nil
}

func getNameservers(mapping map[string]Nameserver, list []string) ([]Nameserver, error) {
	var nameservers []Nameserver
	for _, name := range list {
		ns, ok := mapping[name]
		if !ok {
			return nil, fmt.Errorf("'%s' not found", name)
		}
		nameservers = append(nameservers, ns)
	}
	return nameservers, nil
}
