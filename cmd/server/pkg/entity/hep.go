package entity

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	DefaultTenantID uint64 = 1
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
	IP            uint32   `json:"ip"`
	TenantID      uint64   `bson:"tenant_id"`
	IPs           []string `bson:"ips"`
	IPsV4         []string `bson:"ips_v4,omitempty"`
	IPsV6         []string `bson:"ips_v6,omitempty"`
}

func (HostEndpoint) CollectionName() string {
	return "host_endpoint"
}
