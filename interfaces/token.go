package interfaces

import (
	"context"
	"evertrust-backend-exercise/entity"
	"time"
)

type TokenRepository interface {
	Create(ctx context.Context, hash [32]byte, expiresAt time.Time) error
	Get(ctx context.Context, hash [32]byte) (*entity.Token, error)
	Delete(ctx context.Context, hash [32]byte) error
}
