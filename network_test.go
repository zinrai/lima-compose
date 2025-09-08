package main

import (
	"testing"
)

// Test IP address extraction from ip addr output
// This helps troubleshoot when 'ips' command fails
func TestParseNetworkInfo(t *testing.T) {
	// Typical ip addr show output
	input := `1: lo: <LOOPBACK,UP,LOWER_UP> mtu 65536 qdisc noqueue state UNKNOWN group default qlen 1000
    inet 127.0.0.1/8 scope host lo
       valid_lft forever preferred_lft forever
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 qdisc mq state UP group default qlen 1000
    inet 192.168.5.15/24 brd 192.168.5.255 scope global dynamic eth0
       valid_lft 86389sec preferred_lft 86389sec`

	iface, ip := parseNetworkInfo(input)

	if iface != "eth0" {
		t.Errorf("Expected interface eth0, got %s", iface)
	}

	if ip != "192.168.5.15" {
		t.Errorf("Expected IP 192.168.5.15, got %s", ip)
	}
}
