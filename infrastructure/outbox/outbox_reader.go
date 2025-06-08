package outbox

import (
	"context"
	"errors"

	"github.com/arvinpaundra/cent/user/domain/outbox/constant"
	"github.com/arvinpaundra/cent/user/domain/outbox/entity"
	"github.com/arvinpaundra/cent/user/domain/outbox/repository"
	"github.com/arvinpaundra/cent/user/model"
	"gorm.io/gorm"
)

var _ repository.OutboxReader = (*OutboxReaderRepository)(nil)

type OutboxReaderRepository struct {
	db *gorm.DB
}

func NewOutboxReaderRepository(db *gorm.DB) OutboxReaderRepository {
	return OutboxReaderRepository{db: db}
}

func (r OutboxReaderRepository) FindUnprocessed(ctx context.Context) (*entity.Outbox, error) {
	var outbox entity.Outbox

	err := r.db.WithContext(ctx).
		Model(&model.Outbox{}).
		Where("status = ? AND published_at IS NULL", constant.OutboxStatusPending.String()).
		First(&outbox).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrOutboxNotFound
		}

		return nil, err
	}

	return &outbox, nil
}
