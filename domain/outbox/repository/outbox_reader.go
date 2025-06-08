package repository

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/outbox/entity"
)

type OutboxReader interface {
	FindUnprocessed(ctx context.Context) (*entity.Outbox, error)
}
