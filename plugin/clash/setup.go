package clash

import (
	"fmt"
	"github.com/coredns/caddy"
	"github.com/coredns/coredns/core/dnsserver"
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/component/resource"
	"github.com/coredns/coredns/plugin/clash/config"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"os"
	"path/filepath"
	"time"
)

const (
	_defaultClashConfigUpdateInterval = 24 * time.Hour
)

var (
	log                      = clog.NewWithPlugin(constant.PluginName)
	clashRemoteConfigFetcher *resource.Fetcher[*config.ClashConfig]
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

		clashCfg, err := parseClashConfig(pluginCfg.path)
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

func parseClashConfig(path string) (*config.ClashConfig, error) {
	if common.IsHTTPResource(path) {
		return parseRemoteClashConfig(path)
	} else {
		return parseLocalClashConfig(path)
	}
}

func parseLocalClashConfig(path string) (*config.ClashConfig, error) {
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

func parseRemoteClashConfig(path string) (*config.ClashConfig, error) {
	clashRemoteConfigFetcher = resource.NewFetcher[*config.ClashConfig]("clash-config", path, _defaultClashConfigUpdateInterval, config.Parse, onUpdateClashConfig)
	clashConfig, err := clashRemoteConfigFetcher.Initial()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch clash config file '%s', %v", path, err)
	}

	return clashConfig, nil
}

func UpdateRemoteClashConfig() error {
	clashConfig, same, err := clashRemoteConfigFetcher.Update()
	if same {
		return nil
	}
	if err != nil {
		return err
	}

	onUpdateClashConfig(clashConfig)
	return nil
}

func onUpdateClashConfig(config *config.ClashConfig) {
	log.Warning("Clash Config Updated, OnUpdate method not implemented yet")
}
