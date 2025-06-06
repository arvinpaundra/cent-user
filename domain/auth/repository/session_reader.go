package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
)

type SessionReader interface {
	FindByRefreshToken(ctx context.Context, userId int64, refreshToken string) (*entity.Session, error)
}
