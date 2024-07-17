package clash

import (
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/config"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"os"
	"path/filepath"
	"time"
)

const _pluginName = "clash"

var log = clog.NewWithPlugin(_pluginName)

type PluginConfig struct {
	path string

	modifiedTime time.Time
	size         int64

	clashConfig *config.ClashConfig
}

func init() { plugin.Register(_pluginName, setup) }

func setup(c *caddy.Controller) error {
	clash, err := parseClash(c)
	if err != nil {
		return plugin.Error(_pluginName, err)
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

		cfg, err := parsePluginConfig(c)
		if err != nil {
			return nil, err
		}

		clash, err = NewClash(cfg)
		if err != nil {
			return nil, err
		}
	}

	return clash, nil
}

func parsePluginConfig(c *caddy.Controller) (*PluginConfig, error) {
	config := dnsserver.GetConfig(c)
	pluginConfig := &PluginConfig{}

	args := c.RemainingArgs()
	if len(args) != 1 {
		return nil, c.Errf("invalid number of config files: %d", len(args))
	}
	configFilename := args[0]

	if !filepath.IsAbs(configFilename) && config.Root != "" {
		pluginConfig.path = filepath.Join(config.Root, configFilename)
	}
	s, err := os.Stat(pluginConfig.path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warningf("File does not exist: %s", pluginConfig.path)
		} else {
			return nil, c.Errf("unable to access clash config file '%s': %v", pluginConfig.path, err)
		}
	}
	if s != nil && s.IsDir() {
		log.Warningf("Clash config file %q is a directory", pluginConfig.path)
	}

	if nil != readClashConfig(pluginConfig) {
		return nil, c.Errf("unable to parse clash config file '%s'", pluginConfig.path)
	}

	return pluginConfig, nil
}

func readClashConfig(pluginConfig *PluginConfig) error {
	path := pluginConfig.path
	stat, err := os.Stat(path)
	if err != nil {
		return err
	}
	if pluginConfig.modifiedTime.Equal(stat.ModTime()) && pluginConfig.size == stat.Size() {
		return err
	}
	pluginConfig.modifiedTime = stat.ModTime()
	pluginConfig.size = stat.Size()

	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	pluginConfig.clashConfig, err = config.Parse(data)

	return nil
}
