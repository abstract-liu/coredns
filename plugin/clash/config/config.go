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
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"time"
)

const (
	_defaultClashConfigUpdateInterval = 24 * time.Hour
	_defaultRestfulAPIAddress         = "0.0.0.0:8080"
)

var (
	log                      = clog.NewWithPlugin(constant.PluginName)
	clashRemoteConfigFetcher *resource.Fetcher[*ClashConfig]
	_defaultRawConfig        = RawClashConfig{
		Mode:               constant.RULE,
		ExternalController: _defaultRestfulAPIAddress,

		Nameservers:      []map[string]any{},
		NameserverGroups: []map[string]any{},
		Rules:            []string{},
		Filters:          []map[string][]string{},
		Hosts:            []string{},

		GeoXUrl: GeoXUrl{
			Mmdb:    "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.metadb",
			ASN:     "https://github.com/xishang0128/geoip/releases/download/latest/GeoLite2-ASN.mmdb",
			GeoIp:   "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat",
			GeoSite: "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat",
		},
	}
)

type GeneralConfig struct {
	Mode               constant.TunnelMode
	ExternalController string
}

type ClashConfig struct {
	General     *GeneralConfig
	Nameservers map[string]constant.Nameserver
	Rules       []constant.Rule
	Filters     map[string][]constant.Filter
	Hosts       *constant.HostTable

	GeoXUrl  GeoXUrl
	MMDBPath string

	ModifiedTime time.Time
	Size         int64
	Path         string
}

type RawClashConfig struct {
	Mode               constant.TunnelMode `yaml:"mode"`
	ExternalController string              `yaml:"external-controller"`

	Nameservers      []map[string]any      `yaml:"nameservers"`
	NameserverGroups []map[string]any      `yaml:"nameserver-groups"`
	Rules            []string              `yaml:"rules"`
	Filters          []map[string][]string `yaml:"filters"`
	Hosts            []string              `yaml:"hosts"`

	GeoXUrl GeoXUrl `yaml:"geox-url"`
}

type GeoXUrl struct {
	GeoIp   string `yaml:"geoip" json:"geoip"`
	Mmdb    string `yaml:"mmdb" json:"mmdb"`
	ASN     string `yaml:"asn" json:"asn"`
	GeoSite string `yaml:"geosite" json:"geosite"`
}

func ParseClashConfig(path string) (*ClashConfig, error) {
	if common.IsHTTPResource(path) {
		return parseRemoteClashConfig(path)
	} else {
		return parseLocalClashConfig(path)
	}
}

func parseLocalClashConfig(path string) (*ClashConfig, error) {
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

	clashConfig, err := parse(fileData)
	if nil != err {
		return nil, fmt.Errorf("unable to parse clash config file '%s', %v", path, err)
	}

	return clashConfig, nil
}

func parseRemoteClashConfig(path string) (*ClashConfig, error) {
	clashRemoteConfigFetcher = resource.NewFetcher[*ClashConfig]("clash-config", path, _defaultClashConfigUpdateInterval, parse, onUpdateClashConfig)
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

func onUpdateClashConfig(config *ClashConfig) {
	log.Warning("Clash Config Updated, OnUpdate method not implemented yet")
}

func parse(buf []byte) (*ClashConfig, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	cfg := &ClashConfig{}
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

func UnmarshalRawConfig(buf []byte) (*RawClashConfig, error) {
	rawCfg := _defaultRawConfig
	err := yaml.Unmarshal(buf, &rawCfg)
	if err != nil {
		return nil, err
	}

	return &rawCfg, nil
}

func parseGeneralConfig(cfg *RawClashConfig) (*GeneralConfig, error) {
	generalCfg := &GeneralConfig{
		Mode:               cfg.Mode,
		ExternalController: cfg.ExternalController,
	}
	return generalCfg, nil
}

func parseNameservers(cfg *RawClashConfig) (nameservers map[string]constant.Nameserver, err error) {
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

func parseFilters(rawConfig *RawClashConfig) (map[string][]constant.Filter, error) {
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

func parseHosts(rawConfig *RawClashConfig) (*constant.HostTable, error) {
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
