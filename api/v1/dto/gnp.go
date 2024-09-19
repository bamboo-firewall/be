package dto

import "time"

type GlobalNetworkPolicy struct {
	ID          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Version     uint        `json:"version"`
	Metadata    GNPMetadata `json:"metadata" yaml:"metadata"`
	Spec        GNPSpec     `json:"spec" yaml:"spec"`
	Description string      `json:"description" yaml:"description"`
	CreatedAt   time.Time   `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt" yaml:"updatedAt"`
}

type GNPMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type GNPSpec struct {
	Selector string        `json:"selector" yaml:"selector"`
	Types    []string      `json:"types" yaml:"types"`
	Ingress  []GNPSpecRule `json:"ingress" yaml:"ingress"`
	Egress   []GNPSpecRule `json:"egress" yaml:"egress"`
}

type GNPSpecRule struct {
	Metadata    map[string]string `json:"metadata" yaml:"metadata"`
	Action      string            `json:"action" yaml:"action"`
	Protocol    string            `json:"protocol" yaml:"protocol"`
	NotProtocol string            `json:"notProtocol" yaml:"notProtocol"`
	IPVersion   int               `json:"ipVersion"`
	Source      GNPSpecRuleEntity `json:"source" yaml:"source"`
	Destination GNPSpecRuleEntity `json:"destination" yaml:"destination"`
}

type GNPSpecRuleEntity struct {
	Selector string        `json:"selector" yaml:"selector"`
	Nets     []string      `json:"nets" yaml:"nets"`
	NotNets  []string      `json:"notNets" yaml:"notNets"`
	Ports    []interface{} `json:"ports" yaml:"ports"`
	NotPorts []interface{} `json:"notPorts" yaml:"notPorts"`
}

type CreateGlobalNetworkPolicyInput struct {
	Metadata    GNPMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
	Spec        GNPSpecInput     `json:"spec" yaml:"spec"`
	Description string           `json:"description" yaml:"description"`
}

type GNPMetadataInput struct {
	Name   string            `json:"name" yaml:"name" validate:"required"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type GNPSpecInput struct {
	Selector string             `json:"selector" yaml:"selector"`
	Types    []string           `json:"types" yaml:"types"`
	Ingress  []GNPSpecRuleInput `json:"ingress" yaml:"ingress"`
	Egress   []GNPSpecRuleInput `json:"egress" yaml:"egress"`
}

type GNPSpecRuleInput struct {
	Metadata    map[string]string      `json:"metadata" yaml:"metadata"`
	Action      string                 `json:"action" yaml:"action"`
	Protocol    string                 `json:"protocol" yaml:"protocol"`
	NotProtocol string                 `json:"notProtocol" yaml:"notProtocol"`
	IPVersion   int                    `json:"ipVersion"`
	Source      GNPSpecRuleEntityInput `json:"source" yaml:"source"`
	Destination GNPSpecRuleEntityInput `json:"destination" yaml:"destination"`
}

type GNPSpecRuleEntityInput struct {
	Selector string        `json:"selector" yaml:"selector"`
	Nets     []string      `json:"nets" yaml:"nets"`
	NotNets  []string      `json:"notNets" yaml:"notNets"`
	Ports    []interface{} `json:"ports" yaml:"ports"`
	NotPorts []interface{} `json:"notPorts" yaml:"notPorts"`
}

type DeleteGlobalNetworkPolicyInput struct {
	Metadata GNPMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
}
