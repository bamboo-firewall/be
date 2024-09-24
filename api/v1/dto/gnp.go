package dto

import "time"

type GlobalNetworkPolicy struct {
	ID          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Version     uint        `json:"version"`
	Metadata    GNPMetadata `json:"metadata,omitempty" yaml:"metadata,omitempty"`
	Spec        GNPSpec     `json:"spec" yaml:"spec"`
	Description string      `json:"description,omitempty" yaml:"description"`
	CreatedAt   time.Time   `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt" yaml:"updatedAt"`
}

type GNPMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels,omitempty" yaml:"labels"`
}

type GNPSpec struct {
	Selector string        `json:"selector,omitempty" yaml:"selector"`
	Ingress  []GNPSpecRule `json:"ingress,omitempty" yaml:"ingress"`
	Egress   []GNPSpecRule `json:"egress,omitempty" yaml:"egress"`
}

type GNPSpecRule struct {
	Metadata    map[string]string  `json:"metadata,omitempty" yaml:"metadata"`
	Action      string             `json:"action" yaml:"action"`
	Protocol    string             `json:"protocol,omitempty" yaml:"protocol"`
	NotProtocol string             `json:"notProtocol,omitempty" yaml:"notProtocol"`
	IPVersion   int                `json:"ipVersion"`
	Source      *GNPSpecRuleEntity `json:"source,omitempty" yaml:"source"`
	Destination *GNPSpecRuleEntity `json:"destination,omitempty" yaml:"destination"`
}

type GNPSpecRuleEntity struct {
	Selector string        `json:"selector,omitempty" yaml:"selector"`
	Nets     []string      `json:"nets,omitempty" yaml:"nets"`
	NotNets  []string      `json:"notNets,omitempty" yaml:"notNets"`
	Ports    []interface{} `json:"ports,omitempty" yaml:"ports"`
	NotPorts []interface{} `json:"notPorts,omitempty" yaml:"notPorts"`
}

type CreateGlobalNetworkPolicyInput struct {
	Metadata    GNPMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
	Spec        GNPSpecInput     `json:"spec" yaml:"spec" validate:"required"`
	Description string           `json:"description" yaml:"description"`
}

type GNPMetadataInput struct {
	Name   string            `json:"name" yaml:"name" validate:"required,name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type GNPSpecInput struct {
	Selector string             `json:"selector" yaml:"selector" validate:"omitempty,selector"`
	Ingress  []GNPSpecRuleInput `json:"ingress" yaml:"ingress" validate:"omitempty,min=1,dive"`
	Egress   []GNPSpecRuleInput `json:"egress" yaml:"egress" validate:"omitempty,min=1,dive"`
}

type GNPSpecRuleInput struct {
	Metadata    map[string]string       `json:"metadata" yaml:"metadata"`
	Action      string                  `json:"action" yaml:"action" validate:"required,action"`
	Protocol    string                  `json:"protocol" yaml:"protocol" validate:"omitempty,protocol"`
	NotProtocol string                  `json:"notProtocol" yaml:"notProtocol" validate:"omitempty,protocol"`
	IPVersion   int                     `json:"ipVersion" validate:"required,ip_version"`
	Source      *GNPSpecRuleEntityInput `json:"source" yaml:"source" validate:"omitempty"`
	Destination *GNPSpecRuleEntityInput `json:"destination" yaml:"destination" validate:"omitempty"`
}

type GNPSpecRuleEntityInput struct {
	Selector string        `json:"selector" yaml:"selector" validate:"omitempty,selector"`
	Nets     []string      `json:"nets" yaml:"nets" validate:"omitempty,min=1,unique"`
	NotNets  []string      `json:"notNets" yaml:"notNets" validate:"omitempty,min=1,unique"`
	Ports    []interface{} `json:"ports" yaml:"ports" validate:"omitempty,min=1,unique,dive"`
	NotPorts []interface{} `json:"notPorts" yaml:"notPorts" validate:"omitempty,min=1,unique,dive"`
}

type DeleteGlobalNetworkPolicyInput struct {
	Metadata GNPMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
}
