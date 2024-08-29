package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalNetworkPolicy struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Version     uint               `bson:"version"`
	Metadata    GNPMetadata        `bson:"metadata"`
	Spec        GNPSpec            `bson:"spec"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"createdAt"`
	UpdatedAt   time.Time          `bson:"updatedAt"`
}

type GNPMetadata struct {
	Name   string            `bson:"name"`
	Labels map[string]string `bson:"labels"`
}

type GNPSpec struct {
	Selector string        `bson:"selector"`
	Types    []string      `bson:"types"`
	Ingress  []GNPSpecRule `bson:"ingress"`
	Egress   []GNPSpecRule `bson:"egress"`
}

type GNPSpecRule struct {
	Metadata    map[string]string `bson:"metadata"`
	Action      string            `bson:"action"`
	Protocol    string            `bson:"protocol"`
	Source      GNPSpecRuleEntity `bson:"source"`
	Destination GNPSpecRuleEntity `bson:"destination"`
}

type GNPSpecRuleEntity struct {
	Nets  []string      `bson:"nets"`
	Ports []interface{} `bson:"ports"`
}

func (GlobalNetworkPolicy) CollectionName() string {
	return "global_network_policy"
}
