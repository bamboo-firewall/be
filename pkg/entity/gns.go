package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	GNSEmpty = GlobalNetworkSet{
		ID:      primitive.NewObjectID(),
		UUID:    NewMinifyUUID(),
		Version: 1,
		Metadata: GNSMetadata{
			Name: "default-empty",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
)

type GlobalNetworkSet struct {
	ID          primitive.ObjectID `bson:"_id"`
	UUID        string             `bson:"uuid"`
	Version     uint               `bson:"version"`
	Metadata    GNSMetadata        `bson:"metadata"`
	Spec        GNSSpec            `bson:"spec"`
	Description string             `bson:"description"`
	FilePath    string             `bson:"file_path"`
	CreatedAt   time.Time          `bson:"created_at"`
	UpdatedAt   time.Time          `bson:"updated_at"`
}

type GNSMetadata struct {
	Name   string            `bson:"name"`
	Labels map[string]string `bson:"labels,omitempty"`
}

type GNSSpec struct {
	Nets   []string `bson:"nets"`
	NetsV4 []string `bson:"nets_v4,omitempty"`
	NetsV6 []string `bson:"nets_v6,omitempty"`
}

func (GlobalNetworkSet) CollectionName() string {
	return "global_network_set"
}
