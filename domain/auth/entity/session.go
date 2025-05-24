package entity

import (
	"time"

	"github.com/arvinpaundra/cent/user/model"
	"github.com/guregu/null/v6"
)

type Session struct {
	ID           int64
	UserId       int64
	AccessToken  string
	RefreshToken *string
	DeletedAt    *time.Time
}

func (e *Session) ToModel() model.Session {
	return model.Session{
		ID:           e.ID,
		UserId:       e.UserId,
		AccessToken:  e.AccessToken,
		RefreshToken: null.StringFromPtr(e.RefreshToken),
	}
}
