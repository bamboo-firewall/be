package resourcemanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewHEP() Resource {
	return &hep{}
}

type hep struct {
}

func (h *hep) Create(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) error {
	r := resource.(*dto.CreateHostEndpointInput)
	r.FilePath = filePath
	return apiServer.CreateHEP(ctx, r)
}

func (h *hep) List(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.ListHostEndpointsInput)
	return apiServer.ListHEPs(ctx, r)
}

func (h *hep) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetHostEndpointInput)
	return apiServer.GetHEP(ctx, r)
}

func (h *hep) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteHostEndpointInput)
	return apiServer.DeleteHEP(ctx, r)
}

func (h *hep) Validate(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) (interface{}, error) {
	r := resource.(*dto.CreateHostEndpointInput)
	r.FilePath = filePath
	return apiServer.ValidateHostEndpoint(ctx, r)
}

func (h *hep) GetResourceType() ResourceType {
	return ResourceTypeHEP
}

func (h *hep) GetHeader() []string {
	return []string{"UUID", "NAME", "TENANT_ID", "IP", "IPS", "VERSION"}
}

func (h *hep) GetHeaderMap() map[string]string {
	return map[string]string{
		"UUID":      "{{.UUID}}",
		"NAME":      "{{.Metadata.Name}}",
		"TENANT_ID": "{{.Spec.TenantID}}",
		"IP":        "{{.Spec.IP}}",
		"IPS":       "{{.Spec.IPs}}",
		"VERSION":   "{{.Version}}",
	}
}
