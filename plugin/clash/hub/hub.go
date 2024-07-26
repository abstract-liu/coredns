package hub

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/hub/route"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

var log = clog.NewWithPlugin(constant.PluginName)

func Start(address string) error {
	go route.Start(address)

	return nil
}
