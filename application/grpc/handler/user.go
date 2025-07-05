package handler

import (
	"context"

	"github.com/arvinpaundra/cent/user/application/command/user"
	"github.com/arvinpaundra/cent/user/domain/user/constant"
	"github.com/arvinpaundra/cent/user/domain/user/service"
	userinfra "github.com/arvinpaundra/cent/user/infrastructure/user"
	commonpb "github.com/arvinpaundra/centpb/gen/go/common/v1"
	userpb "github.com/arvinpaundra/centpb/gen/go/user/v1"
	"google.golang.org/grpc/codes"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB

	userpb.UnimplementedUserServiceServer
}

func NewUserService(db *gorm.DB) UserService {
	return UserService{db: db}
}

func (u UserService) FindUserBySlug(ctx context.Context, req *userpb.FindUserBySlugRequest) (*userpb.FindUserBySlugResponse, error) {
	slug := req.GetSlug()

	svc := service.NewFindUserBySlug(
		userinfra.NewUserReaderRepository(u.db),
	)

	result, err := svc.Exec(ctx, slug)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			return &userpb.FindUserBySlugResponse{
				Meta: &commonpb.Meta{
					Code:    codes.NotFound.String(),
					Message: err.Error(),
				},
			}, err
		default:
			return &userpb.FindUserBySlugResponse{
				Meta: &commonpb.Meta{
					Code:    codes.Internal.String(),
					Message: "internal server error",
				},
			}, err
		}
	}

	return &userpb.FindUserBySlugResponse{
		Meta: &commonpb.Meta{
			Code:    codes.OK.String(),
			Message: "success",
		},
		User: &userpb.User{
			Id:       result.ID,
			Email:    result.Email,
			Fullname: result.Fullname,
			Key:      result.Key,
			Image:    result.Image,
		},
	}, nil
}

func (u UserService) FindUserDetail(ctx context.Context, req *userpb.FindUserDetailRequest) (*userpb.FindUserDetailResponse, error) {
	key := req.GetId()

	svc := service.NewFindUserDetail(
		userinfra.NewUserReaderRepository(u.db),
	)

	result, err := svc.Exec(ctx, key)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			return &userpb.FindUserDetailResponse{
				Meta: &commonpb.Meta{
					Code:    codes.NotFound.String(),
					Message: err.Error(),
				},
			}, err
		default:
			return &userpb.FindUserDetailResponse{
				Meta: &commonpb.Meta{
					Code:    codes.Internal.String(),
					Message: "internal server error",
				},
			}, err
		}
	}

	return &userpb.FindUserDetailResponse{
		Meta: &commonpb.Meta{
			Code:    codes.OK.String(),
			Message: "success find user by key",
		},
		User: &userpb.User{
			Id:       result.ID,
			Email:    result.Email,
			Fullname: result.Fullname,
			Key:      result.Key,
			Image:    result.Image,
		},
	}, nil
}

func (u UserService) UpdateUserBalance(ctx context.Context, req *userpb.UpdateUserBalanceRequest) (*userpb.UpdateUserBalanceResponse, error) {
	command := user.UpdateUserBalance{
		UserId: req.GetId(),
		Amount: req.GetAmount(),
	}

	svc := service.NewUpdateUserBalance(
		userinfra.NewUserReaderRepository(u.db),
		userinfra.NewUserWriterRepository(u.db),
		userinfra.NewUnitOfWork(u.db),
	)

	err := svc.Exec(ctx, command)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			return &userpb.UpdateUserBalanceResponse{
				Meta: &commonpb.Meta{
					Code:    codes.InvalidArgument.String(),
					Message: err.Error(),
				},
			}, err
		default:
			return &userpb.UpdateUserBalanceResponse{
				Meta: &commonpb.Meta{
					Code:    codes.Internal.String(),
					Message: "internal server error",
				},
			}, err
		}
	}

	return &userpb.UpdateUserBalanceResponse{
		Meta: &commonpb.Meta{
			Code:    codes.OK.String(),
			Message: "success",
		},
	}, nil
}
