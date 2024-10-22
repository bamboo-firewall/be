package resouremanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

func NewGNS() Resource {
	return &gns{}
}

type gns struct {
}

func (s *gns) Create(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.CreateGlobalNetworkSetInput)
	return apiServer.CreateGNS(ctx, r)
}

func (s *gns) Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error) {
	r := resource.(*dto.GetGNSInput)
	return apiServer.GetGNS(ctx, r)
}

func (s *gns) Delete(ctx context.Context, apiServer APIServer, resource interface{}) error {
	r := resource.(*dto.DeleteGlobalNetworkSetInput)
	return apiServer.DeleteGNS(ctx, r)
}

func (s *gns) GetResourceType() ResourceType {
	return ResourceTypeGNS
}
