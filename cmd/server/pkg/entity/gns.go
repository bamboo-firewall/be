package entity

import (
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var (
	GNSV4Empty = GlobalNetworkSet{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: GNSMetadata{
			Name: "v4-empty",
		},
		Spec: GNSSpec{
			NetsV4: []string{""},
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	GNSV6Empty = GlobalNetworkSet{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: GNSMetadata{
			Name: "v6-empty",
		},
		Spec: GNSSpec{
			NetsV6: []string{""},
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
