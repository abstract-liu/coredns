package tunnel

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
	"strings"
	"sync"
	"time"
)

var (
	log          = clog.NewWithPlugin(constant.PluginName)
	GlobalTunnel = Tunnel{
		mode: constant.RULE,

		nameservers: make(map[string]constant.Nameserver),
		rules:       make([]constant.Rule, 0),
	}
)

type Tunnel struct {
	mode      constant.TunnelMode
	configMux sync.RWMutex

	rules       []constant.Rule
	nameservers map[string]constant.Nameserver
	hosts       *constant.HostTable
}

func (t *Tunnel) Handle(ctx context.Context, r request.Request) int {
	switch t.mode {
	case constant.DIRECT:
	case constant.GLOBAL:
		log.Errorf("mode[%s] not supported yet", t.mode.String())
		return dns.RcodeServerFailure
	default:
	}

	if t.handleStaticHost(ctx, r) {
		return dns.RcodeSuccess
	}

	if t.handleRule(ctx, r) {
		return dns.RcodeSuccess
	}

	return dns.RcodeServerFailure
}

func (t *Tunnel) handleStaticHost(ctx context.Context, req request.Request) bool {
	ips := []net.IP{}
	answers := []dns.RR{}
	switch req.QType() {
	case dns.TypeA:
		ips = t.lookupStaticHost(req.Name(), constant.A)
		answers = common.A(req.Name(), 3600, ips)
	case dns.TypeAAAA:
		ips = t.lookupStaticHost(req.Name(), constant.AAAA)
		answers = common.AAAA(req.Name(), 3600, ips)
	}

	if !t.isStaticHostExist(req.Name()) {
		return false
	}

	m := new(dns.Msg)
	m.SetReply(req.Req)
	m.Authoritative = true
	m.Answer = answers
	req.W.WriteMsg(m)
	log.Infof("Query: [%s]-[%s], Use static host: %v", req.Name(), dns.TypeToString[req.QType()], ips)

	return true
}

func (t *Tunnel) handleRule(ctx context.Context, req request.Request) bool {
	var (
		err         error
		requestMsg  = req.Req
		responseMsg *dns.Msg
	)
	// time how long it take to match the rule
	startTime := time.Now()
	nameserver, matchDetail, err := t.match(requestMsg)
	matchDoneTime := time.Now()
	if err != nil {
		log.Error("match rule error: %v", err)
		return false
	}
	question := requestMsg.Question[0]

	responseMsg, err = nameserver.Query(ctx, requestMsg)
	if err != nil {
		log.Errorf("Query Error: %s, [%s]-[%s], Match rule: [%s], Use ns: [%s]", err, question.Name, dns.TypeToString[question.Qtype], matchDetail, nameserver.Name())
		return false
	} else {
		log.Infof("Query Success: [%s]-[%s], Takes %s to match rule: [%s], Use ns: [%s], Result: %v", question.Name, dns.TypeToString[question.Qtype], matchDoneTime.Sub(startTime).String(), matchDetail, nameserver.Name(), common.RRsToIPs(responseMsg.Answer))
		req.W.WriteMsg(responseMsg)
	}

	return true
}

func (t *Tunnel) lookupStaticHost(host string, hostType constant.HostType) []net.IP {
	host = strings.ToLower(host)
	return t.hosts.LookupHost(host, hostType)
}

func (t *Tunnel) isStaticHostExist(host string) bool {
	host = strings.ToLower(host)
	if len(t.hosts.LookupHost(host, constant.A)) > 0 {
		return true
	}
	if len(t.hosts.LookupHost(host, constant.AAAA)) > 0 {
		return true
	}
	return false
}

func (t *Tunnel) UpdateRules(newRules []constant.Rule) {
	t.configMux.Lock()
	t.rules = newRules
	t.configMux.Unlock()
}

func (t *Tunnel) UpdateNameservers(newNameservers map[string]constant.Nameserver) {
	t.configMux.Lock()
	t.nameservers = newNameservers
	t.configMux.Unlock()
}

func (t *Tunnel) UpdateHosts(newHosts *constant.HostTable) {
	t.configMux.Lock()
	t.hosts = newHosts
	t.configMux.Unlock()
}

func (t *Tunnel) match(msg *dns.Msg) (constant.Nameserver, string, error) {
	t.configMux.RLock()
	defer t.configMux.RUnlock()

	for _, r := range t.rules {
		if matched, matchNS, matchDetail := r.Match(msg); matched {
			return matchNS, matchDetail, nil
		}
	}

	return t.nameservers["DIRECT"], "No Rule Match, Use DIRECT", nil
}
