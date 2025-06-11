package user

import (
	"context"

	"github.com/arvinpaundra/cent/user/domain/user/service"
	userinfra "github.com/arvinpaundra/cent/user/infrastructure/user"
	userpb "github.com/arvinpaundra/centpb/gen/go/user/v1"
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
	handler := service.NewFindUserDetailHandler(
		userinfra.NewUserReaderRepository(u.db),
	)

	result, err := handler.Handle(ctx, req.Slug)
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
