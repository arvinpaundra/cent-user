package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
)

type SessionWriter interface {
	Save(ctx context.Context, refreshToken entity.Session) error
}
