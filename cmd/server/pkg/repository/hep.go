package repository

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/bamboo-firewall/be/cmd/server/pkg/common/errlist"
	"github.com/bamboo-firewall/be/cmd/server/pkg/entity"
	"github.com/bamboo-firewall/be/cmd/server/pkg/httpbase/ierror"
)

func (r *PolicyDB) UpsertHostEndpoint(ctx context.Context, hep *entity.HostEndpoint) *ierror.CoreError {
	filter := bson.D{{Key: "_id", Value: hep.ID}}
	update := bson.D{{Key: "$set", Value: hep}}
	opts := options.Update().SetUpsert(true)
	_, err := r.mongo.Database.Collection(hep.CollectionName()).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}

	return nil
}

func (r *PolicyDB) GetHostEndpointByName(ctx context.Context, name string) (*entity.HostEndpoint, *ierror.CoreError) {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	hep := new(entity.HostEndpoint)
	err := r.mongo.Database.Collection(hep.CollectionName()).FindOne(ctx, filter).Decode(hep)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundHostEndpoint
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return hep, nil
}

func (r *PolicyDB) DeleteHostEndpointByName(ctx context.Context, name string) *ierror.CoreError {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	_, err := r.mongo.Database.Collection(entity.HostEndpoint{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}
	return nil
}

func (r *PolicyDB) ListHostEndpoints(ctx context.Context) ([]*entity.HostEndpoint, *ierror.CoreError) {
	heps := make([]*entity.HostEndpoint, 0)
	cursor, err := r.mongo.Database.Collection(entity.HostEndpoint{}.CollectionName()).Find(ctx, bson.D{})
	if err != nil {
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	if err = cursor.All(ctx, &heps); err != nil {
		return nil, errlist.ErrUnmarshalFailed.WithChild(err)
	}
	return heps, nil
}
