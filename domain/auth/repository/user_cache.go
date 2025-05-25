package repository

import (
	"context"
	"time"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
)

type UserCache interface {
	Get(ctx context.Context, key string) (entity.User, error)
	Set(ctx context.Context, key string, value entity.User, ttl time.Duration) error
	Del(ctx context.Context, key string) error
}
