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
)

func (r *PolicyDB) UpsertGNS(ctx context.Context, gns *entity.GlobalNetworkSet) (*entity.GlobalNetworkSet, *ierror.CoreError) {
	session, err := r.mongo.Database.Client().StartSession()
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	defer session.EndSession(ctx)

	sessionCallback := func(sessionCtx mongo.SessionContext) (interface{}, error) {
		filter := bson.D{{Key: "metadata.name", Value: gns.Metadata.Name}}
		existedGNS := new(entity.GlobalNetworkSet)
		err = r.mongo.Database.Collection(gns.CollectionName()).FindOne(ctx, filter).Decode(existedGNS)
		if err != nil && !errors.Is(mongo.ErrNoDocuments, err) {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("find global network set failed: %w", err))
		}

		// gns is existed
		if !errors.Is(mongo.ErrNoDocuments, err) {
			gns.ID = existedGNS.ID
			gns.UUID = existedGNS.UUID
			gns.Version = existedGNS.Version
			gns.CreatedAt = existedGNS.CreatedAt
		}

		filter = bson.D{{Key: "_id", Value: gns.ID}}
		update := bson.D{{Key: "$set", Value: gns}}
		opts := options.Update().SetUpsert(true)
		_, err = r.mongo.Database.Collection(gns.CollectionName()).UpdateOne(ctx, filter, update, opts)
		if err != nil {
			if mongo.IsDuplicateKeyError(err) {
				return nil, errlist.ErrDuplicateGlobalNetworkSet.
					WithChild(fmt.Errorf("global network set already exists: %w", err))
			}
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update gns failed: %w", err))
		}

		updateVersion := bson.M{
			"$inc": bson.M{
				"version": 1,
			},
		}
		optUpdateVersions := options.FindOneAndUpdate().SetReturnDocument(options.After)
		err = r.mongo.Database.Collection(gns.CollectionName()).FindOneAndUpdate(ctx, filter, updateVersion, optUpdateVersions).Decode(gns)
		if err != nil {
			return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("update version gns failed: %w", err))
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

	return gns, nil
}

func (r *PolicyDB) GetGNSByName(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.CoreError) {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	gns := new(entity.GlobalNetworkSet)
	err := r.mongo.Database.Collection(gns.CollectionName()).FindOne(ctx, filter).Decode(gns)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundGlobalNetworkSet
		}
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("find global network set failed: %w", err))
	}
	return gns, nil
}

func (r *PolicyDB) DeleteGNSByName(ctx context.Context, name string) *ierror.CoreError {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	_, err := r.mongo.Database.Collection(entity.GlobalNetworkSet{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(fmt.Errorf("delete global network set failed: %w", err))
	}
	return nil
}

func (r *PolicyDB) ListGNSs(ctx context.Context) ([]*entity.GlobalNetworkSet, *ierror.CoreError) {
	sets := make([]*entity.GlobalNetworkSet, 0)
	cursor, err := r.mongo.Database.Collection(entity.GlobalNetworkSet{}.CollectionName()).Find(ctx, bson.D{})
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(fmt.Errorf("list global network sets failed: %w", err))
	}
	if err = cursor.All(ctx, &sets); err != nil {
		return nil, errlist.ErrUnmarshalFailed.WithChild(fmt.Errorf("decode global network sets failed: %w", err))
	}
	return sets, nil
}
