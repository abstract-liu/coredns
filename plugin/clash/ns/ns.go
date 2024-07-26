package ns

import (
	"errors"
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/common/structure"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	"github.com/coredns/coredns/plugin/clash/ns/outboundgroup"
	"strings"
)

var (
	errFormat            = errors.New("format error")
	errType              = errors.New("unsupported type")
	errMissProxy         = errors.New("`use` or `proxies` missing")
	errDuplicateProvider = errors.New("duplicate provider name")
)

func ParseNameserver(mapping map[string]any) (constant.Nameserver, error) {
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
		ns  constant.Nameserver
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
	case "https":
		httpsOption := &outbound.HttpsOption{}
		err = decoder.Decode(mapping, httpsOption)
		if err != nil {
			break
		}
		ns, err = outbound.NewHttpsNs(*httpsOption)
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

	return ns, err
}

func ParseNSGroup(config map[string]any, nameservers map[string]constant.Nameserver) (constant.Nameserver, error) {
	decoder := structure.NewDecoder(structure.Option{TagName: "group", WeaklyTypedInput: true})

	groupOption := &outboundgroup.GroupBaseOption{}
	if err := decoder.Decode(config, groupOption); err != nil {
		return nil, errFormat
	}

	if groupOption.Type == "" || groupOption.Name == "" {
		return nil, errFormat
	}

	groupName := groupOption.Name
	if groupOption.Type != "fallback" && len(groupOption.Nameservers) == 0 {
		return nil, fmt.Errorf("%s: %w", groupName, errMissProxy)
	}

	usedNameservers, err := getNameservers(nameservers, groupOption.Nameservers)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", groupName, err)
	}

	var (
		group constant.Nameserver
	)
	switch groupOption.Type {
	case "select":
	case "fallback":
		fallbackOpt := &outboundgroup.FallbackOption{}
		err = decoder.Decode(config, fallbackOpt)
		if err != nil {
			break
		}
		group, err = outboundgroup.NewFallback(fallbackOpt, nameservers)
	case "load-balance":
	case "random":
	case "round-robin":
		roundRobinOpt := &outboundgroup.RoundRobinOption{}
		err = decoder.Decode(config, roundRobinOpt)
		if err != nil {
			break
		}
		group = outboundgroup.NewRoundRobin(roundRobinOpt, usedNameservers)
	case "fast":
		fastOpt := &outboundgroup.FastGroupOption{}
		err = decoder.Decode(config, fastOpt)
		if err != nil {
			break
		}
		group = outboundgroup.NewFastGroup(fastOpt, usedNameservers)
	default:
		return nil, fmt.Errorf("%w: %s", errType, groupOption.Type)
	}

	return group, err
}

func getNameservers(mapping map[string]constant.Nameserver, list []string) ([]constant.Nameserver, error) {
	var nameservers []constant.Nameserver
	for _, name := range list {
		ns, ok := mapping[name]
		if !ok {
			return nil, fmt.Errorf("'%s' not found", name)
		}
		nameservers = append(nameservers, ns)
	}
	return nameservers, nil
}
