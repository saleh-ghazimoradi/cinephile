package service

import (
	"context"
	"github.com/saleh-ghazimoradi/cinephile/internal/repository"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
)

type Movie interface {
	Create(ctx context.Context, movie *service_models.Movie) error
	Get(ctx context.Context, id int64) (*service_models.Movie, error)
	GetAll(ctx context.Context, fq service_models.Filter) ([]*service_models.Movie, error)
	Update(ctx context.Context, movie *service_models.Movie) error
	Delete(ctx context.Context, id int64) error
}

type movieService struct {
	movieRepo repository.Movie
}

func (s *movieService) Create(ctx context.Context, movie *service_models.Movie) error {
	return s.movieRepo.Create(ctx, movie)
}

func (s *movieService) Get(ctx context.Context, id int64) (*service_models.Movie, error) {
	return s.movieRepo.Get(ctx, id)
}

func (s *movieService) Update(ctx context.Context, movie *service_models.Movie) error {
	return s.movieRepo.Update(ctx, movie)
}

func (s *movieService) Delete(ctx context.Context, id int64) error {
	return s.movieRepo.Delete(ctx, id)
}

func (s *movieService) GetAll(ctx context.Context, fq service_models.Filter) ([]*service_models.Movie, error) {
	return s.movieRepo.GetAll(ctx, fq)
}

func NewMovieService(movieRepo repository.Movie) Movie {
	return &movieService{
		movieRepo: movieRepo,
	}
}
