package clash

import "github.com/miekg/dns"

type Proxy struct {
	Client *dns.Client
	host   string
}
