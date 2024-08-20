package constant

type GeneralConfig struct {
	Mode               TunnelMode
	ExternalController string
}

type ClashConfig struct {
	General     *GeneralConfig
	Nameservers map[string]Nameserver
	Rules       []Rule
	Filters     map[string][]Filter
	Hosts       *HostTable

	GeoXUrl  GeoXUrl
	MMDBPath string
}

type RawClashConfig struct {
	Mode               TunnelMode `yaml:"mode"`
	ExternalController string     `yaml:"external-controller"`

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

var (
	ConfigDir string
)
