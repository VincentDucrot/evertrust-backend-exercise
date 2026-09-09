package local

import (
	"context"
	"errors"
	"evertrust-backend-exercise/entity"
	"time"
)

var tokenRepository = make(map[[32]byte]entity.Token)

var errorTokenAlreadyExist = errors.New("token already exist")
var errorTokenNotFound = errors.New("token not found")

type TokenRepository struct {
}

func NewTokenRepository() *TokenRepository {
	return &TokenRepository{}
}

func (repo *TokenRepository) Create(_ context.Context, hash [32]byte, expiresAt time.Time) error {
	if _, ok := tokenRepository[hash]; ok {
		return errorTokenAlreadyExist
	}
	tokenRepository[hash] = entity.Token{
		Hash:      hash[:],
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}
	return nil
}

func (repo *TokenRepository) Get(_ context.Context, hash [32]byte) (*entity.Token, error) {
	token, ok := tokenRepository[hash]
	if !ok {
		return nil, errorTokenNotFound
	}
	return &token, nil
}

func (repo *TokenRepository) Delete(_ context.Context, hash [32]byte) error {
	delete(tokenRepository, hash)
	return nil
}
