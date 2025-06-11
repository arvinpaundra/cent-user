package user

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/user/entity"
	"github.com/arvinpaundra/cent/user/model"
	"gorm.io/gorm"
)

type UserReaderRepository struct {
	db *gorm.DB
}

func NewUserReaderRepository(db *gorm.DB) UserReaderRepository {
	return UserReaderRepository{db: db}
}

func (r UserReaderRepository) FindBySlug(ctx context.Context, slug string) (*entity.User, error) {
	var userModel model.User

	err := r.db.WithContext(ctx).
		Model(&model.User{}).
		Where("slug = ?", slug).
		First(&userModel).
		Error

	if err != nil {
		return nil, err
	}

	user := entity.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Fullname: userModel.Fullname,
		Image:    userModel.Image.Ptr(),
	}

	return &user, nil
}
