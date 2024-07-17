package clash

import (
	"context"
	"github.com/coredns/coredns/plugin"
	"github.com/miekg/dns"
)

func NewClash(cfg *PluginConfig) (*Clash, error) {
	c := &Clash{}
	return c, nil
}

type Clash struct {
	proxies []*Proxy

	Next plugin.Handler
}

func (clash *Clash) Name() string {
	return _pluginName
}

func (clash *Clash) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	return plugin.NextOrFailure(_pluginName, clash.Next, ctx, w, r)
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
