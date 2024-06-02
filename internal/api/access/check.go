package user

import (
	"context"
	desc "github.com/gomscourse/auth/pkg/access_v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (i *Implementation) Check(context.Context, *desc.CheckRequest) (*emptypb.Empty, error) {

	return nil, nil
}
