package storage

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

const (
	mongoDefaultMaxPoolSize     = 100
	mongoDefaultMaxConnIdleTime = 30 * time.Minute

	mongoSocketTimeout = 3 * time.Second
)

type ConfigMongo struct {
	AuthMechanism   string
	AuthSource      string
	Username        string
	Password        string
	URI             string
	SocketTimeout   time.Duration
	MaxPoolSize     uint64
	MaxConnIdleTime time.Duration
}

type PolicyMongo struct {
	Database *mongo.Database
}

func NewPolicyMongo(config ConfigMongo) (*PolicyMongo, error) {
	credential := options.Credential{
		AuthMechanism: config.AuthMechanism,
		AuthSource:    config.AuthSource,
		Username:      config.Username,
		Password:      config.Password,
	}
	opts := options.Client()
	opts.ApplyURI(config.URI)
	opts.SetAuth(credential)
	opts.SetSocketTimeout(mongoSocketTimeout)

	// config connection pool
	if config.MaxPoolSize == 0 {
		config.MaxPoolSize = mongoDefaultMaxPoolSize
	}
	if config.MaxConnIdleTime <= 0 {
		config.MaxConnIdleTime = mongoDefaultMaxConnIdleTime
	}
	opts.SetMaxPoolSize(config.MaxPoolSize)
	opts.SetMaxConnIdleTime(config.MaxConnIdleTime)

	client, err := mongo.Connect(context.Background(), opts)
	if err != nil {
		return nil, err
	}
	if err = client.Ping(context.Background(), readpref.Primary()); err != nil {
		return nil, err
	}
	return &PolicyMongo{Database: client.Database("policy")}, nil
}

func (pm *PolicyMongo) Stop(ctx context.Context) error {
	slog.Info("Stop policy mongo")
	return pm.Database.Client().Disconnect(ctx)
}
