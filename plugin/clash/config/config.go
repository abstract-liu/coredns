package config

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/common"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/ns"
	"github.com/coredns/coredns/plugin/clash/ns/outbound"
	R "github.com/coredns/coredns/plugin/clash/rule"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"gopkg.in/yaml.v3"
	"strings"
)

var log = clog.NewWithPlugin("clash")

var _defaultRawConfig = RawClashConfig{
	Nameservers:      []map[string]any{},
	NameserverGroups: []map[string]any{},
	Rules:            []string{},

	GeoXUrl: GeoXUrl{
		Mmdb:    "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.metadb",
		ASN:     "https://github.com/xishang0128/geoip/releases/download/latest/GeoLite2-ASN.mmdb",
		GeoIp:   "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geoip.dat",
		GeoSite: "https://github.com/MetaCubeX/meta-rules-dat/releases/download/latest/geosite.dat",
	},
}

type ClashConfig struct {
	Nameservers map[string]constant.Nameserver
	Rules       []constant.Rule
}

type RawClashConfig struct {
	Nameservers      []map[string]any `yaml:"nameservers"`
	NameserverGroups []map[string]any `yaml:"nameserver-groups"`
	Rules            []string         `yaml:"rules"`

	GeoXUrl GeoXUrl `yaml:"geox-url"`
}

type GeoXUrl struct {
	GeoIp   string `yaml:"geoip" json:"geoip"`
	Mmdb    string `yaml:"mmdb" json:"mmdb"`
	ASN     string `yaml:"asn" json:"asn"`
	GeoSite string `yaml:"geosite" json:"geosite"`
}

func Parse(buf []byte) (*ClashConfig, error) {
	rawCfg, err := UnmarshalRawConfig(buf)
	if err != nil {
		return nil, err
	}

	return ParseRawConfig(rawCfg)
}

func UnmarshalRawConfig(buf []byte) (*RawClashConfig, error) {
	rawCfg := _defaultRawConfig
	err := yaml.Unmarshal(buf, &rawCfg)
	if err != nil {
		return nil, err
	}

	return &rawCfg, nil
}

func ParseRawConfig(rawCfg *RawClashConfig) (*ClashConfig, error) {
	cfg := &ClashConfig{}

	nameservers, err := parseNameservers(rawCfg)
	if err != nil {
		return nil, err
	}
	cfg.Nameservers = nameservers

	rules, err := parseRules(rawCfg.Rules, nameservers)
	if err != nil {
		return nil, err
	}
	cfg.Rules = rules

	return cfg, nil
}

func parseNameservers(cfg *RawClashConfig) (nameservers map[string]constant.Nameserver, err error) {
	nameservers = make(map[string]constant.Nameserver)
	nameservers["REJECT"] = outbound.NewRejectNs()

	// parse Nameservers
	for idx, mapping := range cfg.Nameservers {
		ns, err := ns.ParseNameserver(mapping)
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

func parseRules(rulesConfig []string, nameservers map[string]constant.Nameserver) ([]constant.Rule, error) {
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
		parsed, parseErr := R.ParseRule(ruleName, payload, target, params)
		if parsed == nil {
			continue
		}
		if parseErr != nil {
			return nil, fmt.Errorf("rule[%d] [%s] error: %s", idx, line, parseErr.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}
