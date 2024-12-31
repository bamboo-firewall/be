package resourcemanager

import (
	"context"

	"github.com/bamboo-firewall/be/api/v1/dto"
)

type ResourceType int

const (
	ResourceTypeNone ResourceType = iota
	ResourceTypeHEP
	ResourceTypeGNS
	ResourceTypeGNP
)

type Resource interface {
	Create(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) error
	List(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error)
	Get(ctx context.Context, apiServer APIServer, resource interface{}) (interface{}, error)
	Delete(ctx context.Context, apiServer APIServer, resource interface{}) error
	Validate(ctx context.Context, apiServer APIServer, filePath string, resource interface{}) (interface{}, error)
	GetResourceType() ResourceType
	GetHeader() []string
	GetHeaderMap() map[string]string
}

type APIServer interface {
	CreateHEP(ctx context.Context, input *dto.CreateHostEndpointInput) error
	ListHEPs(ctx context.Context, input *dto.ListHostEndpointsInput) ([]*dto.HostEndpoint, error)
	GetHEP(ctx context.Context, input *dto.GetHostEndpointInput) (*dto.HostEndpoint, error)
	DeleteHEP(ctx context.Context, input *dto.DeleteHostEndpointInput) error
	CreateGNS(ctx context.Context, input *dto.CreateGlobalNetworkSetInput) error
	ListGNSs(ctx context.Context) ([]*dto.GlobalNetworkSet, error)
	GetGNS(ctx context.Context, input *dto.GetGNSInput) (*dto.GlobalNetworkSet, error)
	DeleteGNS(ctx context.Context, input *dto.DeleteGlobalNetworkSetInput) error
	CreateGNP(ctx context.Context, input *dto.CreateGlobalNetworkPolicyInput) error
	ListGNPs(ctx context.Context, input *dto.ListGNPsInput) ([]*dto.GlobalNetworkPolicy, error)
	GetGNP(ctx context.Context, input *dto.GetGNPInput) (*dto.GlobalNetworkPolicy, error)
	DeleteGNP(ctx context.Context, input *dto.DeleteGlobalNetworkPolicyInput) error
	ValidateHostEndpoint(ctx context.Context, input *dto.CreateHostEndpointInput) (*dto.ValidateHostEndpointOutput, error)
	ValidateGlobalNetworkPolicy(ctx context.Context, input *dto.CreateGlobalNetworkPolicyInput) (*dto.ValidateGlobalNetworkPolicyOutput, error)
	ValidateGlobalNetworkSet(ctx context.Context, input *dto.CreateGlobalNetworkSetInput) (*dto.ValidateGlobalNetworkSetOutput, error)
}
