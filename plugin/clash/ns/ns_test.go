package ns

import (
	"context"
	"github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
	"testing"
)

func exchange(ns Nameserver, domain string, tp uint16) (*dns.Msg, error) {
	query := &dns.Msg{}
	query.SetQuestion(dns.Fqdn(domain), tp)
	return ns.Query(context.Background(), query)
}

func TestUdpNs(t *testing.T) {
	config := map[string]any{
		"name":    "aliyun",
		"address": "udp://223.5.5.5:53",
	}

	ns, err := ParseNameserver(config)
	assert.NoError(t, err)

	rmsg, err := exchange(ns, "1.1.1.1.nip.io", dns.TypeA)
	assert.NoError(t, err)
	assert.NotEmptyf(t, rmsg, "response emty")
	assert.NotZero(t, rmsg.Answer)
	record := rmsg.Answer[0].(*dns.A)
	assert.Equal(t, record.A.String(), "1.1.1.1")

}
