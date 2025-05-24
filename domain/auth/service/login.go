package service

import (
	"context"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/request"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/response"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type LoginHandler struct {
	userReader         repository.UserReader
	refreshTokenWriter repository.SessionWriter
	tokenable          token.Tokenable
}

func NewLoginHandler(
	userReader repository.UserReader,
	refreshTokenWriter repository.SessionWriter,
	tokenable token.Tokenable,
) LoginHandler {
	return LoginHandler{
		userReader:         userReader,
		refreshTokenWriter: refreshTokenWriter,
		tokenable:          tokenable,
	}
}

func (s LoginHandler) Handle(ctx context.Context, payload request.Login) (response.Login, error) {
	user, err := s.userReader.FindByEmail(ctx, payload.Email)
	if err != nil {
		return response.Login{}, err
	}

	if !user.ComparePassword(payload.Password) {
		return response.Login{}, constant.ErrWrongEmailOrPassword
	}

	accessToken, err := s.tokenable.Encode(user.ID, constant.TokenValidFifteenMinutes, constant.TokenValidImmediately)
	if err != nil {
		return response.Login{}, err
	}

	refreshToken, err := s.tokenable.Encode(user.ID, constant.TokenValidSevenDays, constant.TokenValidAfterFifteenMinutes)
	if err != nil {
		return response.Login{}, err
	}

	session := entity.Session{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
	}

	err = s.refreshTokenWriter.Save(ctx, session)
	if err != nil {
		return response.Login{}, err
	}

	res := response.Login{
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
