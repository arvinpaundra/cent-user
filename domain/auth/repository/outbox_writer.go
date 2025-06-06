package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
)

type OutboxWriter interface {
	Save(ctx context.Context, outbox *entity.Outbox) error
}
