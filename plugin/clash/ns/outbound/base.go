package outbound

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin(constant.PluginName)

type Base struct {
	name   string
	addr   string
	nsType constant.NameserverType
}

type BaseOption struct {
	Name    string `ns:"name"`
	Address string `ns:"address"`
	NSType  constant.NameserverType
}

func NewBase(option *BaseOption) *Base {
	return &Base{
		name:   option.Name,
		addr:   option.Address,
		nsType: option.NSType,
	}
}

func (b *Base) Name() string {
	return b.name
}

func (b *Base) Type() constant.NameserverType {
	return b.nsType
}
