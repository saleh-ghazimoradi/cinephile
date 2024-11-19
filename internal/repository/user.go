package repository

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"errors"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"time"
)

type User interface {
	Insert(ctx context.Context, user *service_models.User) error
	GetByEmail(ctx context.Context, email string) (*service_models.User, error)
	Update(ctx context.Context, user *service_models.User) error
	GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*service_models.User, error)
}

type userRepository struct {
	db *sql.DB
}

func (u *userRepository) Insert(ctx context.Context, user *service_models.User) error {
	query := `INSERT INTO users (name, email, password_hash, activated) 
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, version`

	args := []any{user.Name, user.Email, user.Password.Hash, user.Activated}

	err := u.db.QueryRowContext(ctx, query, args...).Scan(&user.ID, &user.CreatedAt, &user.Version)

	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetByEmail(ctx context.Context, email string) (*service_models.User, error) {
	query := `
        SELECT id, created_at, name, email, password_hash, activated, version
        FROM users
        WHERE email = $1`

	var user service_models.User

	err := u.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.CreatedAt, &user.Name, &user.Email, &user.Password.Hash, &user.Activated, &user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}
	return &user, nil
}

func (u *userRepository) Update(ctx context.Context, user *service_models.User) error {
	query := `
        UPDATE users 
        SET name = $1, email = $2, password_hash = $3, activated = $4, version = version + 1
        WHERE id = $5 AND version = $6
        RETURNING version`

	args := []any{
		user.Name,
		user.Email,
		user.Password.Hash,
		user.Activated,
		user.ID,
		user.Version,
	}

	err := u.db.QueryRowContext(ctx, query, args...).Scan(&user.Version)
	if err != nil {
		switch {
		case err.Error() == `pq: duplicate key value violates unique constraint "users_email_key"`:
			return ErrDuplicateEmail
		case errors.Is(err, sql.ErrNoRows):
			return ErrEditConflict
		default:
			return err
		}
	}
	return nil
}

func (u *userRepository) GetForToken(ctx context.Context, tokenScope, tokenPlaintext string) (*service_models.User, error) {

	tokenHash := sha256.Sum256([]byte(tokenPlaintext))

	query := `
        SELECT users.id, users.created_at, users.name, users.email, users.password_hash, users.activated, users.version
        FROM users
        INNER JOIN tokens
        ON users.id = tokens.user_id
        WHERE tokens.hash = $1
        AND tokens.scope = $2 
        AND tokens.expiry > $3`

	args := []any{tokenHash[:], tokenScope, time.Now()}

	var user service_models.User

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err := u.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.Name,
		&user.Email,
		&user.Password.Hash,
		&user.Activated,
		&user.Version,
	)
	if err != nil {
		switch {
		case errors.Is(err, sql.ErrNoRows):
			return nil, ErrRecordNotFound
		default:
			return nil, err
		}
	}

	return &user, nil
}

func NewUserRepository(db *sql.DB) User {
	return &userRepository{
		db: db,
	}
}
