package model

type CreateGlobalNetworkPolicyInput struct {
	Metadata    GNPMetadataInput
	Spec        GNPSpecInput
	Description string
}

type GNPMetadataInput struct {
	Name   string
	Labels map[string]string
}

type GNPSpecInput struct {
	Order    *uint64
	Selector string
	Ingress  []GNPSpecRuleInput
	Egress   []GNPSpecRuleInput
}

type GNPSpecRuleInput struct {
	Metadata    map[string]string
	Action      string
	Protocol    string
	NotProtocol string
	IPVersion   int
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

type ListGNPInput struct {
	IsOrder bool
}
