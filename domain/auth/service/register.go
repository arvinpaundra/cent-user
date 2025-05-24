package service

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/request"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type RegisterHandler struct {
	userReader repository.UserReader
	userWriter repository.UserWriter
}

func NewRegisterHandler(
	userReader repository.UserReader,
	userWriter repository.UserWriter,
) RegisterHandler {
	return RegisterHandler{
		userReader: userReader,
		userWriter: userWriter,
	}
}

func (s RegisterHandler) Handle(ctx context.Context, payload request.Register) error {
	isExist, err := s.userReader.IsEmailExist(ctx, payload.Email)
	if err != nil {
		return err
	}

	if isExist {
		return constant.ErrEmailAlreadyTaken
	}

	user := entity.User{
		Email:    payload.Email,
		Fullname: payload.Fullname,
	}

	err = user.GeneratePassword(payload.Password)
	if err != nil {
		return err
	}

	err = s.userWriter.Save(ctx, user)
	if err != nil {
		return err
	}

	return nil
}
