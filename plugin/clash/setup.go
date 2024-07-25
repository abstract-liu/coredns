package clash

import (
	"fmt"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/config"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"os"
	"path/filepath"
)

var log = clog.NewWithPlugin(constant.PluginName)

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

		clashCfg, err := parseClashConfig(pluginCfg.path)
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
	if !filepath.IsAbs(path) {
		rootCorefilePath := c.Dispenser.File()
		rootPath := filepath.Dir(rootCorefilePath)
		pluginConfig.path = filepath.Join(rootPath, path)
	}
	clog.Infof("Parse Plugin Config Success! Plugin Config Path: %s", pluginConfig.path)

	return pluginConfig, nil
}

func parseClashConfig(path string) (*config.ClashConfig, error) {
	stat, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			log.Warningf("File does not exist: %stat", path)
		} else {
			return nil, fmt.Errorf("unable to access clash config file '%s': %v", path, err)
		}
	}
	if stat != nil && stat.IsDir() {
		return nil, fmt.Errorf("clash config file %s is a directory", path)
	}

	fileData, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	clashConfig, err := config.Parse(fileData)
	if nil != err {
		return nil, fmt.Errorf("unable to parse clash config file '%s', %v", path, err)
	}

	return clashConfig, nil
}
