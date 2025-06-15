package user

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/user/entity"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
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
	if user.IsMarkedToBeUpdated() {
		return r.update(ctx, user)
	}

	return nil
}

func (r UserWriterRepository) update(ctx context.Context, user *entity.User) error {
	userModel := model.User{
		Fullname: user.Fullname,
		Balance:  user.Balance,
		Image:    null.StringFromPtr(user.Image),
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
