package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
)

type Movie interface {
	Create(ctx context.Context, movie *service_models.Movie) error
	Get(ctx context.Context, id int64) (*service_models.Movie, error)
	Update(ctx context.Context, movie *service_models.Movie) error
	Delete(ctx context.Context, id int64) error
}

type movieRepository struct {
	db *sql.DB
}

func (m *movieRepository) Create(ctx context.Context, movie *service_models.Movie) error {
	return nil
}

func (m *movieRepository) Get(ctx context.Context, id int64) (*service_models.Movie, error) {
	return nil, nil
}

func (m *movieRepository) Update(ctx context.Context, movie *service_models.Movie) error {
	return nil
}

func (m *movieRepository) Delete(ctx context.Context, id int64) error {
	return nil
}

func NewMovieRepository(db *sql.DB) Movie {
	return &movieRepository{
		db: db,
	}
}
