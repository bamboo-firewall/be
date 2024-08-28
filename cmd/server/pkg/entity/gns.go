package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalNetworkSet struct {
	ID        primitive.ObjectID `bson:"_id"`
	UUID      string             `bson:"uuid"`
	Version   uint               `bson:"version"`
	Metadata  GNSMetadata        `bson:"metadata"`
	Spec      GNSSpec            `bson:"spec"`
	CreatedAt time.Time          `bson:"created_at"`
	UpdatedAt time.Time          `bson:"updated_at"`
}

type GNSMetadata struct {
	Name   string            `bson:"name"`
	Labels map[string]string `bson:"labels"`
}

type GNSSpec struct {
	Nets []string `bson:"nets"`
}

func (GlobalNetworkSet) CollectionName() string {
	return "global_network_set"
}
