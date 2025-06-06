package auth

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/redis/go-redis/v9"
)

var _ repository.UserCache = (*UserCacheRepository)(nil)

type UserCacheRepository struct {
	rdb *redis.Client
}

func NewUserCacheRepository(rdb *redis.Client) UserCacheRepository {
	return UserCacheRepository{rdb: rdb}
}

func (r UserCacheRepository) Get(ctx context.Context, key string) (*entity.User, error) {
	var user entity.User

	err := r.rdb.Get(ctx, key).Scan(&user)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			return nil, constant.ErrUserNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r UserCacheRepository) Set(ctx context.Context, key string, value *entity.User, ttl time.Duration) error {
	valb, _ := json.Marshal(value)

	err := r.rdb.Set(ctx, key, valb, ttl).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r UserCacheRepository) Del(ctx context.Context, key string) error {
	panic("not implemented")
}
