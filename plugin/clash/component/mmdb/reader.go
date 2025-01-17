package mmdb

import (
	"fmt"
	"github.com/oschwald/maxminddb-golang"
	"net"
	"strings"
)

type geoip2Country struct {
	Country struct {
		IsoCode string `maxminddb:"iso_code"`
	} `maxminddb:"country"`
}

type IPReader struct {
	*maxminddb.Reader
	databaseType
}

func (r IPReader) LookupCode(ipAddress net.IP) []string {
	switch r.databaseType {
	case typeMaxmind:
		var country geoip2Country
		_ = r.Lookup(ipAddress, &country)
		if country.Country.IsoCode == "" {
			return []string{}
		}
		return []string{strings.ToLower(country.Country.IsoCode)}

	case typeSing:
		var code string
		_ = r.Lookup(ipAddress, &code)
		if code == "" {
			return []string{}
		}
		return []string{code}

	case typeMetaV0:
		var record any
		_ = r.Lookup(ipAddress, &record)
		switch record := record.(type) {
		case string:
			return []string{record}
		case []any: // lookup returned type of slice is []any
			result := make([]string, 0, len(record))
			for _, item := range record {
				result = append(result, item.(string))
			}
			return result
		}
		return []string{}

	default:
		panic(fmt.Sprint("unknown geoip database type:", r.databaseType))
	}
}
