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

func (r *PolicyMongo) UpsertGNS(ctx context.Context, gns *entity.GlobalNetworkSet) *ierror.CoreError {
	filter := bson.D{{Key: "_id", Value: gns.ID}}
	update := bson.D{{Key: "$set", Value: gns}}
	opts := options.Update().SetUpsert(true)
	_, err := r.mongo.Database.Collection(gns.CollectionName()).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}

	return nil
}

func (r *PolicyMongo) GetGNSByName(ctx context.Context, name string) (*entity.GlobalNetworkSet, *ierror.CoreError) {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	gns := new(entity.GlobalNetworkSet)
	err := r.mongo.Database.Collection(gns.CollectionName()).FindOne(ctx, filter).Decode(gns)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundHostEndpoint
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return gns, nil
}

func (r *PolicyMongo) DeleteGNSByName(ctx context.Context, name string) *ierror.CoreError {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	_, err := r.mongo.Database.Collection(entity.GlobalNetworkPolicy{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}
	return nil
}
