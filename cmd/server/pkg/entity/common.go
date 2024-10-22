package entity

type IPVersion int

const (
	IPVersion4 IPVersion = 4
	IPVersion6 IPVersion = 6
)

type Protocol string

const (
	ProtocolTCP  Protocol = "tcp"
	ProtocolUDP  Protocol = "udp"
	ProtocolICMP Protocol = "icmp"
	ProtocolSCTP Protocol = "sctp"
)
