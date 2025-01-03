package resourcemanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewGNS() Resource {
	return &gns{}
}

type gns struct {
}

func (s *gns) Create(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) error {
	r := resource.(*dto.CreateGlobalNetworkSetInput)
	r.FilePath = filePath
	return apiServer.CreateGNS(ctx, r)
}

func (s *gns) List(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	return apiServer.ListGNSs(ctx)
}

func (s *gns) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetGNSInput)
	return apiServer.GetGNS(ctx, r)
}

func (s *gns) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteGlobalNetworkSetInput)
	return apiServer.DeleteGNS(ctx, r)
}

func (s *gns) Validate(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) (interface{}, error) {
	r := resource.(*dto.CreateGlobalNetworkSetInput)
	r.FilePath = filePath
	return apiServer.ValidateGlobalNetworkSet(ctx, r)
}

func (s *gns) GetResourceType() ResourceType {
	return ResourceTypeGNS
}

func (s *gns) GetHeader() []string {
	return []string{"UUID", "NAME", "NETS", "VERSION"}
}

func (s *gns) GetHeaderMap() map[string]string {
	return map[string]string{
		"UUID":    "{{.UUID}}",
		"NAME":    "{{.Metadata.Name}}",
		"NETS":    "{{.Spec.Nets}}",
		"VERSION": "{{.Version}}",
	}
}
