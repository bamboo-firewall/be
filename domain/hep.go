package domain

import (
	"context"

	models "github.com/bamboo-firewall/watcher/model"
	"go.mongodb.org/mongo-driver/bson"
)

const (
	CollectionHEP = "hostendpoints"
)

type HEPRepository interface {
	Fetch(c context.Context) ([]models.HostEndPoint, error)
	Search(c context.Context, options bson.M) ([]models.HostEndPoint, error)
	GetTotal(c context.Context) (int64, error)
	GetProjectSummary(c context.Context) ([]ProjectSummary, error)
	AggGroupBy(c context.Context, query bson.M, key string, jsonPath string) ([]Option, error)
}

type HEPUsecase interface {
	Fetch(c context.Context) ([]models.HostEndPoint, error)
	Search(c context.Context, options []Option) ([]models.HostEndPoint, error)
	GetOptions(c context.Context, filter []Option, key string) ([]Option, error)
}
