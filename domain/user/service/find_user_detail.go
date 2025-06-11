package service

import (
	"context"

	userres "github.com/arvinpaundra/cent/user/application/response/user"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
)

type FindUserDetailHandler struct {
	userReader repository.UserReader
}

func NewFindUserDetailHandler(userReader repository.UserReader) FindUserDetailHandler {
	return FindUserDetailHandler{
		userReader: userReader,
	}
}

func (s FindUserDetailHandler) Handle(ctx context.Context, slug string) (userres.FindUserDetail, error) {
	user, err := s.userReader.FindBySlug(ctx, slug)
	if err != nil {
		return userres.FindUserDetail{}, err
	}

	res := userres.FindUserDetail{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
		Image:    user.Image,
	}

	return res, nil
}
