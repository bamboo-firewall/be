package storage

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/x/mongo/driver/connstring"
)

type PolicyMongo struct {
	Database *mongo.Database
}

func NewPolicyMongo(uri string) (*PolicyMongo, error) {
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
	return &PolicyMongo{Database: client.Database(cs.Database)}, nil
}

func (pm *PolicyMongo) Stop(ctx context.Context) error {
	slog.Info("Stop policy mongo")
	return pm.Database.Client().Disconnect(ctx)
}
