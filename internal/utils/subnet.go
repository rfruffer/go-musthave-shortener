package utils

import (
	"net"
)

// IsIPInTrustedSubnet проверяет, находится ли IP адрес в доверенной подсети
func IsIPInTrustedSubnet(ipStr, cidr string) (bool, error) {
	if cidr == "" {
		return false, nil
	}

	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false, nil
	}

	_, subnet, err := net.ParseCIDR(cidr)
	if err != nil {
		return false, err
	}

	return subnet.Contains(ip), nil
}
