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

func (r *PolicyDB) UpsertGroupPolicy(ctx context.Context, gnp *entity.GlobalNetworkPolicy) (*entity.GlobalNetworkPolicy, *ierror.CoreError) {
	session, err := r.mongo.Database.Client().StartSession()
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	defer session.EndSession(ctx)

	sessionCallback := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.D{{Key: "metadata.name", Value: gnp.Metadata.Name}}
		existedGNP := new(entity.GlobalNetworkPolicy)
		err = r.mongo.Database.Collection(gnp.CollectionName()).FindOne(ctx, filter).Decode(existedGNP)
		if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("find global network policy failed: %w", err))
		}

		// gnp is existed
		if !errors.Is(mongo.ErrNoDocuments, err) {
			gnp.ID = existedGNP.ID
			gnp.UUID = existedGNP.UUID
			gnp.Version = existedGNP.Version
			gnp.CreatedAt = existedGNP.CreatedAt
		}

		filter = bson.D{{Key: "_id", Value: gnp.ID}}
		update := bson.D{{Key: "$set", Value: gnp}}
		opts := options.Update().SetUpsert(true)
		_, err = r.mongo.Database.Collection(gnp.CollectionName()).UpdateOne(ctx, filter, update, opts)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, errlist.ErrDuplicateGlobalNetworkPolicy.WithChild(fmt.Errorf("global network policy already exists: %w", err))
			}
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update gnp failed: %w", err))
		}

		updateVersion := bson.M{
			"$inc": bson.M{
				"version": 1,
			},
		}
		optUpdateVersions := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = r.mongo.Database.Collection(gnp.CollectionName()).FindOneAndUpdate(ctx, filter, updateVersion, optUpdateVersions).Decode(gnp)
		if err != nil {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update version gnp failed: %w", err))
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

	return gnp, nil
}

func (r *PolicyDB) GetGNPByName(ctx context.Context, name string) (*entity.GlobalNetworkPolicy, *ierror.CoreError) {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	gnp := new(entity.GlobalNetworkPolicy)
	err := r.mongo.Database.Collection(gnp.CollectionName()).FindOne(ctx, filter).Decode(gnp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundGlobalNetworkPolicy
		}
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("find global network policy failed: %w", err))
	}
	return gnp, nil
}

func (r *PolicyDB) DeleteGNPByName(ctx context.Context, name string) *ierror.CoreError {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	_, err := r.mongo.Database.Collection(entity.GlobalNetworkPolicy{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(fmt.Errorf("delete global network policy failed: %w", err))
	}
	return nil
}

func (r *PolicyDB) ListGNPs(ctx context.Context, input *model.ListGNPsInput) ([]*entity.GlobalNetworkPolicy, *ierror.CoreError) {
	var opts []*options.FindOptions
	if input != nil {
		if input.IsOrder {
			opts = append(opts, options.Find().SetSort(bson.D{{"spec.order", 1}}))
		}
	}
	policies := make([]*entity.GlobalNetworkPolicy, 0)
	cursor, err := r.mongo.Database.Collection(entity.GlobalNetworkPolicy{}.CollectionName()).Find(ctx, bson.D{}, opts...)
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("list global network policies failed: %w", err))
	}
	if err = cursor.All(ctx, &policies); err != nil {
		return nil, errlist.ErrUnmarshalFailed.WithChild(fmt.Errorf("decode global network policies failed: %w", err))
	}
	return policies, nil
}
