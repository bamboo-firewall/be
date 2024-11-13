package repository

import (
	"context"
	"errors"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"

	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
	"github.com/bamboo-firewall/be/domain/model"
)

func (r *PolicyDB) UpsertHostEndpoint(ctx context.Context, hep *entity.HostEndpoint) (*entity.HostEndpoint, *ierror.CoreError) {
	session, err := r.mongo.Database.Client().StartSession()
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	defer session.EndSession(ctx)

	sessionCallback := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.D{{Key: "spec.tenant_id", Value: hep.Spec.TenantID}, {Key: "spec.ip", Value: hep.Spec.IP}}
		existedHEP := new(entity.HostEndpoint)
		err = r.mongo.Database.Collection(hep.CollectionName()).FindOne(ctx, filter).Decode(existedHEP)
		if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("find host endpoint failed: %w", err))
		}

		// hep is not exist
		if !errors.Is(mongo.ErrNoDocuments, err) {
			hep.ID = existedHEP.ID
			hep.UUID = existedHEP.UUID
			hep.Version = existedHEP.Version
			hep.CreatedAt = existedHEP.CreatedAt
		}

		filter = bson.D{{Key: "_id", Value: hep.ID}}
		update := bson.D{{Key: "$set", Value: hep}}
		opts := options.Update().SetUpsert(true)
		_, err = r.mongo.Database.Collection(hep.CollectionName()).UpdateOne(ctx, filter, update, opts)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, errlist.ErrDuplicateHostEndpoint.WithChild(fmt.Errorf("host endpoint already exists: %w", err))
			}
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update host endpoint failed: %w", err))
		}

		updateVersion := bson.M{
			"$inc": bson.M{
				"version": 1,
			},
		}
		optUpdateVersions := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = r.mongo.Database.Collection(hep.CollectionName()).FindOneAndUpdate(ctx, filter, updateVersion, optUpdateVersions).Decode(hep)
		if err != nil {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update version host endpoint failed: %w", err))
		}

		return nil, nil
	}

	opts := options.Transaction().SetWriteConcern(writeconcern.Majority()).SetReadConcern(readconcern.Snapshot())
	_, sessionErr := session.WithTransaction(ctx, sessionCallback, opts)
	if sessionErr != nil {
		var coreErr *ierror.CoreError
		if errors.As(sessionErr, &coreErr) {
			return nil, coreErr
		}
		return nil, errlist.ErrDatabase.WithChild(sessionErr)
	}

	return hep, nil
}

func (r *PolicyDB) GetHostEndpoint(ctx context.Context, input *model.GetHostEndpointInput) (*entity.HostEndpoint, *ierror.CoreError) {
	var filter bson.D
	if input != nil {
		filter = bson.D{
			{Key: "spec.tenant_id", Value: input.TenantID},
			{Key: "spec.ip", Value: input.IP},
		}
	}

	hep := new(entity.HostEndpoint)
	err := r.mongo.Database.Collection(hep.CollectionName()).FindOne(ctx, filter).Decode(hep)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundHostEndpoint
		}
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("get host endpoint failed: %w", err))
	}
	return hep, nil
}

func (r *PolicyDB) DeleteHostEndpoint(ctx context.Context, tenantID uint64, ip uint32) *ierror.CoreError {
	filter := bson.D{{Key: "spec.tenant_id", Value: tenantID}, {Key: "spec.ip", Value: ip}}

	_, err := r.mongo.Database.Collection(entity.HostEndpoint{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(fmt.Errorf("delete host endpoint failed: %w", err))
	}
	return nil
}

func (r *PolicyDB) ListHostEndpoints(ctx context.Context, input *model.ListHostEndpointsInput) ([]*entity.HostEndpoint, *ierror.CoreError) {
	filter := bson.D{}
	if input != nil {
		if input.TenantID != nil {
			filter = append(filter, bson.E{Key: "spec.tenant_id", Value: *input.TenantID})
		}
		if input.IP != nil {
			filter = append(filter, bson.E{Key: "spec.ip", Value: *input.IP})
		}
	}

	heps := make([]*entity.HostEndpoint, 0)
	cursor, err := r.mongo.Database.Collection(entity.HostEndpoint{}.CollectionName()).Find(ctx, filter)
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("list host endpoints failed: %w", err))
	}
	if err = cursor.All(ctx, &heps); err != nil {
		return nil, errlist.ErrUnmarshalFailed.WithChild(fmt.Errorf("decode host endpoints failed: %w", err))
	}
	return heps, nil
}
