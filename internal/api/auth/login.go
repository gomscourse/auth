package user

import (
	"context"
	desc "github.com/gomscourse/auth/pkg/auth_v1"
	"github.com/pkg/errors"
)

func (i *Implementation) Login(ctx context.Context, req *desc.LoginRequest) (*desc.LoginResponse, error) {
	refreshToken, err := i.authService.Login(ctx, req.GetUsername(), req.GetPassword())
	if err != nil {
		return nil, errors.New("failed to log in")
	}

	return &desc.LoginResponse{
		RefreshToken: refreshToken,
	}, nil
}
