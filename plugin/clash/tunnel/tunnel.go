package tunnel

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns/outboundgroup"
	"github.com/coredns/coredns/plugin/clash/rule/logic"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/miekg/dns"
	"sync"
)

var (
	mode      = constant.RULE
	configMux sync.RWMutex

	rules       []constant.Rule
	nameservers = make(map[string]constant.Nameserver)

	log = clog.NewWithPlugin(constant.PluginName)
)

type Tunnel struct {
}

var GlobalTunnel Tunnel

func (t *Tunnel) Handle(ctx context.Context, msg *dns.Msg) (*dns.Msg, error) {
	ns, r, err := matchRuleNs(msg)
	if err != nil {
		return nil, err
	}
	logMatchData(msg, ns, r)

	if r.RuleType() == constant.FALLBACK {
		fallbackRule := r.(*logic.Fallback)
		fallbackNS := ns.(*outboundgroup.Fallback)
		// TODO: reformat to channel style
		defaultAnswer, err := fallbackNS.DefaultQuery(ctx, msg)
		if err == nil {
			if !fallbackRule.ShouldFallback(defaultAnswer) {
				log.Infof("fallback normal, use default ns: [%s]", fallbackNS.DefaultNS.Name())
				return defaultAnswer, nil
			}
		}
		log.Infof("fallback filter miss, use fallback ns: [%s]", fallbackNS.FallbackNS.Name())
		return fallbackNS.FallbackQuery(ctx, msg)
	}
	return ns.Query(ctx, msg)
}

func UpdateRules(newRules []constant.Rule) {
	configMux.Lock()
	rules = newRules
	configMux.Unlock()
}

func UpdateNameservers(newNameservers map[string]constant.Nameserver) {
	configMux.Lock()
	nameservers = newNameservers
	configMux.Unlock()
}

func logMatchData(msg *dns.Msg, ns constant.Nameserver, r constant.Rule) {
	question := msg.Question[0]
	log.Infof("query: [%s]-[%s], match rule: [%s], use ns: [%s]", question.Name, dns.TypeToString[question.Qtype], r.RuleType().String(), ns.Name())
}

func matchRuleNs(r *dns.Msg) (ns constant.Nameserver, rule constant.Rule, err error) {
	switch mode {
	case constant.DIRECT:
	case constant.GLOBAL:
		log.Debug("mode not supported now")
	default:
		ns, rule, err = match(r)
	}
	return
}

func match(msg *dns.Msg) (constant.Nameserver, constant.Rule, error) {
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
			matchNS, ok := nameservers[ada]
			if !ok {
				continue
			}

			return matchNS, r, nil
		}
	}

	return nameservers["DIRECT"], nil, nil
}
