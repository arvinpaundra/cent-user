package model

import (
	"time"

	"github.com/guregu/null/v6"
)

type User struct {
	ID        int64
	Email     string
	Fullname  string
	Balance   float64
	Currency  string
	Key       string
	Password  null.String
	Image     null.String
	Slug      null.String
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt null.Time
}
