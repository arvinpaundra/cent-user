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

var _ repository.UserReader = (*UserReaderRepository)(nil)

type UserReaderRepository struct {
	db *gorm.DB
}

func NewUserReaderRepository(db *gorm.DB) UserReaderRepository {
	return UserReaderRepository{db: db}
}

func (r UserReaderRepository) FindByEmail(ctx context.Context, email string) (entity.User, error) {
	var userModel model.User

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("email = ?", email).
		Take(&userModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, constant.ErrUserNotFound
		}

		return entity.User{}, err
	}

	user := entity.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Password: userModel.Password.Ptr(),
		Fullname: userModel.Fullname,
		Image:    userModel.Image.Ptr(),
	}

	return user, nil
}

func (r UserReaderRepository) FindById(ctx context.Context, id int64) (entity.User, error) {
	var userModel model.User

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		First(&userModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity.User{}, constant.ErrUserNotFound
		}

		return entity.User{}, err
	}

	user := entity.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Password: userModel.Password.Ptr(),
		Fullname: userModel.Fullname,
		Image:    userModel.Image.Ptr(),
	}

	return user, nil
}

func (r UserReaderRepository) IsEmailExist(ctx context.Context, email string) (bool, error) {
	var total int64

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Select("id").
		Where("email = ?", email).
		Count(&total).
		Error

	if err != nil {
		return false, err
	}

	return total > 0, nil
}
