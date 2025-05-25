package auth

import (
	"context"
	"errors"

	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/arvinpaundra/cent/user/model"
	"gorm.io/gorm"
)

var _ repository.SessionReader = (*SessionReaderRepository)(nil)

type SessionReaderRepository struct {
	db *gorm.DB
}

func NewSessionReaderRepository(db *gorm.DB) SessionReaderRepository {
	return SessionReaderRepository{db: db}
}

func (r SessionReaderRepository) FindByRefreshToken(ctx context.Context, userId int64, refreshToken string) (entity.Session, error) {
	var session entity.Session

	err := r.db.WithContext(ctx).
		Model(&model.Session{}).
		Select("id", "user_id", "refresh_token").
		Where("user_id = ? AND refresh_token = ? AND deleted_at IS NULL", userId, refreshToken).
		Take(&session).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.Session{}, constant.ErrSessionNotFound
		}

		return entity.Session{}, err
	}

	return session, nil
}
