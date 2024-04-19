package user

import (
	"context"
	"github.com/gomscourse/auth/internal/converter"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	info := req.GetInfo()
	userID := req.GetId()

	err := i.userService.Update(ctx, converter.ToUserUpdateInfoFromDesc(info, userID))
	if err != nil {
		return nil, err
	}

	return nil, nil
}
