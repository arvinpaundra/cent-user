package service

import (
	"context"
	"strconv"

	authcmd "github.com/arvinpaundra/cent/user/application/command/auth"
	authres "github.com/arvinpaundra/cent/user/application/response/auth"
	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/entity"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type RefreshToken struct {
	userReader    repository.UserReader
	sessionReader repository.SessionReader
	sessionWriter repository.SessionWriter
	unitOfWork    repository.UnitOfWork
	userCache     repository.UserCache
	tokenable     token.Tokenable
}

func NewRefreshToken(
	userReader repository.UserReader,
	sessionReader repository.SessionReader,
	sessionWriter repository.SessionWriter,
	unitOfWork repository.UnitOfWork,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) RefreshToken {
	return RefreshToken{
		userReader:    userReader,
		sessionReader: sessionReader,
		sessionWriter: sessionWriter,
		unitOfWork:    unitOfWork,
		userCache:     userCache,
		tokenable:     tokenable,
	}
}

func (s RefreshToken) Exec(ctx context.Context, payload authcmd.RefreshToken) (authres.RefreshToken, error) {
	claims, err := s.tokenable.Decode(payload.RefreshToken)
	if err != nil {
		return authres.RefreshToken{}, constant.ErrTokenInvalid
	}

	user, err := s.userReader.FindById(ctx, claims.Identifier)
	if err != nil {
		return authres.RefreshToken{}, err
	}

	session, err := s.sessionReader.FindByRefreshToken(ctx, user.ID, payload.RefreshToken)
	if err != nil {
		return authres.RefreshToken{}, err
	}

	accessToken, err := s.tokenable.Encode(session.UserId, constant.TokenValidFifteenMinutes, constant.TokenValidImmediately)
	if err != nil {
		return authres.RefreshToken{}, err
	}

	refreshToken, err := s.tokenable.Encode(session.UserId, constant.TokenValidSevenDays, constant.TokenValidAfterFifteenMinutes)
	if err != nil {
		return authres.RefreshToken{}, err
	}

	tx, err := s.unitOfWork.Begin()
	if err != nil {
		return authres.RefreshToken{}, nil
	}

	newSession := entity.Session{
		UserId:       user.ID,
		AccessToken:  accessToken,
		RefreshToken: &refreshToken,
	}

	err = tx.SessionWriter().Save(ctx, &newSession)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return authres.RefreshToken{}, uowErr
		}

		return authres.RefreshToken{}, err
	}

	session.SetDeletedAt()

	err = s.sessionWriter.Revoke(ctx, session)
	if err != nil {
		if uowErr := tx.Rollback(); uowErr != nil {
			return authres.RefreshToken{}, uowErr
		}

		return authres.RefreshToken{}, err
	}

	identifierStr := strconv.Itoa(int(claims.Identifier))
	key := constant.UserCachedKey + identifierStr

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	if uowErr := tx.Commit(); uowErr != nil {
		return authres.RefreshToken{}, uowErr
	}

	res := authres.RefreshToken{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	return res, nil
}
