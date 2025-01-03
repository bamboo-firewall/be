package storage

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"

	entity2 "github.com/bamboo-firewall/be/pkg/entity"
)

type PolicyDB struct {
	Database *mongo.Database
}

func NewPolicyDB(uri string) (*PolicyDB, error) {
	opts := options.Client()
	opts.ApplyURI(uri)
	cs, cErr := connstring.ParseAndValidate(uri)
	if cErr != nil {
		return nil, cErr
	}
	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}
	pm := &PolicyDB{
		Database: client.Database(cs.Database),
	}
	if err = pm.createIndexes(); err != nil {
		return nil, err
	}
	return pm, nil
}

func (pm *PolicyDB) createIndexes() error {
	indexMap := map[string][]mongo.IndexModel{
		entity2.HostEndpoint{}.CollectionName(): {
			{
				Keys:    bson.D{{Key: "spec.tenant_id", Value: 1}, {Key: "spec.ip", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "uuid", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		entity2.GlobalNetworkSet{}.CollectionName(): {
			{
				Keys:    bson.D{{Key: "metadata.name", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys:    bson.D{{Key: "uuid", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
		entity2.GlobalNetworkPolicy{}.CollectionName(): {
			{
				Keys:    bson.D{{Key: "metadata.name", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
			{
				Keys: bson.D{{Key: "spec.order", Value: 1}},
			},
			{
				Keys:    bson.D{{Key: "uuid", Value: 1}},
				Options: options.Index().SetUnique(true),
			},
		},
	}
	for collectName, indexes := range indexMap {
		_, err := pm.Database.Collection(collectName).Indexes().CreateMany(context.TODO(), indexes)
		if err != nil {
			return fmt.Errorf("create index: %w", err)
		}
	}

	return nil
}

func (pm *PolicyDB) Stop(ctx context.Context) error {
	slog.Info("Stop policy mongo")
	return pm.Database.Client().Disconnect(ctx)
}
