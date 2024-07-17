package adapter

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/adapter/outbound"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/common/structure"
	"strings"
)

type Nameserver interface {
	Name() string
	Type() constant.NameserverType
}

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
	case "tcp":
	default:
		return nil, fmt.Errorf("unsupport nameserver type: %s", nsType)
	}

	return ns, nil
}
