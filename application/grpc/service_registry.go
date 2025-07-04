package grpc

import (
	"github.com/arvinpaundra/cent/user/application/grpc/handler"
	"github.com/arvinpaundra/cent/user/core/validator"
	"github.com/arvinpaundra/centpb/gen/go/auth/v1"
	"github.com/arvinpaundra/centpb/gen/go/user/v1"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
	"gorm.io/gorm"
)

func Register(srv *grpc.Server, db *gorm.DB, rdb *redis.Client, vld *validator.Validator) {
	auth.RegisterAuthenticateServiceServer(srv, handler.NewAuthService(db, rdb, vld))
	user.RegisterUserServiceServer(srv, handler.NewUserService(db))
}
