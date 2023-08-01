package usecase

import (
	"context"
	"time"

	models "github.com/bamboo-firewall/watcher/model"

	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/internal/optionutil"
)

type hepUsecase struct {
	hepRepository  domain.HEPRepository
	contextTimeout time.Duration
}

var HEPMapping = map[string]string{
	"name":      "$spec.node",
	"ip":        "$spec.expectedIPs",
	"namespace": "$metadata.labels.namespace",
	"project":   "$metadata.labels.project",
	"role":      "$metadata.labels.role",
	"zone":      "$metadata.labels.zone",
}

func (hu *hepUsecase) Fetch(c context.Context) ([]models.HostEndPoint, error) {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()
	return hu.hepRepository.Fetch(ctx)
}

func (hu *hepUsecase) Search(c context.Context, options []domain.Option) ([]models.HostEndPoint, error) {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()
	query := optionutil.ConvertToBsonM(options, HEPMapping)
	return hu.hepRepository.Search(ctx, query)
}

func (hu *hepUsecase) GetOptions(c context.Context, filter []domain.Option, key string) ([]domain.Option, error) {
	ctx, cancel := context.WithTimeout(c, hu.contextTimeout)
	defer cancel()
	query := optionutil.ConvertToBsonM(filter, HEPMapping)
	return hu.hepRepository.AggGroupBy(ctx, query, key, HEPMapping[key])
}

func NewHEPUsecase(hepRepository domain.HEPRepository, timeout time.Duration) domain.HEPUsecase {
	return &hepUsecase{
		hepRepository:  hepRepository,
		contextTimeout: timeout,
	}
}
