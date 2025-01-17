package clash

import (
	"context"
	"fmt"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/component/mmdb"
	"github.com/coredns/coredns/plugin/clash/hub"
	"github.com/coredns/coredns/plugin/clash/metrics"
	"github.com/coredns/coredns/plugin/clash/tunnel"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"os"
)

var log = clog.NewWithPlugin(constant.PluginName)

type Clash struct {
	tunnel *tunnel.Tunnel
	config *constant.ClashConfig

	Next plugin.Handler
}

func NewClash(cfg *constant.ClashConfig) (*Clash, error) {
	c := &Clash{
		config: cfg,
		tunnel: &tunnel.GlobalTunnel,
	}
	return c, nil
}

func (c *Clash) Name() string {
	return constant.PluginName
}

func (c *Clash) ServeDNS(ctx context.Context, w dns.ResponseWriter, r *dns.Msg) (int, error) {
	if len(r.Question) == 0 {
		log.Error("No question in the request")
		return plugin.NextOrFailure(constant.PluginName, c.Next, ctx, w, r)
	} else if len(r.Question) > 1 {
		log.Error("Multiple questions in the request")
		return plugin.NextOrFailure(constant.PluginName, c.Next, ctx, w, r)
	}

	state := request.Request{W: w, Req: r}
	metrics.Report(state)

	return c.tunnel.Handle(ctx, request.Request{W: w, Req: r}), nil
}

// OnStartup starts a goroutines for all proxies.
func (c *Clash) OnStartup() (err error) {
	c.tunnel.ApplyConfig(c.config)

	if err = c.initMMDB(); err != nil {
		return fmt.Errorf("unable to init mmdb, %v", err)
	}

	if err = hub.Start(c.config.General.ExternalController); err != nil {
		return fmt.Errorf("unable to start hub, %v", err)

	}

	log.Info("Initializing CoreDNS 'Clash' list update routines...")
	return nil
}

// OnShutdown stops all configured proxies.
func (c *Clash) OnShutdown() error {
	return nil
}

func (c *Clash) initMMDB() error {
	if c.config.MMDBPath != "" {
		constant.MMDB_PATH = c.config.MMDBPath
	} else {
		constant.MMDB_PATH = constant.ConfigDir + constant.MMDB_FILE_DEFAULT_PATH
	}
	constant.MMDB_URL = c.config.GeoXUrl.Mmdb

	if _, err := os.Stat(constant.MMDB_PATH); os.IsNotExist(err) {
		log.Infof("Can't find MMDB, start download")
		if err := mmdb.DownloadMMDB(); err != nil {
			return err
		}
	} else {
		log.Infof("Load MMDB Success! MMDB File Path: %s", constant.MMDB_PATH)
	}
	return nil
}
