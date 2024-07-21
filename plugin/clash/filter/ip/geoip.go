package ip

import (
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/coredns/coredns/plugin/clash/component/mmdb"
	clog "github.com/coredns/coredns/plugin/pkg/log"
	"net/netip"
	"strings"
)

type GEOIP struct {
	*Base
	code string
}

func (i *GEOIP) FilterType() constant.FilterType {
	return constant.GEOIP
}

func (i *GEOIP) Match(addr netip.Addr) bool {
	codes := mmdb.IPInstance().LookupCode(addr.AsSlice())
	for _, code := range codes {
		if strings.EqualFold(code, i.code) && !addr.IsPrivate() {
			return true
		}
	}
	clog.Infof("[GEOIP] Match failed: %s, ip codes:[%s]", addr.String(), strings.Join(codes, ","))
	return false
}

func NewGEOIP(s string) (*GEOIP, error) {
	ipcidr := &GEOIP{
		Base: &Base{},
		code: s,
	}

	return ipcidr, nil
}
