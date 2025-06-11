package auth

import (
	"context"

	"github.com/arvinpaundra/cent/user/core/token"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/arvinpaundra/cent/user/domain/auth/service"
	authinfra "github.com/arvinpaundra/cent/user/infrastructure/auth"
	authpb "github.com/arvinpaundra/centpb/gen/go/auth/v1"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
)

type AuthService struct {
	db  *gorm.DB
	rdb *redis.Client
	vld *validator.Validator

	authpb.UnimplementedAuthenticateServiceServer
}

func NewAuthService(db *gorm.DB, rdb *redis.Client, vld *validator.Validator) AuthService {
	return AuthService{
		db:  db,
		rdb: rdb,
		vld: vld,
	}
}

func (a AuthService) CheckSession(ctx context.Context, req *authpb.CheckSessionRequest) (*authpb.CheckSessionResponse, error) {
	tokenStr := req.GetToken()

	handler := service.NewAuthenticateHandler(
		authinfra.NewUserReaderRepository(a.db),
		authinfra.NewUserCacheRepository(a.rdb),
		token.NewJWT(viper.GetString("JWT_SECRET")),
	)

	result, err := handler.Handle(ctx, tokenStr)
	if err != nil {
		return nil, err
	}

	resp := authpb.CheckSessionResponse{
		UserId: result.UserID,
	}

	return &resp, nil
}
