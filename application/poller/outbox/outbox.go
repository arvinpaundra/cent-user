package outbox

import (
	"context"
	"errors"

	"github.com/arvinpaundra/cent/user/core/poller"
	"github.com/arvinpaundra/cent/user/domain/outbox/constant"
	"github.com/arvinpaundra/cent/user/domain/outbox/service"
	outboxinfra "github.com/arvinpaundra/cent/user/infrastructure/outbox"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

type Outbox struct {
	db *gorm.DB
	nc *nats.Conn
}

func NewOutbox(db *gorm.DB, nc *nats.Conn) Outbox {
	return Outbox{
		db: db,
		nc: nc,
	}
}

func (o Outbox) OutboxProcessor() error {
	handler := service.NewOutboxProcessorHandler(
		outboxinfra.NewOutboxReaderRepository(o.db),
		outboxinfra.NewOutboxWriterRepository(o.db),
		outboxinfra.NewUnitOfWork(o.db),
		outboxinfra.NewMessaging(o.nc),
	)

	err := handler.Handle(context.Background())
	if err != nil {
		if errors.Is(err, constant.ErrOutboxNotFound) {
			return poller.ErrNoData
		}

		return err
	}

	return nil
}
