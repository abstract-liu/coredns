package tunnel

import "C"
import (
	"context"
	"github.com/coredns/coredns/plugin/clash/adapter"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/rule"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"sync"
)

var (
	mode      = constant.Rule
	configMux sync.RWMutex

	rules       []rule.Rule
	nameservers = make(map[string]adapter.Nameserver)

	log = clog.NewWithPlugin(constant.PluginName)
)

type Tunnel struct {
}

var GlobalTunnel Tunnel

func (t *Tunnel) Handle(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	ns, r, err := resolveDnsMsg(msg)
	if err != nil {
		return nil, err
	}
	logMatchData(msg, ns, r)

	return ns.Query(ctx, msg)
}

func UpdateRules(newRules []rule.Rule) {
	configMux.Lock()
	rules = newRules
	configMux.Unlock()
}

func UpdateNameservers(newNameservers map[string]adapter.Nameserver) {
	configMux.Lock()
	nameservers = newNameservers
	configMux.Unlock()
}

func logMatchData(msg *dns.Msg, ns adapter.Nameserver, r rule.Rule) {
	question := msg.Question[0]
	log.Infof("query: [%s]-[%s], match rule: [%s], use ns: [%s]", question.Name, dns.TypeToString[question.Qtype], r.RuleType().String(), ns.Name())
}

func resolveDnsMsg(r *dns.Msg) (ns adapter.Nameserver, rule rule.Rule, err error) {
	switch mode {
	case constant.Direct:
	case constant.Global:
		log.Debug("mode not supported now")
	default:
		ns, rule, err = match(r)
	}
	return
}

func match(msg *dns.Msg) (ns adapter.Nameserver, rule rule.Rule, err error) {
	configMux.RLock()
	defer configMux.RUnlock()
	/*
		if node, ok := resolver.DefaultHosts.Search(metadata.Host, false); ok {
			metadata.DstIP, _ = node.RandIP()
			resolved = true
		}
	*/

	if len(msg.Question) != 1 {
		log.Error("dns query more than one")
		return nameservers["DIRECT"], nil, nil
	}

	for _, r := range rules {
		if matched, ada := r.Match(msg); matched {
			adapter, ok := nameservers[ada]
			if !ok {
				continue
			}

			return adapter, r, nil
		}
	}

	return nameservers["DIRECT"], nil, nil
}
