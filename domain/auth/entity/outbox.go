package entity

import (
	"time"

	"github.com/arvinpaundra/cent/user/domain/auth/constant"
)

type Outbox struct {
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
