package ip

import (
	"github.com/stretchr/testify/assert"
	"net/netip"
	"testing"
)

func TestIPCIDR_Match(t *testing.T) {
	ipcidr, _ := NewIPCIDR("192.168.2.1/24")
	addr1, _ := netip.ParseAddr("192.168.2.4")
	assert.True(t, ipcidr.Match(addr1))
	addr2, _ := netip.ParseAddr("192.168.3.1")
	assert.False(t, ipcidr.Match(addr2))
}
