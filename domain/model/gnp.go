package model

import "github.com/bamboo-firewall/be/pkg/entity"

type CreateGlobalNetworkPolicyInput struct {
	Metadata    GNPMetadataInput
	Spec        GNPSpecInput
	Description string
	FilePath    string
}

type GNPMetadataInput struct {
	Name   string
	Labels map[string]string
}

type GNPSpecInput struct {
	Order    *uint32
	Selector string
	Ingress  []GNPSpecRuleInput
	Egress   []GNPSpecRuleInput
}

type GNPSpecRuleInput struct {
	Metadata    map[string]string
	Action      string
	Protocol    interface{}
	NotProtocol interface{}
	IPVersion   *int
	Source      *GNPSpecRuleEntityInput
	Destination *GNPSpecRuleEntityInput
}

type GNPSpecRuleEntityInput struct {
	Selector string
	Nets     []string
	NotNets  []string
	Ports    []interface{}
	NotPorts []interface{}
}

type ListGNPsInput struct {
	IsOrder bool
}

type PolicyWithRelatedHostEndpoint struct {
	GNP        *entity.GlobalNetworkPolicy
	ParsedHEPs []*ParsedHEP
}
