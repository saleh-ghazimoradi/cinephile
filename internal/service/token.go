package service

import (
	"context"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"time"
)

type Token interface {
	Insert(ctx context.Context, token *service_models.Token) error
	DeleteAllForUser(scope string, userID int64) error
	New(userID int64, ttl time.Duration, scope string) (*service_models.Token, error)
}

type tokenService struct {
	tokenRepo repository.Token
}

func (t *tokenService) Insert(ctx context.Context, token *service_models.Token) error {
	return t.tokenRepo.Insert(ctx, token)
}

func (t *tokenService) DeleteAllForUser(scope string, userID int64) error {
	return t.tokenRepo.DeleteAllForUser(scope, userID)
}

func (t *tokenService) New(userID int64, ttl time.Duration, scope string) (*service_models.Token, error) {
	return t.tokenRepo.New(userID, ttl, scope)
}

func NewTokenService(tokenRepo repository.Token) Token {
	return &tokenService{
		tokenRepo: tokenRepo,
	}
}
