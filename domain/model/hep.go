package model

import "github.com/bamboo-firewall/be/cmd/server/pkg/entity"

type CreateHostEndpointInput struct {
	Metadata    HostEndpointMetadataInput
	Spec        HostEndpointSpecInput
	Description string
}

type HostEndpointMetadataInput struct {
	Name   string
	Labels map[string]string
}

type HostEndpointSpecInput struct {
	InterfaceName string
	IP            string
	TenantID      uint64
	IPs           []string
	Ports         []HostEndpointSpecPortInput
}

type HostEndpointSpecPortInput struct {
	Name     string
	Port     int
	Protocol string
}

type FetchHostEndpointPolicyInput struct {
	Name string
}

type HostEndPointPolicy struct {
	MetaData   HostEndPointPolicyMetadata
	HEP        *entity.HostEndpoint
	ParsedGNPs []*ParsedGNP
	ParsedHEPs []*ParsedHEP
	ParsedGNSs []*ParsedGNS
}

type HostEndPointPolicyMetadata struct {
	GNPVersions map[string]uint
	HEPVersions map[string]uint
	GNSVersions map[string]uint
}

type ParsedGNP struct {
	UUID          string
	Version       uint
	Name          string
	InboundRules  []*ParsedRule
	OutboundRules []*ParsedRule
}

type ParsedRule struct {
	Action             string
	IPVersion          int
	Protocol           string
	IsProtocolNegative bool
	SrcNets            []string
	IsSrcNetNegative   bool
	SrcGNSUUIDs        []string
	SrcHEPUUIDs        []string
	SrcPorts           []string
	IsSrcPortNegative  bool
	DstNets            []string
	IsDstNetNegative   bool
	DstGNSUUIDs        []string
	DstHEPUUIDs        []string
	DstPorts           []string
	IsDstPortNegative  bool
}

type ParsedHEP struct {
	UUID     string
	Name     string
	TenantID uint64
	IP       string
	IPsV4    []string
	IPsV6    []string
}

type ParsedGNS struct {
	UUID   string
	Name   string
	NetsV4 []string
	NetsV6 []string
}
