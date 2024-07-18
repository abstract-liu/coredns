package clash

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/config"
	"github.com/coredns/coredns/plugin/clash/tunnel"
	"github.com/miekg/dns"
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

func (clash *Clash) Name() string {
	return constant.PluginName
}

func (clash *Clash) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		log.Error("No question in the request")
		return plugin.NextOrFailure(constant.PluginName, clash.Next, ctx, w, r)
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
