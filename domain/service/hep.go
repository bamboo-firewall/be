package service

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/bamboo-firewall/be"
	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/cmd/server/pkg/repository"
	"github.com/bamboo-firewall/be/domain/model"
)

func NewHEP(policyMongo *repository.PolicyDB) *hep {
	return &hep{
		storage: policyMongo,
	}
}

type hep struct {
	storage be.Storage
}

func (ds *hep) Create(ctx context.Context, input *model.CreateHostEndpointInput) (*entity.HostEndpoint, *ierror.Error) {
	// ToDo: use transaction and lock row
	hepExisted, coreErr := ds.storage.GetHostEndpointByName(ctx, input.Metadata.Name)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
		return nil, httpbase.ErrDatabase(ctx, "get host endpoint failed").SetSubError(coreErr)
	}

	var ports []entity.HostEndpointSpecPort
	for _, port := range input.Spec.Ports {
		ports = append(ports, entity.HostEndpointSpecPort{
			Name:     port.Name,
			Port:     port.Port,
			Protocol: port.Protocol,
		})
	}

	hepEntity := &entity.HostEndpoint{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: entity.HostEndpointMetadata{
			Name:   input.Metadata.Name,
			Labels: input.Metadata.Labels,
		},
		Spec: entity.HostEndpointSpec{
			InterfaceName: input.Spec.InterfaceName,
			IPs:           input.Spec.IPs,
			Ports:         ports,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if !errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
		hepEntity.ID = hepExisted.ID
		hepEntity.UUID = hepExisted.UUID
		hepEntity.Version = hepExisted.Version + 1
		hepEntity.CreatedAt = hepExisted.CreatedAt
	}

	if coreErr = ds.storage.UpsertHostEndpoint(ctx, hepEntity); coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "Create host endpoint failed").SetSubError(coreErr)
	}
	return hepEntity, nil
}

func (ds *hep) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteHostEndpointByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "Delete host endpoint failed").SetSubError(coreErr)
	}
	return nil
}

func (ds *hep) FetchPolicies(ctx context.Context, input *model.FetchPoliciesInput) (*entity.HostEndpoint, []*entity.GlobalNetworkPolicy, []*entity.GlobalNetworkSet, *ierror.Error) {
	hepEntity, coreErr := ds.storage.GetHostEndpointByName(ctx, input.Name)
	if coreErr != nil {
		if errors.Is(coreErr, errlist.ErrNotFoundHostEndpoint) {
			return nil, nil, nil, httpbase.ErrNotFound(ctx, "NotFound").SetSubError(coreErr)
		}
		return nil, nil, nil, httpbase.ErrDatabase(ctx, "Get host endpoint failed").SetSubError(coreErr)
	}

	policies, err := ds.storage.ListGNP(ctx)
	if err != nil {
		return nil, nil, nil, httpbase.ErrDatabase(ctx, "List policies failed").SetSubError(coreErr)
	}

	sets, err := ds.storage.ListGNS(ctx)
	if err != nil {
		return nil, nil, nil, httpbase.ErrDatabase(ctx, "List sets failed").SetSubError(coreErr)
	}

	return hepEntity, policies, sets, nil
}
