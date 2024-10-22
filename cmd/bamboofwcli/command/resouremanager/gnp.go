package resouremanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewGNP() Resource {
	return &gnp{}
}

type gnp struct {
}

func (p *gnp) Create(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.CreateGlobalNetworkPolicyInput)
	return apiServer.CreateGNP(ctx, r)
}

func (p *gnp) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetGNPInput)
	return apiServer.GetGNP(ctx, r)
}

func (p *gnp) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteGlobalNetworkPolicyInput)
	return apiServer.DeleteGNP(ctx, r)
}

func (p *gnp) GetResourceType() ResourceType {
	return ResourceTypeGNP
}
