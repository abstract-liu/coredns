package config

import (
	"fmt"
	"github.com/coredns/coredns/plugin/clash/adapter"
	"github.com/coredns/coredns/plugin/clash/common"
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
	nameservers map[string]adapter.Nameserver
	rules       []Rule
}

type Rule interface {
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
	cfg.nameservers = nameservers

	rules, err := parseRules(rawCfg.Rules, nameservers)
	if err != nil {
		return nil, err
	}
	cfg.rules = rules

	return cfg, nil
}

func parseNameservers(cfg *RawClashConfig) (nameservers map[string]adapter.Nameserver, err error) {
	nameservers = make(map[string]adapter.Nameserver)

	// parse nameservers
	for idx, mapping := range cfg.Nameservers {
		ns, err := adapter.ParseNameserver(mapping)
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
		groupName, existName := mapping["name"].(string)
		if !existName {
			return nil, fmt.Errorf("ns group %d: missing name", idx)
		}
		log.Debug("group name: ", groupName)
	}

	return nameservers, nil
}

func parseRules(rulesConfig []string, nameservers map[string]adapter.Nameserver) ([]Rule, error) {
	var rules []Rule

	// parse rules
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
			return nil, fmt.Errorf("Rule[%d] [%s] error: %s", idx, line, parseErr.Error())
		}

		rules = append(rules, parsed)
	}

	return rules, nil
}
