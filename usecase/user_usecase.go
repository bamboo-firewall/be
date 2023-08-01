package usecase

import (
	"context"
	"time"

	"github.com/bamboo-firewall/be/domain"
	"github.com/casbin/casbin/v2"
)

type userUsecase struct {
	userRepository domain.UserRepository
	contextTimeout time.Duration
}

// Update implements domain.UserUsecase.
func (su *userUsecase) Update(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.Update(ctx, user)
}

func (su *userUsecase) DeleteById(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.DeleteById(ctx, id)
}

func (su *userUsecase) Create(c context.Context, user *domain.User) error {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.Create(ctx, user)
}

func (su *userUsecase) Fetch(c context.Context) ([]domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.Fetch(ctx)
}

func (su *userUsecase) GetUserByEmail(c context.Context, email string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByEmail(ctx, email)
}

func (su *userUsecase) GetUserByID(c context.Context, id string) (domain.User, error) {
	ctx, cancel := context.WithTimeout(c, su.contextTimeout)
	defer cancel()
	return su.userRepository.GetByID(ctx, id)
}

func NewUserUsecase(userRepository domain.UserRepository, enforcer *casbin.Enforcer, timeout time.Duration) domain.UserUsecase {
	return &userUsecase{
		userRepository: userRepository,
		contextTimeout: timeout,
	}
}
