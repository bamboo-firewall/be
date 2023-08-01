package usecase

import (
	"context"
	"time"

	models "github.com/bamboo-firewall/watcher/model"

	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/internal/optionutil"
)

type gnsUsecase struct {
	gnsRepository  domain.GNSRepository
	contextTimeout time.Duration
}

var GNSMapping = map[string]string{
	"name": "$metadata.name",
	"zone": "$metadata.labels.zone",
}

func (gu *gnsUsecase) Fetch(c context.Context) ([]models.GlobalNetworkSet, error) {
	ctx, cancel := context.WithTimeout(c, gu.contextTimeout)
	defer cancel()
	return gu.gnsRepository.Fetch(ctx)
}

func (u *gnsUsecase) Search(c context.Context, options []domain.Option) ([]models.GlobalNetworkSet, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	findOptions := optionutil.ConvertToBsonM(options, GNSMapping)
	return u.gnsRepository.Search(ctx, findOptions)
}

func (u *gnsUsecase) GetOptions(c context.Context, filter []domain.Option, key string) ([]domain.Option, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	query := optionutil.ConvertToBsonM(filter, GNSMapping)
	return u.gnsRepository.AggGroupBy(ctx, query, key, GNSMapping[key])
}

func NewGNSUsecase(gnsRepository domain.GNSRepository, timeout time.Duration) domain.GNSUsecase {
	return &gnsUsecase{
		gnsRepository:  gnsRepository,
		contextTimeout: timeout,
	}
}
