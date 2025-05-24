package model

import (
	"time"

	"github.com/guregu/null/v6"
)

type User struct {
	ID        int64
	Email     string
	Password  null.String
	Fullname  string
	Image     null.String
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}
