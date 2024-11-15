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
	ProtocolTCP  = "tcp"
	ProtocolUDP  = "udp"
	ProtocolICMP = "icmp"
	// ProtocolSCTP Stream Control Transmission Protocol (SCTP) is a network protocol that allows for the reliable transmission of data between two endpoints in a computer network.
	// The Stream Control Transmission Protocol (SCTP) is a computer networking communications protocol in the transport layer of the Internet protocol suite.
	// Originally intended for Signaling System 7 (SS7) message transport in telecommunication, the protocol provides the message-oriented feature of the User Datagram Protocol (UDP),
	// while ensuring reliable, in-sequence transport of messages with congestion control like the Transmission Control Protocol (TCP).
	//Unlike UDP and TCP, the protocol supports multihoming and redundant paths to increase resilience and reliability.
	ProtocolSCTP    = "sctp"
	ProtocolUDPLite = "udplite"

	ProtocolNumTCP  = 6
	ProtocolNumUDP  = 17
	ProtocolNumSCTP = 132
)

func NewMinifyUUID() string {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}
