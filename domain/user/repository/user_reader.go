package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/user/entity"
)

type UserReader interface {
	FindBySlug(ctx context.Context, slug string) (*entity.User, error)
}
