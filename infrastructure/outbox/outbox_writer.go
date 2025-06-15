package outbox

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/outbox/entity"
	"github.com/arvinpaundra/cent/user/domain/outbox/repository"
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
	} else if outbox.IsMarkedToBeUpdated() {
		return r.update(ctx, outbox)
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

func (r OutboxWriterRepository) update(ctx context.Context, outbox *entity.Outbox) error {
	outboxModel := model.Outbox{
		Status:      outbox.Status.String(),
		Event:       outbox.Event.String(),
		PublishedAt: null.TimeFromPtr(outbox.PublishedAt),
		Error:       null.StringFromPtr(outbox.Error),
		Payload:     outbox.Payload,
	}

	err := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("id = ?", outbox.ID).
		Updates(&outboxModel).
		Error

	if err != nil {
		return err
	}

	return nil
}
