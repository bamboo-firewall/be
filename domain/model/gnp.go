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
	Selector string
	Types    []string
	Ingress  []GNPSpecRuleInput
	Egress   []GNPSpecRuleInput
}

type GNPSpecRuleInput struct {
	Metadata    map[string]string
	Action      string
	Protocol    string
	Source      GNPSpecRuleEntityInput
	Destination GNPSpecRuleEntityInput
}

type GNPSpecRuleEntityInput struct {
	Nets  []string
	Ports []interface{}
}
