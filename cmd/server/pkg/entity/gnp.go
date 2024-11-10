package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	PolicyOrderLowest = ^uint64(0)
)

type RuleAction string

const (
	RuleActionAllow RuleAction = "allow"
	RuleActionDeny  RuleAction = "deny"
	RuleActionLog   RuleAction = "log"
	RuleActionPass  RuleAction = "pass"
)

type GlobalNetworkPolicy struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Version     uint               `bson:"version"`
	Metadata    GNPMetadata        `bson:"metadata"`
	Spec        GNPSpec            `bson:"spec"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type GNPMetadata struct {
	Name   string            `bson:"name"`
	Labels map[string]string `bson:"labels"`
}

type GNPSpec struct {
	Order    uint64        `bson:"order"`
	Selector string        `bson:"selector,omitempty"`
	Ingress  []GNPSpecRule `bson:"ingress,omitempty"`
	Egress   []GNPSpecRule `bson:"egress,omitempty"`
}

type GNPSpecRule struct {
	Metadata    map[string]string  `bson:"metadata,omitempty"`
	Action      string             `bson:"action"`
	IPVersion   IPVersion          `bson:"ip_version"`
	Protocol    string             `bson:"protocol,omitempty"`
	NotProtocol string             `bson:"not_protocol,omitempty"`
	Source      *GNPSpecRuleEntity `bson:"source,omitempty"`
	Destination *GNPSpecRuleEntity `bson:"destination,omitempty"`
}

type GNPSpecRuleEntity struct {
	Selector string        `bson:"selector,omitempty"`
	Nets     []string      `bson:"nets,omitempty"`
	NotNets  []string      `bson:"not_nets,omitempty"`
	Ports    []interface{} `bson:"ports,omitempty"`
	NotPorts []interface{} `bson:"not_ports,omitempty"`
}

func (GlobalNetworkPolicy) CollectionName() string {
	return "global_network_policy"
}
