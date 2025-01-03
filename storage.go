package be

import (
	"context"

	"github.com/bamboo-firewall/be/domain/model"
	"github.com/bamboo-firewall/be/pkg/entity"
	"github.com/bamboo-firewall/be/pkg/httpbase/ierror"
)

type Storage interface {
	UpsertHostEndpoint(ctx context.Context, hep *entity.HostEndpoint) *ierror.CoreError
	GetHostEndpoint(ctx context.Context, input *model.GetHostEndpointInput) (*entity.HostEndpoint, *ierror.CoreError)
	DeleteHostEndpoint(ctx context.Context, tenantID uint64, ip uint32) *ierror.CoreError
	ListHostEndpoints(ctx context.Context, input *model.ListHostEndpointsInput) ([]*entity.HostEndpoint, *ierror.CoreError)
	UpsertGroupPolicy(ctx context.Context, gnp *entity.GlobalNetworkPolicy) *ierror.CoreError
	GetGNPByName(ctx context.Context, name string) (*entity.GlobalNetworkPolicy, *ierror.CoreError)
	DeleteGNPByName(ctx context.Context, name string) *ierror.CoreError
	ListGNPs(ctx context.Context, input *model.ListGNPsInput) ([]*entity.GlobalNetworkPolicy, *ierror.CoreError)
	UpsertGNS(ctx context.Context, gns *entity.GlobalNetworkSet) *ierror.CoreError
	GetGNSByName(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.CoreError)
	DeleteGNSByName(ctx context.Context, name string) *ierror.CoreError
	ListGNSs(ctx context.Context) ([]*entity.GlobalNetworkSet, *ierror.CoreError)
}
