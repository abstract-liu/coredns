package config

import (
	"os"
	"testing"
)

func TestParse(t *testing.T) {
	testFile := "config.yaml"
	buf, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatal(err)
	}
	config, err := Parse(buf)
	if err != nil {
		t.Fatal(err)
	}
	if len(config.Nameservers) != 4 {
		t.Fatalf("expected 4 nameserver, got %d", len(config.Nameservers))
	}
}
