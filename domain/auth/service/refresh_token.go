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

type RefreshTokenHandler struct {
	userReader    repository.UserReader
	sessionReader repository.SessionReader
	sessionWriter repository.SessionWriter
	userCache     repository.UserCache
	tokenable     token.Tokenable
}

func NewRefreshTokenHandler(
	userReader repository.UserReader,
	sessionReader repository.SessionReader,
	sessionWriter repository.SessionWriter,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) RefreshTokenHandler {
	return RefreshTokenHandler{
		userReader:    userReader,
		sessionReader: sessionReader,
		sessionWriter: sessionWriter,
		userCache:     userCache,
		tokenable:     tokenable,
	}
}

func (s RefreshTokenHandler) Handle(ctx context.Context, payload request.RefreshToken) (response.RefreshToken, error) {
	claims, err := s.tokenable.Decode(payload.RefreshToken)
	if err != nil {
		return response.RefreshToken{}, constant.ErrTokenInvalid
	}

	user, err := s.userReader.FindById(ctx, claims.Identifier)
	if err != nil {
		return response.RefreshToken{}, err
	}

	session, err := s.sessionReader.FindByRefreshToken(ctx, user.ID, payload.RefreshToken)
	if err != nil {
		return response.RefreshToken{}, err
	}

	accessToken, err := s.tokenable.Encode(session.UserId, constant.TokenValidFifteenMinutes, constant.TokenValidImmediately)
	if err != nil {
		return response.RefreshToken{}, err
	}

	refreshToken, err := s.tokenable.Encode(session.UserId, constant.TokenValidSevenDays, constant.TokenValidAfterFifteenMinutes)
	if err != nil {
		return response.RefreshToken{}, err
	}

	newSession := entity.Session{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
	}

	err = s.sessionWriter.Save(ctx, newSession)
	if err != nil {
		return response.RefreshToken{}, err
	}

	session.SetDeletedAt()

	err = s.sessionWriter.Revoke(ctx, session)
	if err != nil {
		return response.RefreshToken{}, err
	}

	identifierStr := strconv.Itoa(int(claims.Identifier))
	key := constant.UserCachedKey + identifierStr

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	res := response.RefreshToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
