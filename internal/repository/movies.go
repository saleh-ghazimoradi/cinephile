package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"time"
)

type Movie interface {
	Create(ctx context.Context, movie *service_models.Movie) error
	Get(ctx context.Context, id int64) (*service_models.Movie, error)
	GetAll(ctx context.Context, fq service_models.Filter) ([]*service_models.Movie, error)
	Update(ctx context.Context, movie *service_models.Movie) error
	Delete(ctx context.Context, id int64) error
}

type movieRepository struct {
	db *sql.DB
}

func (m *movieRepository) Create(ctx context.Context, movie *service_models.Movie) error {
	query := `
        INSERT INTO movies (title, year, runtime, genres) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	args := []any{movie.Title, movie.Year, movie.Runtime, pq.Array(movie.Genres)}

	return m.db.QueryRowContext(ctx, query, args...).Scan(&movie.ID, &movie.CreatedAt, &movie.Version)
}

func (m *movieRepository) Get(ctx context.Context, id int64) (*service_models.Movie, error) {
	query := `
        SELECT id, created_at, title, year, runtime, genres, version
        FROM movies
        WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	var movie service_models.Movie

	err := m.db.QueryRowContext(ctx, query, id).Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(&movie.Genres), &movie.Version)

	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &movie, nil
}

func (m *movieRepository) GetAll(ctx context.Context, fq service_models.Filter) ([]*service_models.Movie, error) {
	query := `SELECT id, created_at, title, year, runtime, genres, version 
              FROM movies 
              ORDER BY id ` + fq.Sort + ` 
              LIMIT $1 OFFSET $2`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	rows, err := m.db.QueryContext(ctx, query, fq.Limit, fq.Offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []*service_models.Movie
	for rows.Next() {
		var movie service_models.Movie
		err = rows.Scan(&movie.ID, &movie.CreatedAt, &movie.Title, &movie.Year, &movie.Runtime, pq.Array(movie.Genres), &movie.Version)
		if err != nil {
			return nil, err
		}
		movies = append(movies, &movie)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return movies, nil
}

func (m *movieRepository) Update(ctx context.Context, movie *service_models.Movie) error {
	query := `
        UPDATE movies 
        SET title = $1, year = $2, runtime = $3, genres = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []any{
		movie.Title,
		movie.Year,
		movie.Runtime,
		pq.Array(movie.Genres),
		movie.ID,
		movie.Version,
	}

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	if err := m.db.QueryRowContext(ctx, query, args...).Scan(&movie.Version); err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (m *movieRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM movies WHERE id = $1`

	ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
	defer cancel()

	result, err := m.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return ErrRecordNotFound
	}

	return nil
}

func NewMovieRepository(db *sql.DB) Movie {
	return &movieRepository{
		db: db,
	}
}
