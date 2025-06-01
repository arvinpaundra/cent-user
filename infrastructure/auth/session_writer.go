package auth

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/arvinpaundra/cent/user/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.SessionWriter = (*SessionWriterRepository)(nil)

type SessionWriterRepository struct {
	db *gorm.DB
}

func NewSessionWriterRepository(db *gorm.DB) SessionWriterRepository {
	return SessionWriterRepository{db: db}
}

func (r SessionWriterRepository) Save(ctx context.Context, session entity.Session) error {
	if session.IsNew() {
		return r.insert(ctx, session)
	}

	return nil
}

func (r SessionWriterRepository) insert(ctx context.Context, session entity.Session) error {
	sessionModel := model.Session{
		UserId:       session.UserId,
		AccessToken:  session.AccessToken,
		RefreshToken: null.StringFromPtr(session.RefreshToken),
	}

	err := r.db.WithContext(ctx).
		Model(&model.Session{}).
		Create(&sessionModel).
		Error

	if err != nil {
		return err
	}

	return nil
}

func (r SessionWriterRepository) Revoke(ctx context.Context, session entity.Session) error {
	sessionModel := model.Session{
		UserId:       session.UserId,
		AccessToken:  session.AccessToken,
		RefreshToken: null.StringFromPtr(session.RefreshToken),
		DeletedAt:    null.TimeFromPtr(session.DeletedAt),
	}

	err := r.db.WithContext(ctx).
		Model(&model.Session{}).
		Where("id = ?", sessionModel.ID).
		Updates(&sessionModel).
		Error

	if err != nil {
		return err
	}

	return nil
}
