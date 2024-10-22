package resouremanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewHEP() Resource {
	return &hep{}
}

type hep struct {
}

func (h *hep) Create(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.CreateHostEndpointInput)
	return apiServer.CreateHEP(ctx, r)
}

func (h *hep) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetHostEndpointInput)
	return apiServer.GetHEP(ctx, r)
}

func (h *hep) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteHostEndpointInput)
	return apiServer.DeleteHEP(ctx, r)
}

func (h *hep) GetResourceType() ResourceType {
	return ResourceTypeHEP
}
