package clash

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/config"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"path/filepath"
)

type PluginConfig struct {
	path string
}

func init() { plugin.Register(constant.PluginName, setup) }

func setup(c *caddy.Controller) error {
	clash, err := parseClash(c)
	if err != nil {
		return plugin.Error(constant.PluginName, err)
	}

	c.OnStartup(func() error {
		return clash.OnStartup()
	})

	c.OnShutdown(func() error {
		return clash.OnShutdown()
	})

	dnsserver.GetConfig(c).AddPlugin(func(next plugin.Handler) plugin.Handler {
		clash.Next = next
		return clash
	})

	return nil
}

func parseClash(c *caddy.Controller) (*Clash, error) {
	var (
		clash *Clash
		i     int
	)

	for c.Next() {
		if i > 0 {
			return nil, plugin.ErrOnce
		}
		i += 1

		pluginCfg, err := parsePluginConfig(c)
		if err != nil {
			return nil, err
		}

		clashCfg, err := config.ParseClashConfig(pluginCfg.path)
		if err != nil {
			return nil, err
		}

		clash, err = NewClash(clashCfg)
		if err != nil {
			return nil, err
		}
	}

	return clash, nil
}

func parsePluginConfig(c *caddy.Controller) (*PluginConfig, error) {
	pluginConfig := &PluginConfig{}

	args := c.RemainingArgs()
	if len(args) != 1 {
		return nil, c.Errf("invalid number of pluginConfig files: %d", len(args))
	}
	path := args[0]
	pluginConfig.path = path
	if !filepath.IsAbs(path) && !common.IsHTTPResource(path) {
		rootCorefilePath := c.Dispenser.File()
		rootPath := filepath.Dir(rootCorefilePath)
		pluginConfig.path = filepath.Join(rootPath, path)
	}
	clog.Infof("Parse Plugin Config Success! Plugin Config Path: %s", pluginConfig.path)

	return pluginConfig, nil
}
