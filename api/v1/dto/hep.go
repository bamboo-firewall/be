package dto

import (
	"time"
)

type HostEndpoint struct {
	ID          string               `json:"id" yaml:"id"`
	UUID        string               `json:"uuid" yaml:"uuid"`
	Version     uint                 `json:"version" yaml:"version"`
	Metadata    HostEndpointMetadata `json:"metadata" yaml:"metadata"`
	Spec        HostEndpointSpec     `json:"spec" yaml:"spec"`
	Description string               `json:"description" yaml:"description"`
	CreatedAt   time.Time            `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time            `json:"updatedAt" yaml:"updatedAt"`
}

type HostEndpointMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type HostEndpointSpec struct {
	InterfaceName string   `json:"interfaceName" yaml:"interfaceName"`
	TenantID      uint64   `json:"tenantID" yaml:"tenantID"`
	IP            string   `json:"ip" yaml:"ip"`
	IPs           []string `json:"ips" yaml:"ips"`
}

type CreateHostEndpointInput struct {
	Metadata    HostEndpointMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
	Spec        HostEndpointSpecInput     `json:"spec" yaml:"spec" validate:"required"`
	Description string                    `json:"description" yaml:"description"`
}

type HostEndpointMetadataInput struct {
	Name   string            `json:"name" yaml:"name" validate:"omitempty,name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type HostEndpointSpecInput struct {
	InterfaceName string   `json:"interfaceName" yaml:"interfaceName"`
	TenantID      uint64   `json:"tenantID" yaml:"tenantID" validate:"omitempty"`
	IP            string   `json:"ip" yaml:"ip" validate:"omitempty,ip"`
	IPs           []string `json:"ips" yaml:"ips" validate:"min=1,unique,dive,ip"`
}

type ListHEPsInput struct {
	TenantID *uint64 `form:"tenantID" yaml:"tenantID" validate:"omitempty"`
	IP       *string `form:"ip" yaml:"ip" validate:"omitempty,ip"`
}

type GetHostEndpointInput struct {
	TenantID uint64 `uri:"tenantID" yaml:"tenantID" validate:"required"`
	IP       string `uri:"ip" yaml:"ip" validate:"required,ip"`
}

type DeleteHostEndpointInput struct {
	Spec HostEndpointSpecInput `json:"spec" yaml:"spec" validate:"required"`
}

type FetchHostEndpointPoliciesInput struct {
	TenantID *uint64 `form:"tenantID" yaml:"tenantID" validate:"omitempty"`
	IP       *string `form:"ip" yaml:"ip" validate:"omitempty,ip"`
}

type HostEndpointPolicy struct {
	MetaData   HostEndpointPolicyMetadata `json:"metadata"`
	HEP        *HostEndpoint              `json:"hostEndpoint"`
	ParsedGNPs []*ParsedGNP               `json:"parsedGNPs"`
	ParsedHEPs []*ParsedHEP               `json:"parsedHEPs"`
	ParsedGNSs []*ParsedGNS               `json:"parsedGNSs"`
}

type HostEndpointPolicyMetadata struct {
	HEPVersions map[string]uint `json:"hepVersions"`
	GNPVersions map[string]uint `json:"gnpVersions"`
	GNSVersions map[string]uint `json:"gnsVersions"`
}

type ParsedGNP struct {
	UUID          string        `json:"uuid"`
	Version       uint          `json:"version"`
	Name          string        `json:"name"`
	InboundRules  []*ParsedRule `json:"inboundRules"`
	OutboundRules []*ParsedRule `json:"outboundRules"`
}

type ParsedRule struct {
	Action             string      `json:"action"`
	IPVersion          *int        `json:"ipVersion"`
	Protocol           interface{} `json:"protocol"`
	IsProtocolNegative bool        `json:"isProtocolNegative"`
	SrcNets            []string    `json:"srcNets"`
	IsSrcNetNegative   bool        `json:"isSrcNetNegative"`
	SrcGNSUUIDs        []string    `json:"srcGNSUUIDs"`
	SrcHEPUUIDs        []string    `json:"srcHEPUUIDs"`
	SrcPorts           []string    `json:"srcPorts"`
	IsSrcPortNegative  bool        `json:"isSrcPortNegative"`
	DstNets            []string    `json:"dstNets"`
	IsDstNetNegative   bool        `json:"isDstNetNegative"`
	DstGNSUUIDs        []string    `json:"dstGNSUUIDs"`
	DstHEPUUIDs        []string    `json:"dstHEPUUIDs"`
	DstPorts           []string    `json:"dstPorts"`
	IsDstPortNegative  bool        `json:"isDstPortNegative"`
}

type ParsedHEP struct {
	UUID     string   `json:"uuid"`
	TenantID uint64   `json:"tenantID"`
	Name     string   `json:"name"`
	IP       string   `json:"ip"`
	IPsV4    []string `json:"ipsV4"`
	IPsV6    []string `json:"ipsV6"`
}

type ParsedGNS struct {
	UUID   string   `json:"uuid"`
	Name   string   `json:"name"`
	NetsV4 []string `json:"netsV4"`
	NetsV6 []string `json:"netsV6"`
}
