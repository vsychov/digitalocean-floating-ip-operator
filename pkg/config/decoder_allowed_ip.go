package config

import (
	"fmt"
	"net"
	"strings"
)

type allowedIngressIpDecoder map[string]bool

func (ipd *allowedIngressIpDecoder) Decode(value string) error {
	input := strings.Split(value, ",")
	ips := map[string]bool{}

	for _, ip := range input {
		netIp := net.ParseIP(ip)
		if netIp == nil {
			return fmt.Errorf("IP is wrong: %s", ip)
		}

		if netIp.To4() == nil {
			return fmt.Errorf("IP is wrong: %s, only ipv4 supported now", ip)
		}

		ips[ip] = true
	}

	if len(ips) == 0 {
		return fmt.Errorf("IP list should not be empty")
	}

	*ipd = allowedIngressIpDecoder(ips)
	return nil
}
