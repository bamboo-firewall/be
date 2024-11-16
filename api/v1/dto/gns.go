package dto

import "time"

type GlobalNetworkSet struct {
	ID          string      `json:"id" yaml:"id"`
	UUID        string      `json:"uuid" yaml:"uuid"`
	Version     uint        `json:"version" yaml:"version"`
	Metadata    GNSMetadata `json:"metadata" yaml:"metadata"`
	Spec        GNSSpec     `json:"spec" yaml:"spec"`
	Description string      `json:"description" yaml:"description"`
	CreatedAt   time.Time   `json:"createdAt" yaml:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt" yaml:"updatedAt"`
}

type GNSMetadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type GNSSpec struct {
	Nets []string `json:"nets" yaml:"nets"`
}

type CreateGlobalNetworkSetInput struct {
	Metadata    GNSMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
	Spec        GNSSpecInput     `json:"spec" yaml:"spec"`
	Description string           `json:"description" yaml:"description"`
}

type GNSMetadataInput struct {
	Name   string            `json:"name" yaml:"name" validate:"required,name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type GNSSpecInput struct {
	Nets []string `json:"nets" yaml:"nets" validate:"min=1,unique"`
}

type ListGNSsInput struct{}

type GetGNSInput struct {
	Name string `uri:"name" validate:"required"`
}

type DeleteGlobalNetworkSetInput struct {
	Metadata GNSMetadataInput `json:"metadata" yaml:"metadata" validate:"required"`
}
