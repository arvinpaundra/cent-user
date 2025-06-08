package poller

import (
	"github.com/arvinpaundra/cent/user/application/poller/outbox"
	"github.com/arvinpaundra/cent/user/core/poller"
	"github.com/nats-io/nats.go"
	"gorm.io/gorm"
)

func StartWorker(p *poller.Poller, db *gorm.DB, nc *nats.Conn) {
	outbox := outbox.NewOutbox(db, nc)

	p.Spawn(outbox.OutboxProcessor)
}
