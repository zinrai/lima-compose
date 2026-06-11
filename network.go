package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

// ipInterface mirrors a subset of `ip -j -4 addr show` output.
type ipInterface struct {
	Ifname   string       `json:"ifname"`
	AddrInfo []ipAddrInfo `json:"addr_info"`
}

type ipAddrInfo struct {
	Local string `json:"local"`
}

// PrintHostsFormat prints all IPv4 addresses of all instances in /etc/hosts format.
func PrintHostsFormat(compose *Compose) {
	for name := range compose.Instances {
		ifaces, err := getInstanceInterfaces(name)
		if err != nil {
			continue
		}

		for _, iface := range ifaces {
			for _, info := range iface.AddrInfo {
				fmt.Printf("%s\t%s-%s\n", info.Local, name, iface.Ifname)
			}
		}
	}
}

func getInstanceInterfaces(name string) ([]ipInterface, error) {
	cmd := exec.Command("limactl", "shell", name, "ip", "-j", "-4", "addr", "show")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get network info for instance %s: %w", name, err)
	}

	var ifaces []ipInterface
	if err := json.Unmarshal(output, &ifaces); err != nil {
		return nil, fmt.Errorf("failed to parse ip output for instance %s: %w", name, err)
	}

	return ifaces, nil
}
