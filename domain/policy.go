package domain

import (
	"context"

	models "github.com/bamboo-firewall/watcher/model"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionPolicy = "globalnetworkpolicies"
)

type PolicyRepository interface {
	Fetch(c context.Context) ([]models.GlobalNetworkPolicies, error)
	Search(c context.Context, options bson.M) ([]models.GlobalNetworkPolicies, error)
	GetTotal(c context.Context) (int64, error)
	AggGroupBy(c context.Context, query bson.M, key string, jsonPath string) ([]Option, error)
}

type PolicyUsecase interface {
	Fetch(c context.Context) ([]models.GlobalNetworkPolicies, error)
	Search(c context.Context, options []Option) ([]models.GlobalNetworkPolicies, error)
	GetOptions(c context.Context, filter []Option, key string) ([]Option, error)
}
