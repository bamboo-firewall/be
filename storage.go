package be

import (
	"context"

	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
)

type Storage interface {
	UpsertHostEndpoint(ctx context.Context, hep *entity.HostEndpoint) *ierror.CoreError
	GetHostEndpointByName(ctx context.Context, name string) (*entity.HostEndpoint, *ierror.CoreError)
	DeleteHostEndpointByName(ctx context.Context, name string) *ierror.CoreError
	UpsertGroupPolicy(ctx context.Context, gnp *entity.GlobalNetworkPolicy) *ierror.CoreError
	GetGNPByName(ctx context.Context, name string) (*entity.GlobalNetworkPolicy, *ierror.CoreError)
	DeleteGNPByName(ctx context.Context, name string) *ierror.CoreError
	ListGNP(ctx context.Context) ([]*entity.GlobalNetworkPolicy, *ierror.CoreError)
	UpsertGNS(ctx context.Context, gns *entity.GlobalNetworkSet) *ierror.CoreError
	GetGNSByName(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.CoreError)
	DeleteGNSByName(ctx context.Context, name string) *ierror.CoreError
	ListGNS(ctx context.Context) ([]*entity.GlobalNetworkSet, *ierror.CoreError)
}
