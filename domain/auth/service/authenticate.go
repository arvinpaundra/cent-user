package service

import (
	"context"
	"strconv"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	"github.com/arvinpaundra/cent/user/domain/auth/dto/response"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type AuthenticateHandler struct {
	userReader repository.UserReader
	userCache  repository.UserCache
	tokenable  token.Tokenable
}

func NewAuthenticateHandler(
	userReader repository.UserReader,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) AuthenticateHandler {
	return AuthenticateHandler{
		userReader: userReader,
		userCache:  userCache,
		tokenable:  tokenable,
	}
}

func (s AuthenticateHandler) Handle(ctx context.Context, token string) (response.UserAuthenticated, error) {
	claims, err := s.tokenable.Decode(token)
	if err != nil {
		return response.UserAuthenticated{}, err
	}

	var res response.UserAuthenticated

	identifierStr := strconv.Itoa(int(claims.Identifier))
	key := constant.UserCachedKey + identifierStr

	userCached, err := s.userCache.Get(ctx, key)
	if err != nil && err != constant.ErrUserNotFound {
		return response.UserAuthenticated{}, nil
	}

	if !userCached.IsEmpty() {
		res = response.UserAuthenticated{
			UserID: userCached.ID,
		}

		return res, nil
	}

	user, err := s.userReader.FindById(ctx, claims.Identifier)
	if err != nil {
		return response.UserAuthenticated{}, err
	}

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	res = response.UserAuthenticated{
		UserID: user.ID,
	}

	return res, nil
}
