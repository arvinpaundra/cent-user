package service

import (
	"context"
	"strconv"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	authcmd"github.com/arvinpaundra/cent/user/application/command/auth"
	authres"github.com/arvinpaundra/cent/user/application/response/auth"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type Login struct {
	userReader    repository.UserReader
	sessionWriter repository.SessionWriter
	userCache     repository.UserCache
	tokenable     token.Tokenable
}

func NewLogin(
	userReader repository.UserReader,
	sessionWriter repository.SessionWriter,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) Login {
	return Login{
		userReader:    userReader,
		sessionWriter: sessionWriter,
		userCache:     userCache,
		tokenable:     tokenable,
	}
}

func (s Login) Exec(ctx context.Context, payload authcmd.Login) (authres.Login, error) {
	user, err := s.userReader.FindByEmail(ctx, payload.Email)
	if err != nil {
		return authres.Login{}, err
	}

	if !user.ComparePassword(payload.Password) {
		return authres.Login{}, constant.ErrWrongEmailOrPassword
	}

	accessToken, err := s.tokenable.Encode(user.ID, constant.TokenValidFifteenMinutes, constant.TokenValidImmediately)
	if err != nil {
		return authres.Login{}, err
	}

	refreshToken, err := s.tokenable.Encode(user.ID, constant.TokenValidSevenDays, constant.TokenValidAfterFifteenMinutes)
	if err != nil {
		return authres.Login{}, err
	}

	session := entity.Session{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
	}

	err = s.sessionWriter.Save(ctx, &session)
	if err != nil {
		return authres.Login{}, err
	}

	identifierStr := strconv.Itoa(int(user.ID))
	key := constant.UserCachedKey + identifierStr

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	res := authres.Login{
		Email:        user.Email,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
