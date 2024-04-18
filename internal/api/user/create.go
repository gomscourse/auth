package user

import (
	"context"
	"github.com/gomscourse/auth/internal/converter"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"log"
)

func (i *Implementation) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info := req.GetInfo()
	userID, err := i.userService.Create(ctx, converter.ToUserCreateInfoFromDesc(info))
	if err != nil {
		//return &desc.CreateResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
		return nil, err
	}

	log.Printf("inserted user with id: %d", userID)
	return &desc.CreateResponse{Id: userID}, nil
}
