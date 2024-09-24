package net

import (
	"encoding/json"
	"net"
)

type IPNet struct {
	net.IPNet
}

func (i IPNet) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

func (i *IPNet) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	ip, ipnet, err := ParseCIDROrIP(s)
	if err != nil {
		return err
	}
	i.IP = ip.IP
	i.Mask = ipnet.Mask
	return nil
}

func (i IPNet) Version() int {
	if i.IP.To4() != nil {
		return 4
	} else if i.IP.To16() != nil {
		return 6
	}
	return 0
}

func ParseCIDR(cidr string) (*IP, *IPNet, error) {
	netIP, netIPNet, err := net.ParseCIDR(cidr)
	if netIPNet == nil || err != nil {
		return nil, nil, err
	}
	ip := &IP{netIP}
	ipNet := &IPNet{*netIPNet}

	if ipv4 := ip.IP.To4(); ipv4 != nil {
		ip.IP = ipv4
	}
	return ip, ipNet, nil
}

func ParseCIDROrIP(cidr string) (*IP, *IPNet, error) {
	// parse CIDR
	ip, ipnet, err := ParseCIDR(cidr)
	if err == nil {
		return ip, ipnet, nil
	}

	// if err, parse ip
	ip = ParseIP(cidr)
	if ip == nil {
		return nil, nil, err
	}
	ipnet = ip.Network()
	return ip, ipnet, nil
}
