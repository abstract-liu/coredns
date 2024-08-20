package config

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/component/resource"
	"github.com/coredns/coredns/plugin/clash/filter"
	"github.com/coredns/coredns/plugin/clash/host"
	"github.com/coredns/coredns/plugin/clash/metrics"
	"github.com/coredns/coredns/plugin/clash/ns"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	R "github.com/coredns/coredns/plugin/clash/rule"
	"github.com/coredns/coredns/plugin/clash/tunnel"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"gopkg.in/yaml.v3"
	"strings"
	"time"
)

const (
	_defaultClashConfigUpdateInterval = 24 * time.Hour
	_defaultRestfulAPIAddress         = "0.0.0.0:8080"
)

var (
	log                      = clog.NewWithPlugin(constant.PluginName)
	clashRemoteConfigFetcher *resource.Fetcher[*constant.ClashConfig]
	_defaultRawConfig        = constant.RawClashConfig{
		Mode:               constant.RULE,
		ExternalController: _defaultRestfulAPIAddress,

		Nameservers:      []map[string]any{},
		NameserverGroups: []map[string]any{},
		Rules:            []string{},
		Filters:          []map[string][]string{},
		Hosts:            []string{},

		GeoXUrl: constant.GeoXUrl{
			Mmdb:    "https://minio.abstract-liu.dev/general/geoip.metadb",
			ASN:     "https://github.com/xishang0128/geoip/releases/download/latest/GeoLite2-ASN.mmdb",
			GeoIp:   "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat",
			GeoSite: "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat",
		},
	}
)

func ParseClashConfig(path string) (*constant.ClashConfig, error) {
	clashRemoteConfigFetcher = resource.NewFetcher[*constant.ClashConfig]("clash-config", path, _defaultClashConfigUpdateInterval, parse, onUpdateClashConfig)
	clashConfig, err := clashRemoteConfigFetcher.Initial()
	if err != nil {
		return nil, fmt.Errorf("unable to fetch clash config file '%s', %v", path, err)
	}

	return clashConfig, nil
}

func UpdateRemoteClashConfig() error {
	clashConfig, same, err := clashRemoteConfigFetcher.Update()
	if same {
		log.Debugf("Clash Config doesn't change")
		return nil
	}
	if err != nil {
		return err
	}

	onUpdateClashConfig(clashConfig)
	return nil
}

func onUpdateClashConfig(config *constant.ClashConfig) {
	tunnel.GlobalTunnel.UpdateNameservers(config.Nameservers)
	tunnel.GlobalTunnel.UpdateRules(config.Rules)
	tunnel.GlobalTunnel.UpdateHosts(config.Hosts)
}

func parse(buf []byte) (*constant.ClashConfig, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	cfg := &constant.ClashConfig{}
	generalConfig, err := parseGeneralConfig(rawCfg)
	if err != nil {
		return nil, err
	}
	cfg.General = generalConfig

	nameservers, err := parseNameservers(rawCfg)
	if err != nil {
		return nil, err
	}
	cfg.Nameservers = nameservers

	filters, err := parseFilters(rawCfg)
	if err != nil {
		return nil, err
	}
	cfg.Filters = filters

	rules, err := parseRules(rawCfg.Rules, nameservers, filters)
	if err != nil {
		return nil, err
	}
	cfg.Rules = rules

	hosts, err := parseHosts(rawCfg)
	if err != nil {
		return nil, err
	}
	cfg.Hosts = hosts

	cfg.GeoXUrl = rawCfg.GeoXUrl

	clog.Infof("Parse Clash Config Success! Total with %d nameservers, %d rules, %d filters, %d hosts", len(nameservers), len(rules), len(filters), hosts.Size())
	return cfg, nil
}

func UnmarshalRawConfig(buf []byte) (*constant.RawClashConfig, error) {
	rawCfg := _defaultRawConfig
	err := yaml.Unmarshal(buf, &rawCfg)
	if err != nil {
		return nil, err
	}

	return &rawCfg, nil
}

func parseGeneralConfig(cfg *constant.RawClashConfig) (*constant.GeneralConfig, error) {
	generalCfg := &constant.GeneralConfig{
		Mode:               cfg.Mode,
		ExternalController: cfg.ExternalController,
	}
	return generalCfg, nil
}

func parseNameservers(cfg *constant.RawClashConfig) (nameservers map[string]constant.Nameserver, err error) {
	nameservers = make(map[string]constant.Nameserver)
	nameservers["REJECT"] = outbound.NewRejectNs()
	nameservers["reject"] = outbound.NewRejectNs()

	// parse Nameservers
	for idx, mapping := range cfg.Nameservers {
		ns, err := ns.ParseNameserver(mapping)
		if nil == ns {
			continue
		}
		if err != nil {
			return nil, fmt.Errorf("nameserver %d: %w", idx, err)
		}

		if _, exist := nameservers[ns.Name()]; exist {
			return nil, fmt.Errorf("ns %s is the duplicate name", ns.Name())
		}
		nameservers[ns.Name()] = ns
	}

	// parse nameserver groups
	for idx, mapping := range cfg.NameserverGroups {
		group, err := ns.ParseNSGroup(mapping, nameservers)
		if err != nil {
			return nil, fmt.Errorf("nsgroup[%d]: %w", idx, err)
		}

		groupName := group.Name()
		if _, exist := nameservers[groupName]; exist {
			return nil, fmt.Errorf("nsgroup %s: the duplicate name", groupName)
		}

		nameservers[groupName] = group
	}

	return nameservers, nil
}

func parseRules(rulesConfig []string, nameservers map[string]constant.Nameserver, filters map[string][]constant.Filter) ([]constant.Rule, error) {
	var rules []constant.Rule

	// parse Rules
	// rule in format: ruleType(aka:ruleName), payload, target, params...
	for idx, line := range rulesConfig {
		rule := common.TrimArr(strings.Split(line, ","))
		var (
			payload  string
			target   string
			params   []string
			ruleName = strings.ToUpper(rule[0])
		)

		l := len(rule)

		if l < 2 {
			return nil, fmt.Errorf("Rule[%d] [%s] error: format invalid", idx, line)
		}
		if l < 4 {
			rule = append(rule, make([]string, 4-l)...)
		}
		if ruleName == "MATCH" {
			l = 2
		}
		if l >= 3 {
			l = 3
			payload = rule[1]
		}
		target = rule[l-1]
		params = rule[l:]

		params = common.TrimArr(params)
		parsed, parseErr := R.ParseRule(ruleName, payload, target, params, nameservers, filters)
		if parseErr != nil {
			return nil, fmt.Errorf("rule[%d] [%s] error: %s", idx, line, parseErr.Error())
		}
		if parsed == nil {
			continue
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}

func parseFilters(rawConfig *constant.RawClashConfig) (map[string][]constant.Filter, error) {
	filters := make(map[string][]constant.Filter, len(rawConfig.Filters))
	for _, filterGroup := range rawConfig.Filters {
		// check element in filterGroup only one and get key
		if len(filterGroup) != 1 {
			return nil, fmt.Errorf("filter group format invalid")
		}

		filterName := extractFilterName(filterGroup)
		filterNum := len(filterGroup[filterName])
		if filterNum == 0 {
			return nil, fmt.Errorf("filter group %s is empty", filterName)
		}

		fs := make([]constant.Filter, filterNum)
		for idx, rawFilter := range filterGroup[filterName] {
			rawFilterArray := common.TrimArr(strings.Split(rawFilter, ","))
			if len(rawFilterArray) != 2 {
				return nil, fmt.Errorf("filter[%d] %s format invalid", idx, rawFilter)
			}

			filterType := strings.ToUpper(rawFilterArray[0])
			payload := rawFilterArray[1]
			f, err := filter.ParseFilter(filterType, payload)
			if err != nil {
				return nil, fmt.Errorf("filter[%d] %s error: %s", idx, rawFilter, err.Error())
			}
			fs[idx] = f
		}

		filters[filterName] = fs
	}
	return filters, nil
}

func parseHosts(rawConfig *constant.RawClashConfig) (*constant.HostTable, error) {
	hosts := constant.NewHostTable()
	for idx, rawHost := range rawConfig.Hosts {
		hostElements := common.TrimArr(strings.Split(rawHost, ","))
		if len(hostElements) == 1 {
			clog.Warningf("File hosts not supported yet: %s", rawHost)
			continue
		} else if len(hostElements) == 2 {
			parsedHost := host.ParseHost(hostElements[0], hostElements[1])
			if parsedHost != nil {
				hosts.AddHost(parsedHost.Hostname(), parsedHost.IPs(), parsedHost.Type())
			} else {
				return nil, fmt.Errorf("host[%d] %s format invalid", idx, rawHost)
			}
		} else {
			return nil, fmt.Errorf("host[%d] %s format invalid", idx, rawHost)
		}
	}
	metrics.HostEntries.WithLabelValues().Set(float64(hosts.Size()))
	return hosts, nil
}

func extractFilterName(filterGroup map[string][]string) string {
	var filterName string
	for k := range filterGroup {
		filterName = k
		break
	}
	return filterName
}
