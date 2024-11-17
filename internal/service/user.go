package service

import (
	"context"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
)

type User interface {
	Insert(ctx context.Context, user *service_models.User) error
	GetByEmail(ctx context.Context, email string) (*service_models.User, error)
	Update(ctx context.Context, user *service_models.User) error
}

type userService struct {
	userRepo repository.User
}

func (u *userService) Insert(ctx context.Context, user *service_models.User) error {
	return u.userRepo.Insert(ctx, user)
}

func (u *userService) GetByEmail(ctx context.Context, email string) (*service_models.User, error) {
	return u.userRepo.GetByEmail(ctx, email)
}

func (u *userService) Update(ctx context.Context, user *service_models.User) error {
	return u.userRepo.Update(ctx, user)
}

func NewUserService(userRepo repository.User) User {
	return &userService{
		userRepo: userRepo,
	}
}
