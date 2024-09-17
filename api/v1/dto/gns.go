package dto

import "time"

type GlobalNetworkSet struct {
	ID          string      `json:"id"`
	UUID        string      `json:"uuid"`
	Version     uint        `json:"version"`
	Metadata    GNSMetadata `json:"metadata"`
	Spec        GNSSpec     `json:"spec"`
	Description string      `json:"description"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
}

type GNSMetadata struct {
	Name      string            `json:"name"`
	IPVersion int               `json:"ipVersion"`
	Labels    map[string]string `json:"labels"`
}

type GNSSpec struct {
	Nets []string `json:"nets"`
}

type CreateGlobalNetworkSetInput struct {
	Metadata    GNSMetadataInput `json:"metadata" validate:"required"`
	Spec        GNSSpecInput     `json:"spec"`
	Description string           `json:"description"`
}

type GNSMetadataInput struct {
	Name      string            `json:"name" validate:"required"`
	IPVersion int               `json:"ipVersion"`
	Labels    map[string]string `json:"labels"`
}

type GNSSpecInput struct {
	Nets []string `json:"nets"`
}

type DeleteGlobalNetworkSetInput struct {
	Metadata GNSMetadataInput `json:"metadata" validate:"required"`
}
