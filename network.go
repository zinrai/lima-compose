package main

import (
	"fmt"
	"os/exec"
	"strings"
)

// GetInstanceIP retrieves the IP address of a Lima instance
func GetInstanceIP(name string) (string, error) {
	// Get all network interfaces
	cmd := exec.Command("limactl", "shell", name, "ip", "-4", "addr", "show")

	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get network info for instance %s: %w", name, err)
	}

	// Parse and find the first non-loopback IP
	iface, ip := parseNetworkInfo(string(output))
	if ip == "" {
		return "", fmt.Errorf("no non-loopback IP address found for instance %s", name)
	}

	// Return formatted string: "instance-iface IP"
	return fmt.Sprintf("%s-%s %s", name, iface, ip), nil
}

// parseNetworkInfo extracts the first non-loopback interface and IP from ip addr output
func parseNetworkInfo(output string) (string, string) {
	lines := strings.Split(output, "\n")
	currentIface := ""

	for _, line := range lines {
		// Check for interface line (e.g., "2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP>")
		if iface := extractInterface(line); iface != "" {
			currentIface = iface
			continue
		}

		// Check for IP line
		if currentIface != "" && currentIface != "lo" {
			if ip := extractIP(line); ip != "" {
				return currentIface, ip
			}
		}
	}

	return "", ""
}

// extractInterface extracts interface name from interface line
func extractInterface(line string) string {
	// Interface lines don't start with space and contain colons
	if strings.HasPrefix(line, " ") {
		return ""
	}

	parts := strings.Split(line, ":")
	if len(parts) < 3 {
		return ""
	}

	return strings.TrimSpace(parts[1])
}

// extractIP extracts IP address from inet line
func extractIP(line string) string {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "inet ") {
		return ""
	}

	fields := strings.Fields(line)
	if len(fields) < 2 {
		return ""
	}

	// IP address is in format "192.168.105.2/24"
	ipWithMask := fields[1]
	ip := strings.Split(ipWithMask, "/")[0]

	// Skip localhost
	if strings.HasPrefix(ip, "127.") {
		return ""
	}

	return ip
}

// PrintHostsFormat prints instance IPs in /etc/hosts format
func PrintHostsFormat(compose *Compose) {
	for name := range compose.Instances {
		info, err := GetInstanceIP(name)
		if err != nil {
			// Skip instances without IP
			continue
		}

		// info format is "instance-iface IP"
		parts := strings.Split(info, " ")
		if len(parts) == 2 {
			// Print in hosts format: IP<tab>hostname
			fmt.Printf("%s\t%s\n", parts[1], parts[0])
		}
	}
}
