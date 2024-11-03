package service

import (
	"context"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
)

type Movie interface {
	Create(ctx context.Context, movie *service_models.Movie) error
	Get(ctx context.Context, id int64) (*service_models.Movie, error)
	Update(ctx context.Context, movie *service_models.Movie) error
	Delete(ctx context.Context, id int64) error
}

type movieService struct {
	movieRepo repository.Movie
}

func (s *movieService) Create(ctx context.Context, movie *service_models.Movie) error {
	return nil
}

func (s *movieService) Get(ctx context.Context, id int64) (*service_models.Movie, error) {
	return nil, nil
}

func (s *movieService) Update(ctx context.Context, movie *service_models.Movie) error {
	return nil
}

func (s *movieService) Delete(ctx context.Context, id int64) error {
	return nil
}

func NewMovieService(movieRepo repository.Movie) Movie {
	return &movieService{
		movieRepo: movieRepo,
	}
}
