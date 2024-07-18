package ns

import (
	"context"
	"github.com/coredns/coredns/plugin/clash/common/constant"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exchange(ns constant.Nameserver, domain string, tp uint16) (*dns.Msg, error) {
	query := &dns.Msg{}
	query.SetQuestion(dns.Fqdn(domain), tp)
	return ns.Query(context.Background(), query)
}

func generalTest(ns constant.Nameserver, t *testing.T) {
	rmsg, err := exchange(ns, "1.1.1.1.nip.io", dns.TypeA)
	assert.NoError(t, err)
	assert.NotEmptyf(t, rmsg, "response emty")
	assert.NotZero(t, rmsg.Answer)
	record := rmsg.Answer[0].(*dns.A)
	assert.Equal(t, record.A.String(), "1.1.1.1")

	rmsg, err = exchange(ns, "www.baidu.com", dns.TypeA)
	assert.NoError(t, err)
	assert.NotEmptyf(t, rmsg, "response emty")
	assert.NotZero(t, rmsg.Answer)
}

func TestUDPNS(t *testing.T) {
	config := map[string]any{
		"name":    "aliyun",
		"address": "udp://223.5.5.5:53",
	}

	ns, err := ParseNameserver(config)
	assert.NoError(t, err)
	generalTest(ns, t)
}

func TestRoundRobin(t *testing.T) {
	aliyun1Config := map[string]any{
		"name":    "aliyun-1",
		"address": "udp://223.5.5.5",
	}
	aliyun2Config := map[string]any{
		"name":    "aliyun-2",
		"address": "udp://223.6.6.6",
	}
	ns1, _ := ParseNameserver(aliyun1Config)
	ns2, _ := ParseNameserver(aliyun2Config)
	nss := map[string]constant.Nameserver{
		"aliyun-1": ns1,
		"aliyun-2": ns2,
	}

	config := map[string]any{
		"name": "round-robin",
		"type": "round-robin",
		"nameservers": []string{
			"aliyun-1",
			"aliyun-2",
		},
	}
	ns, err := ParseNSGroup(config, nss)
	assert.NoError(t, err)
	generalTest(ns, t)
}
