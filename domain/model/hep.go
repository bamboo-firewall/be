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
	IPs           []string
	Ports         []HostEndpointSpecPortInput
}

type HostEndpointSpecPortInput struct {
	Name     string
	Port     int
	Protocol string
}

type FetchPoliciesInput struct {
	Name string
}

type HostEndPointPolicy struct {
	MetaData       HostEndPointPolicyMetadata
	HEP            *entity.HostEndpoint
	ParsedPolicies []*ParsedPolicy
	ParsedSets     []*ParsedSet
}

type HostEndPointPolicyMetadata struct {
	HEPVersion  uint
	GNPVersions map[string]uint
	GNSVersions map[string]uint
}

type ParsedPolicy struct {
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
	SrcGNSNetNames     []string
	SrcPorts           []string
	IsSrcPortNegative  bool
	DstNets            []string
	IsDstNetNegative   bool
	DstGNSNetNames     []string
	DstPorts           []string
	IsDstPortNegative  bool
}

type ParsedSet struct {
	Name      string
	IPVersion int
	Nets      []string
}
