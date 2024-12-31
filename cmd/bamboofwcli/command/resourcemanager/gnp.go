package resourcemanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewGNP() Resource {
	return &gnp{}
}

type gnp struct {
}

func (p *gnp) Create(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) error {
	r := resource.(*dto.CreateGlobalNetworkPolicyInput)
	r.FilePath = filePath
	return apiServer.CreateGNP(ctx, r)
}

func (p *gnp) List(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.ListGNPsInput)
	return apiServer.ListGNPs(ctx, r)
}

func (p *gnp) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetGNPInput)
	return apiServer.GetGNP(ctx, r)
}

func (p *gnp) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteGlobalNetworkPolicyInput)
	return apiServer.DeleteGNP(ctx, r)
}

func (p *gnp) Validate(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) (interface{}, error) {
	r := resource.(*dto.CreateGlobalNetworkPolicyInput)
	r.FilePath = filePath
	return apiServer.ValidateGlobalNetworkPolicy(ctx, r)
}

func (p *gnp) GetResourceType() ResourceType {
	return ResourceTypeGNP
}

func (p *gnp) GetHeader() []string {
	return []string{"UUID", "NAME", "ORDER", "VERSION"}
}

func (p *gnp) GetHeaderMap() map[string]string {
	return map[string]string{
		"UUID":    "{{.UUID}}",
		"NAME":    "{{.Metadata.Name}}",
		"ORDER":   "{{.Spec.Order}}",
		"VERSION": "{{.Version}}",
	}
}
