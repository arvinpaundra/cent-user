package user

import (
	"context"

	"github.com/arvinpaundra/cent/user/application/command/user"
	"github.com/arvinpaundra/cent/user/domain/user/constant"
	"github.com/arvinpaundra/cent/user/domain/user/service"
	userinfra "github.com/arvinpaundra/cent/user/infrastructure/user"
	"github.com/arvinpaundra/centpb/gen/go/common/v1"
	userpb "github.com/arvinpaundra/centpb/gen/go/user/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"gorm.io/gorm"
)

type UserService struct {
	db *gorm.DB

	userpb.UnimplementedUserServiceServer
}

func NewUserService(db *gorm.DB) UserService {
	return UserService{db: db}
}

func (u UserService) FindUserDetail(ctx context.Context, req *userpb.FindUserDetailRequest) (*userpb.FindUserDetailResponse, error) {
	slug := req.GetSlug()

	handler := service.NewFindUserDetailHandler(
		userinfra.NewUserReaderRepository(u.db),
	)

	result, err := handler.Handle(ctx, slug)
	if err != nil {
		return nil, err
	}

	var image *wrapperspb.StringValue

	if result.Image != nil {
		image = wrapperspb.String(*result.Image)
	}

	res := userpb.FindUserDetailResponse{
		Id:       result.ID,
		Email:    result.Email,
		Fullname: result.Fullname,
		Image:    image,
	}

	return &res, nil
}

func (u UserService) UpdateUserBalance(ctx context.Context, req *userpb.UpdateUserBalanceRequest) (*userpb.UpdateUserBalanceResponse, error) {
	command := user.UpdateUserBalance{
		UserId: req.GetId(),
		Amount: req.GetAmount(),
	}

	handler := service.NewUpdateUserBalanceHandler(
		userinfra.NewUserReaderRepository(u.db),
		userinfra.NewUserWriterRepository(u.db),
		userinfra.NewUnitOfWork(u.db),
	)

	err := handler.Handle(ctx, command)
	if err != nil {
		switch err {
		case constant.ErrUserNotFound:
			return &userpb.UpdateUserBalanceResponse{
				Meta: &common.Meta{
					Code:    codes.InvalidArgument.String(),
					Message: err.Error(),
				},
			}, err
		default:
			return &userpb.UpdateUserBalanceResponse{
				Meta: &common.Meta{
					Code:    codes.Internal.String(),
					Message: "internal server error",
				},
			}, err
		}
	}

	return &userpb.UpdateUserBalanceResponse{
		Meta: &common.Meta{
			Code:    codes.OK.String(),
			Message: "success",
		},
	}, nil
}
