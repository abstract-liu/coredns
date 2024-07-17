package clash

import (
	"testing"

	"github.com/coredns/caddy"
)

func TestSetup(t *testing.T) {
	tests := []struct {
		input     string
		shouldErr bool
	}{
		{"clash clash_config.yml", false},
	}

	for i, test := range tests {
		c := caddy.NewTestController("dns", test.input)
		_, err := parseClash(c)
		if test.shouldErr && err == nil {
			t.Errorf("Test %d: expected error but found %s for input %s", i, err, test.input)
		}
	}
}
