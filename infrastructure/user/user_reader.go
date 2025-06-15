package user

import (
	"context"
	"errors"

	"github.com/arvinpaundra/cent/user/domain/user/constant"
	"github.com/arvinpaundra/cent/user/domain/user/entity"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
	"github.com/arvinpaundra/cent/user/model"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var _ repository.UserReader = (*UserReaderRepository)(nil)

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

func (r UserReaderRepository) FindUserByIdForUpdate(ctx context.Context, id int64) (*entity.User, error) {
	var userModel model.User

	err := r.db.Clauses(clause.Locking{Strength: "UPDATE"}).
		WithContext(ctx).
		Model(&model.User{}).
		Where("id = ?", id).
		First(&userModel).
		Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, constant.ErrUserNotFound
		}

		return nil, err
	}

	user := entity.User{
		ID:       userModel.ID,
		Email:    userModel.Email,
		Fullname: userModel.Fullname,
		Balance:  userModel.Balance,
		Currency: userModel.Currency,
		Slug:     userModel.Slug.Ptr(),
	}

	return &user, nil
}
