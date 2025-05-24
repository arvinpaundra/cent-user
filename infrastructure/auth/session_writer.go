package auth

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/arvinpaundra/cent/user/model"
	"gorm.io/gorm"
)

var _ repository.SessionWriter = (*SessionWriterRepository)(nil)

type SessionWriterRepository struct {
	db *gorm.DB
}

func NewSessionWriterRepository(db *gorm.DB) SessionWriterRepository {
	return SessionWriterRepository{db: db}
}

func (r SessionWriterRepository) Save(ctx context.Context, refreshToken entity.Session) error {
	refreshTokenModel := refreshToken.ToModel()

	err := r.db.WithContext(ctx).
		Model(&model.Session{}).
		Create(&refreshTokenModel).
		Error

	if err != nil {
		return err
	}

	return nil
}
