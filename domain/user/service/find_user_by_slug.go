package service

import (
	"context"

	userres "github.com/arvinpaundra/cent/user/application/response/user"
	"github.com/arvinpaundra/cent/user/domain/user/repository"
)

type FindUserBySlug struct {
	userReader repository.UserReader
}

func NewFindUserBySlug(userReader repository.UserReader) FindUserBySlug {
	return FindUserBySlug{
		userReader: userReader,
	}
}

func (s FindUserBySlug) Exec(ctx context.Context, slug string) (userres.FindUserBySlug, error) {
	user, err := s.userReader.FindBySlug(ctx, slug)
	if err != nil {
		return userres.FindUserBySlug{}, err
	}

	res := userres.FindUserBySlug{
		ID:       user.ID,
		Email:    user.Email,
		Fullname: user.Fullname,
		Key:      user.Key,
		Image:    user.Image,
	}

	return res, nil
}
