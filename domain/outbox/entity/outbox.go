package entity

import (
	"time"

	"github.com/arvinpaundra/cent/user/core/trait"
	"github.com/arvinpaundra/cent/user/domain/outbox/constant"
)

type Outbox struct {
	trait.Updateable

	ID          int64
	Event       constant.OutboxEvent
	Status      constant.OutboxStatus
	Payload     []byte
	PublishedAt *time.Time
	Error       *string
}

func (e *Outbox) IsNew() bool {
	return e.ID <= 0
}

func (e *Outbox) SetStatus(status constant.OutboxStatus) {
	e.Status = status
}

func (e *Outbox) SetPublishedAt() {
	now := time.Now().UTC()
	e.PublishedAt = &now
}

func (e *Outbox) SetError(err string) {
	e.Error = &err
}
