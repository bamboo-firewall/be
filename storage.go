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
}
