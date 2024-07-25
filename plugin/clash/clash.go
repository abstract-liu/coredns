package clash

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/config"
	"github.com/coredns/coredns/plugin/clash/metrics"
	"github.com/coredns/coredns/plugin/clash/tunnel"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"net"
	"strings"
)

func NewClash(cfg *PluginConfig) (*Clash, error) {
	c := &Clash{
		config: cfg,
		tunnel: &tunnel.GlobalTunnel,
	}
	applyConfig(cfg.clashConfig)
	return c, nil
}

func applyConfig(cfg *config.ClashConfig) {
	tunnel.UpdateRules(cfg.Rules)
	tunnel.UpdateNameservers(cfg.Nameservers)
}

type Clash struct {
	tunnel *tunnel.Tunnel
	config *PluginConfig

	Next plugin.Handler
}

func (clash *Clash) LookupStaticHost(host string, hostType constant.HostType) []net.IP {
	host = strings.ToLower(host)
	return clash.config.clashConfig.Hosts.LookupHost(host, hostType)
}

func (clash *Clash) StaticHostExist(host string) bool {
	host = strings.ToLower(host)
	if len(clash.config.clashConfig.Hosts.LookupHost(host, constant.A)) > 0 {
		return true
	}
	if len(clash.config.clashConfig.Hosts.LookupHost(host, constant.AAAA)) > 0 {
		return true
	}
	return false
}

func (clash *Clash) Name() string {
	return constant.PluginName
}

func (clash *Clash) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		log.Error("No question in the request")
		return plugin.NextOrFailure(constant.PluginName, clash.Next, ctx, w, r)
	}

	state := request.Request{W: w, Req: r}
	metrics.Report(state)

	succHandleStatic := clash.handleStaticHost(state, w)
	if succHandleStatic {
		return dns.RcodeSuccess, nil
	}

	response, err := clash.tunnel.Handle(ctx, r)
	if err != nil {
		return plugin.NextOrFailure(constant.PluginName, clash.Next, ctx, w, r)
	}

	if response != nil {
		err = w.WriteMsg(response)
	}

	if clash.Next != nil {
		return plugin.NextOrFailure(constant.PluginName, clash.Next, ctx, w, r)
	}

	return 0, nil
}

func (c *Clash) handleStaticHost(state request.Request, w dns.ResponseWriter) bool {
	ips := []net.IP{}
	answers := []dns.RR{}
	switch state.QType() {
	case dns.TypeA:
		ips = c.LookupStaticHost(state.Name(), constant.A)
		answers = a(state.Name(), 3600, ips)
	case dns.TypeAAAA:
		ips = c.LookupStaticHost(state.Name(), constant.AAAA)
		answers = aaaa(state.Name(), 3600, ips)
	}

	if !c.StaticHostExist(state.Name()) {
		return false
	}

	m := new(dns.Msg)
	m.SetReply(state.Req)
	m.Authoritative = true
	m.Answer = answers
	w.WriteMsg(m)
	return true
}

// OnStartup starts a goroutines for all proxies.
func (c *Clash) OnStartup() (err error) {
	c.start()
	return nil
}

// OnShutdown stops all configured proxies.
func (c *Clash) OnShutdown() error {
	return nil
}

func (c *Clash) start() {
	log.Info("Initializing CoreDNS 'Clash' list update routines...")
	// TODO: Implement the start function, updater
}

func a(zone string, ttl uint32, ips []net.IP) []dns.RR {
	answers := make([]dns.RR, len(ips))
	for i, ip := range ips {
		r := new(dns.A)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeA, Class: dns.ClassINET, Ttl: ttl}
		r.A = ip
		answers[i] = r
	}
	return answers
}

// aaaa takes a slice of net.IPs and returns a slice of AAAA RRs.
func aaaa(zone string, ttl uint32, ips []net.IP) []dns.RR {
	answers := make([]dns.RR, len(ips))
	for i, ip := range ips {
		r := new(dns.AAAA)
		r.Hdr = dns.RR_Header{Name: zone, Rrtype: dns.TypeAAAA, Class: dns.ClassINET, Ttl: ttl}
		r.AAAA = ip
		answers[i] = r
	}
	return answers
}
