package domain

import (
	"context"

	models "github.com/bamboo-firewall/watcher/model"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionGNS = "globalnetworksets"
)

type GNSRepository interface {
	Fetch(c context.Context) ([]models.GlobalNetworkSet, error)
	Search(c context.Context, options bson.M) ([]models.GlobalNetworkSet, error)
	GetTotal(c context.Context) (int64, error)
	AggGroupBy(c context.Context, query bson.M, key string, jsonPath string) ([]Option, error)
}

type GNSUsecase interface {
	Fetch(c context.Context) ([]models.GlobalNetworkSet, error)
	Search(c context.Context, options []Option) ([]models.GlobalNetworkSet, error)
	GetOptions(c context.Context, filter []Option, key string) ([]Option, error)
}
