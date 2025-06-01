package entity

import "time"

type Session struct {
	ID           int64
	UserId       int64
	AccessToken  string
	RefreshToken *string
	DeletedAt    *time.Time
}

func (e *Session) IsNew() bool {
	return e.ID <= 0
}

func (e *Session) SetDeletedAt() {
	now := time.Now().UTC()
	e.DeletedAt = &now
}
