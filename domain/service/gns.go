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

func NewGNS(policyMongo *repository.PolicyMongo) *gns {
	return &gns{
		storage: policyMongo,
	}
}

type gns struct {
	storage be.Storage
}

func (ds *gns) Create(ctx context.Context, input *model.CreateGlobalNetworkSetInput) (*entity.GlobalNetworkSet, *ierror.Error) {
	// ToDo: use transaction and lock row
	gnsExisted, coreErr := ds.storage.GetGNSByName(ctx, input.Metadata.Name)
	if coreErr != nil && !errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkSet) {
		return nil, httpbase.ErrDatabase(ctx, "get global network set failed").SetSubError(coreErr)
	}

	gnsEntity := &entity.GlobalNetworkSet{
		ID:      primitive.NewObjectID(),
		UUID:    uuid.New().String(),
		Version: 1,
		Metadata: entity.GNSMetadata{
			Name: input.Metadata.Name,
		},
		Spec: entity.GNSSpec{
			Nets: input.Spec.Nets,
		},
		Description: input.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
	if !errors.Is(coreErr, errlist.ErrNotFoundGlobalNetworkSet) {
		gnsEntity.ID = gnsExisted.ID
		gnsEntity.UUID = gnsExisted.UUID
		gnsEntity.Version = gnsExisted.Version + 1
		gnsEntity.CreatedAt = gnsExisted.CreatedAt
	}

	if coreErr = ds.storage.UpsertGNS(ctx, gnsEntity); coreErr != nil {
		return nil, httpbase.ErrDatabase(ctx, "Create global network set failed").SetSubError(coreErr)
	}
	return gnsEntity, nil
}

func (ds *gns) Delete(ctx context.Context, name string) *ierror.Error {
	if coreErr := ds.storage.DeleteGNSByName(ctx, name); coreErr != nil {
		return httpbase.ErrDatabase(ctx, "Delete global network set failed").SetSubError(coreErr)
	}
	return nil
}
