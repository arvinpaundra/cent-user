package auth

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/arvinpaundra/cent/user/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.OutboxWriter = (*OutboxWriterRepository)(nil)

type OutboxWriterRepository struct {
	db *gorm.DB
}

func NewOutboxWriterRepository(db *gorm.DB) OutboxWriterRepository {
	return OutboxWriterRepository{db: db}
}

func (r OutboxWriterRepository) Save(ctx context.Context, outbox *entity.Outbox) error {
	if outbox.IsNew() {
		return r.insert(ctx, outbox)
	}

	return nil
}

func (r OutboxWriterRepository) insert(ctx context.Context, outbox *entity.Outbox) error {
	outboxModel := model.Outbox{
		Status:      outbox.Status.String(),
		Event:       outbox.Event.String(),
		Payload:     outbox.Payload,
		PublishedAt: null.TimeFromPtr(outbox.PublishedAt),
	}

	err := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Create(&outboxModel).
		Error

	if err != nil {
		return err
	}

	return nil
}
