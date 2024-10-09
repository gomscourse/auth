package user

import (
	"context"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) CheckUsersExistence(ctx context.Context, req *desc.CheckUsersExistenceRequest) (
	*emptypb.Empty,
	error,
) {
	users := req.GetUsernames()
	err := i.userService.CheckUsersExistence(ctx, users)
	if err != nil {
		//return &desc.CreateResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
