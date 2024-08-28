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

func NewHEP(policyMongo *repository.PolicyMongo) *HEP {
	return &HEP{
		storage: policyMongo,
	}
}

type HEP struct {
	storage be.Storage
}

func (ds *HEP) CreateHEP(ctx context.Context, input *model.CreateHostEndpointInput) (*entity.HostEndpoint, *ierror.Error) {
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

	hep := &entity.HostEndpoint{
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
		hep.ID = hepExisted.ID
		hep.UUID = hepExisted.UUID
		hep.Version = hepExisted.Version + 1
		hep.CreatedAt = hepExisted.CreatedAt
	}

	if coreErr = ds.storage.UpsertHostEndpoint(ctx, hep); coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "Create host endpoint failed").SetSubError(coreErr)
	}
	return hep, nil
}

func (ds *HEP) DeleteHEP(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteHostEndpointByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "Delete host endpoint failed").SetSubError(coreErr)
	}
	return nil
}
