package net

import (
	"encoding/json"
	"net"
)

type IP struct {
	net.IP
}

func (i IP) MarshalJSON() ([]byte, error) {
	s, err := i.MarshalText()
	if err != nil {
		return nil, err
	}
	return json.Marshal(string(s))
}

func (i *IP) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}
	if err := i.UnmarshalText([]byte(s)); err != nil {
		return err
	}
	if ipv4 := i.To4(); ipv4 != nil {
		i.IP = ipv4
	}
	return nil
}

func ParseIP(s string) *IP {
	addr := net.ParseIP(s)
	if addr == nil {
		return nil
	}
	if addr4 := addr.To4(); addr4 != nil {
		addr = addr4
	}
	return &IP{addr}
}

func (i IP) Version() int {
	if i.To4() != nil {
		return 4
	} else if i.To16() != nil {
		return 6
	}
	return 0
}

func (i IP) Network() *IPNet {
	ipnet := &IPNet{}
	if ipv4 := i.IP.To4(); ipv4 != nil {
		ipnet.IP = ipv4
		ipnet.Mask = net.CIDRMask(net.IPv4len*8, net.IPv4len*8)
	} else {
		ipnet.IP = i.IP
		ipnet.Mask = net.CIDRMask(net.IPv6len*8, net.IPv6len*8)
	}
	return ipnet
}
