package usecase

import (
	"context"
	"time"

	"github.com/bamboo-firewall/be/domain"
)

type statisticUsecase struct {
	policyRepository domain.PolicyRepository
	gnsRepository    domain.GNSRepository
	hepRepository    domain.HEPRepository
	userRepository   domain.UserRepository
	contextTimeout   time.Duration
}

// GetSummary implements domain.StatisticUsecase.
func (su *statisticUsecase) GetSummary(c context.Context) (domain.Summary, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	totalGns, err := su.gnsRepository.GetTotal(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	totalHep, err := su.hepRepository.GetTotal(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	totalPolicy, err := su.policyRepository.GetTotal(ctx)
	if err != nil {
		return domain.Summary{}, err
	}
	totalUser, err := su.userRepository.GetTotal(ctx)
	if err != nil {
		return domain.Summary{}, err
	}

	return domain.Summary{
		TotalGlobalNetworkSet: totalGns,
		TotalHostEndpoint:     totalHep,
		TotalPolicy:           totalPolicy,
		TotalUser:             totalUser,
	}, err
}

func (su *statisticUsecase) GetProjectSummary(c context.Context) ([]domain.ProjectSummary, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	projects, err := su.hepRepository.GetProjectSummary(ctx)
	if err != nil {
		return nil, err
	}
	return projects, nil
}

func NewStatisticUsecase(
	policyRepository domain.PolicyRepository,
	gnsRepository domain.GNSRepository,
	hepRepository domain.HEPRepository,
	userRepository domain.UserRepository,
	timeout time.Duration,
) domain.StatisticUsecase {
	return &statisticUsecase{
		policyRepository: policyRepository,
		gnsRepository:    gnsRepository,
		hepRepository:    hepRepository,
		userRepository:   userRepository,
		contextTimeout:   timeout,
	}
}
