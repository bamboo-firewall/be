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
	InterfaceName string                 `bson:"interfaceName"`
	IPs           []string               `bson:"ips"`
	Ports         []HostEndpointSpecPort `bson:"ports"`
}

type HostEndpointSpecPort struct {
	Name     string `bson:"name"`
	Port     int    `bson:"port"`
	Protocol string `bson:"protocol"`
}

func (HostEndpoint) CollectionName() string {
	return "host_endpoint"
}
