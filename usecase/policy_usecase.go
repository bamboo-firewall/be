package usecase

import (
	"context"
	"time"

	models "github.com/bamboo-firewall/watcher/model"

	"github.com/bamboo-firewall/be/domain"
	"github.com/bamboo-firewall/be/internal/optionutil"
)

type policyUsecase struct {
	policyRepository domain.PolicyRepository
	contextTimeout   time.Duration
}

var PolicyMapping = map[string]string{
	"name": "$metadata.name",
}

func (u *policyUsecase) Search(c context.Context, options []domain.Option) ([]models.GlobalNetworkPolicies, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	findOptions := optionutil.ConvertToBsonM(options, PolicyMapping)
	return u.policyRepository.Search(ctx, findOptions)
}

func (u *policyUsecase) Fetch(c context.Context) ([]models.GlobalNetworkPolicies, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	return u.policyRepository.Fetch(ctx)
}

func (u *policyUsecase) GetOptions(c context.Context, filter []domain.Option, key string) ([]domain.Option, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()
	query := optionutil.ConvertToBsonM(filter, PolicyMapping)
	return u.policyRepository.AggGroupBy(ctx, query, key, PolicyMapping[key])
}

func NewPolicyUsecase(policyRepository domain.PolicyRepository, timeout time.Duration) domain.PolicyUsecase {
	return &policyUsecase{
		policyRepository: policyRepository,
		contextTimeout:   timeout,
	}
}
