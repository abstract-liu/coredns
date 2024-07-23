package metrics

import (
	"github.com/coredns/coredns/plugin"
	"github.com/coredns/coredns/request"
	"github.com/miekg/dns"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"net"
)

var (
	RequestByHostCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Namespace: plugin.Namespace,
		Subsystem: "clash",
		Name:      "requests_by_host_total",
		Help:      "Counter of DNS requests per hostname",
	}, []string{"hostname", "type", "remote_addr"})

	HostEntries = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: plugin.Namespace,
		Subsystem: "clash",
		Name:      "host_entries",
		Help:      "The combined number of entries in hosts and Corefile.",
	}, []string{})
)

func Report(req request.Request) {
	if len(req.Req.Question) == 0 {
		return
	}

	qType := dns.Type(req.QType()).String()
	hostname := req.Req.Question[0].Name
	remoteHost, _, _ := net.SplitHostPort(req.RemoteAddr())
	RequestByHostCount.WithLabelValues(hostname, qType, remoteHost).Inc()
}
