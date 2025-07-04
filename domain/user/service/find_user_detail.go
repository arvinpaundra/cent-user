package service

import (
	"context"

	userres "github.com/arvinpaundra/cent/user/application/response/user"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
)

type FindUserDetail struct {
	userReader repository.UserReader
}

func NewFindUserDetail(userReader repository.UserReader) FindUserDetail {
	return FindUserDetail{
		userReader: userReader,
	}
}

func (s FindUserDetail) Exec(ctx context.Context, id int64) (userres.FindUserDetail, error) {
	user, err := s.userReader.FindById(ctx, id)
	if err != nil {
		return userres.FindUserDetail{}, nil
	}

	res := userres.FindUserDetail{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
		Key:      user.Key,
		Image:    user.Image,
	}

	return res, nil
}
