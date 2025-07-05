package auth

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
	"github.com/arvinpaundra/cent/user/model"
	"github.com/guregu/null/v6"
	"gorm.io/gorm"
)

var _ repository.UserWriter = (*UserWriterRepository)(nil)

type UserWriterRepository struct {
	db *gorm.DB
}

func NewUserWriterRepository(db *gorm.DB) UserWriterRepository {
	return UserWriterRepository{db: db}
}

func (r UserWriterRepository) Save(ctx context.Context, user *entity.User) error {
	if user.IsNew() {
		return r.insert(ctx, user)
	} else if user.IsMarkedToBeUpdated() {
		return r.update(ctx, user)
	}

	return nil
}

func (r UserWriterRepository) insert(ctx context.Context, user *entity.User) error {
	userModel := model.User{
		Email:    user.Email,
		Password: null.StringFromPtr(user.Password),
		Fullname: user.Fullname,
		Key:      user.Key,
	}

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Create(&userModel).
		Error

	if err != nil {
		return err
	}

	user.ID = userModel.ID

	return nil
}

func (r UserWriterRepository) update(ctx context.Context, user *entity.User) error {
	userModel := model.User{
		Fullname: user.Fullname,
		Image:    null.StringFromPtr(user.Image),
		Slug:     null.StringFromPtr(user.Slug),
	}

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", user.ID).
		Updates(&userModel).
		Error

	if err != nil {
		return err
	}

	return nil
}
