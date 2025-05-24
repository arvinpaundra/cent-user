package model

import (
	"time"

	"github.com/guregu/null/v6"
)

type Session struct {
	ID           int64
	UserId       int64
	AccessToken  string
	RefreshToken null.String
	CreatedAt    time.Time
	UpdatedAt    time.Time
	DeletedAt    null.Time
}
