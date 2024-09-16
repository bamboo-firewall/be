package model

type CreateGlobalNetworkSetInput struct {
	Metadata    GNSMetadataInput `json:"metadata" validate:"required"`
	Spec        GNSSpecInput     `json:"spec"`
	Description string           `json:"description"`
}

type GNSMetadataInput struct {
	Name      string            `json:"name" validate:"required"`
	IPVersion int               `json:"ip_version"`
	Labels    map[string]string `json:"labels"`
}

type GNSSpecInput struct {
	Nets []string `json:"nets"`
}
