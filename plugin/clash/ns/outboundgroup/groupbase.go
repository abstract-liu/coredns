package outboundgroup

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin(constant.PluginName)

type GroupBase struct {
	*outbound.Base
	nameservers []constant.Nameserver
}

type GroupBaseOption struct {
	Name        string   `group:"name"`
	Type        string   `group:"type"`
	Nameservers []string `group:"nameservers,omitempty"`
	Use         []string `group:"use,omitempty"`
}
