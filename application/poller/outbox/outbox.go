package outbox

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/outbox/service"
	"github.com/arvinpaundra/cent/user/infrastructure/outbox"
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
		outbox.NewOutboxReaderRepository(o.db),
		outbox.NewOutboxWriterRepository(o.db),
		outbox.NewUnitOfWork(o.db),
		outbox.NewMessaging(o.nc),
	)

	err := handler.Handle(context.Background())
	if err != nil {
		return err
	}

	return nil
}
