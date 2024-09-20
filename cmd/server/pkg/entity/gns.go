package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GlobalNetworkSet struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Version     uint               `bson:"version"`
	Metadata    GNSMetadata        `bson:"metadata"`
	Spec        GNSSpec            `bson:"spec"`
	Description string             `bson:"description"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type GNSMetadata struct {
	Name      string            `bson:"name"`
	IPVersion int               `bson:"ip_version"`
	Labels    map[string]string `bson:"labels"`
}

type GNSSpec struct {
	Nets []string `bson:"nets"`
}

func (GlobalNetworkSet) CollectionName() string {
	return "global_network_set"
}
