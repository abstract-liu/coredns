package single

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
)

var log = clog.NewWithPlugin(constant.PluginName)

type Base struct {
	RT constant.RuleType
}

func (b *Base) RuleType() constant.RuleType {
	return b.RT
}

func (b *Base) NS() constant.Nameserver {
	log.Errorf("BaseRule NS function not implemented")
	return nil
}

func (b *Base) Match(msg *dns.Msg) (bool, constant.Nameserver, string) {
	log.Errorf("BaseRule Match function not implemented")
	return false, nil, ""
}
