package mysql

import (
	"context"
	"evertrust-backend-exercise/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type TokenRepo struct {
	db *sqlx.DB
}

func NewTokenRepo(db *sqlx.DB) *TokenRepo {
	return &TokenRepo{
		db: db,
	}
}

func (repo *TokenRepo) Create(ctx context.Context, hash [32]byte, expiresAt time.Time) error {
	_, err := repo.db.ExecContext(ctx, `INSERT INTO token(hash, expires_at) VALUES (?, ?)`, hash[:], expiresAt)
	return err
}

func (repo *TokenRepo) Get(ctx context.Context, hash [32]byte) (*entity.Token, error) {
	var token entity.Token
	err := repo.db.GetContext(ctx, &token, "SELECT hash, expires_at, created_at FROM token WHERE hash = ?", hash[:])
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (repo *TokenRepo) Delete(ctx context.Context, hash [32]byte) error {
	_, err := repo.db.ExecContext(ctx, `DELETE FROM token WHERE hash=?`, hash[:])
	return err
}
