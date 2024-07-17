package outbound

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
)

type Base struct {
	name   string
	addr   string
	nsType constant.NameserverType
}

type BasicOption struct {
	Name    string `ns:"name"`
	Address string `ns:"address"`
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Type() constant.NameserverType {
	return b.nsType
}
