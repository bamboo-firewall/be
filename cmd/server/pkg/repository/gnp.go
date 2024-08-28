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

func (r *PolicyMongo) UpsertGroupPolicy(ctx context.Context, gnp *entity.GlobalNetworkPolicy) *ierror.CoreError {
	filter := bson.D{{Key: "_id", Value: gnp.ID}}
	update := bson.D{{Key: "$set", Value: gnp}}
	opts := options.Update().SetUpsert(true)
	_, err := r.mongo.Database.Collection(gnp.CollectionName()).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}

	return nil
}

func (r *PolicyMongo) GetGNPByName(ctx context.Context, name string) (*entity.GlobalNetworkPolicy, *ierror.CoreError) {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	gnp := new(entity.GlobalNetworkPolicy)
	err := r.mongo.Database.Collection(gnp.CollectionName()).FindOne(ctx, filter).Decode(gnp)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, errlist.ErrNotFoundHostEndpoint
		}
		return nil, errlist.ErrDatabase.WithChild(err)
	}
	return gnp, nil
}

func (r *PolicyMongo) DeleteGNPByName(ctx context.Context, name string) *ierror.CoreError {
	filter := bson.D{{Key: "metadata.name", Value: name}}

	_, err := r.mongo.Database.Collection(entity.GlobalNetworkPolicy{}.CollectionName()).DeleteOne(ctx, filter)
	if err != nil {
		return errlist.ErrDatabase.WithChild(err)
	}
	return nil
}
