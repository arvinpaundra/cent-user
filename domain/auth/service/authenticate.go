package service

import (
	"context"
	"strconv"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/domain/auth/constant"
	authres"github.com/arvinpaundra/cent/user/application/response/auth"
	"github.com/arvinpaundra/cent/user/domain/auth/repository"
)

type Authenticate struct {
	userReader repository.UserReader
	userCache  repository.UserCache
	tokenable  token.Tokenable
}

func NewAuthenticate(
	userReader repository.UserReader,
	userCache repository.UserCache,
	tokenable token.Tokenable,
) Authenticate {
	return Authenticate{
		userReader: userReader,
		userCache:  userCache,
		tokenable:  tokenable,
	}
}

func (s Authenticate) Exec(ctx context.Context, token string) (authres.UserAuthenticated, error) {
	claims, err := s.tokenable.Decode(token)
	if err != nil {
		return authres.UserAuthenticated{}, err
	}

	var res authres.UserAuthenticated

	identifierStr := strconv.Itoa(int(claims.Identifier))
	key := constant.UserCachedKey + identifierStr

	userCached, err := s.userCache.Get(ctx, key)
	if err != nil && err != constant.ErrUserNotFound {
		return authres.UserAuthenticated{}, nil
	}

	if !userCached.IsEmpty() {
		res = authres.UserAuthenticated{
			UserID: userCached.ID,
		}

		return res, nil
	}

	user, err := s.userReader.FindById(ctx, claims.Identifier)
	if err != nil {
		return authres.UserAuthenticated{}, err
	}

	_ = s.userCache.Set(ctx, key, user, constant.TTLFiveMinutes)

	res = authres.UserAuthenticated{
		UserID: user.ID,
	}

	return res, nil
}
