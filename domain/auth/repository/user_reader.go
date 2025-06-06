package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
)

type UserReader interface {
	FindByEmail(ctx context.Context, email string) (*entity.User, error)
	FindById(ctx context.Context, id int64) (*entity.User, error)
	IsEmailExist(ctx context.Context, email string) (bool, error)
}
