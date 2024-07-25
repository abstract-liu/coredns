package hub

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/hub/route"
	clog "github.com/coredns/coredns/plugin/pkg/log"
)

const (
	_defaultRestfulAPIAddress = "0.0.0.0:8080"
)

var log = clog.NewWithPlugin(constant.PluginName)

func Start() error {
	log.Warning("TODO: hub restful api address not set, use default address")
	go route.Start(_defaultRestfulAPIAddress)

	return nil
}
