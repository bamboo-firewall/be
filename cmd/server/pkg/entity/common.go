package entity

import (
	"strings"

	"github.com/google/uuid"
)

const (
	IPVersion4 = 4
	IPVersion6 = 6
)

const (
	ProtocolTCP     = "tcp"
	ProtocolUDP     = "udp"
	ProtocolICMP    = "icmp"
	ProtocolSCTP    = "sctp"
	ProtocolUDPLite = "udplite"

	ProtocolNumTCP  = 6
	ProtocolNumUDP  = 17
	ProtocolNumSCTP = 132
)

func NewMinifyUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
