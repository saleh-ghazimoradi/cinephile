package repository

import (
	"context"
	"database/sql"
	"github.com/saleh-ghazimoradi/cinephile/internal/service/service_models"
	"time"
)

type Token interface {
	Insert(ctx context.Context, token *service_models.Token) error
	DeleteAllForUser(scope string, userID int64) error
	New(userID int64, ttl time.Duration, scope string) (*service_models.Token, error)
}

type tokenRepo struct {
	db *sql.DB
}

func (t *tokenRepo) Insert(ctx context.Context, token *service_models.Token) error {
	query := `
        INSERT INTO tokens (hash, user_id, expiry, scope) 
        VALUES ($1, $2, $3, $4)`

	args := []any{token.Hash, token.UserID, token.Expiry, token.Scope}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, args...)
	return err

}

func (t *tokenRepo) DeleteAllForUser(scope string, userID int64) error {
	query := `
        DELETE FROM tokens 
        WHERE scope = $1 AND user_id = $2`

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	_, err := t.db.ExecContext(ctx, query, scope, userID)
	return err
}

func (t *tokenRepo) New(userID int64, ttl time.Duration, scope string) (*service_models.Token, error) {
	token, err := service_models.GenerateToken(userID, ttl, scope)
	if err != nil {
		return nil, err
	}

	if err = t.Insert(context.Background(), token); nil != err {
		return nil, err
	}

	return token, nil
}

func NewTokenRepository(db *sql.DB) Token {
	return &tokenRepo{
		db: db,
	}
}
