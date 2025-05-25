package service

import (
	"context"
	"strconv"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/request"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/response"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type LoginHandler struct {
	userReader    repository.UserReader
	sessionWriter repository.SessionWriter
	userCache     repository.UserCache
	tokenable     token.Tokenable
}

func NewLoginHandler(
	userReader repository.UserReader,
	sessionWriter repository.SessionWriter,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) LoginHandler {
	return LoginHandler{
		userReader:    userReader,
		sessionWriter: sessionWriter,
		userCache:     userCache,
		tokenable:     tokenable,
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

	err = s.sessionWriter.Save(ctx, session)
	if err != nil {
		return response.Login{}, err
	}

	identifierStr := strconv.Itoa(int(user.ID))
	key := constant.UserCachedKey + identifierStr

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	res := response.Login{
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
