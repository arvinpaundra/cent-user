package model

import (
	"time"

	"github.com/guregu/null/v6"
	"gorm.io/datatypes"
)

type Outbox struct {
	ID          int64
	Event       string
	Status      string
	Payload     datatypes.JSON
	PublishedAt null.Time
	Error       null.String
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func (Outbox) TableName() string {
	return "outbox"
}
