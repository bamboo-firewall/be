package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type HostEndpoint struct {
	ID          primitive.ObjectID   `bson:"_id"`
	UUID        string               `bson:"uuid"`
	Version     uint                 `bson:"version"`
	Metadata    HostEndpointMetadata `bson:"metadata"`
	Spec        HostEndpointSpec     `bson:"spec"`
	Description string               `bson:"description"`
	CreatedAt   time.Time            `bson:"created_at"`
	UpdatedAt   time.Time            `bson:"updated_at"`
}

type HostEndpointMetadata struct {
	Name   string            `bson:"name"`
	Labels map[string]string `bson:"labels"`
}

type HostEndpointSpec struct {
	InterfaceName string   `bson:"interface_name"`
	IPs           []string `bson:"ips"`
	IPsV4         []string `bson:"ips_v4"`
	IPsV6         []string `bson:"ips_v6"`
}

func (HostEndpoint) CollectionName() string {
	return "host_endpoint"
}
