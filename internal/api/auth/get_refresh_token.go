package user

import (
	"context"
	desc "github.com/gomscourse/auth/pkg/auth_v1"
)

func (i *Implementation) GetRefreshToken(
	ctx context.Context,
	req *desc.GetRefreshTokenRequest,
) (*desc.GetRefreshTokenResponse, error) {
	refreshToken, err := i.authService.GetRefreshToken(ctx, req.GetRefreshToken())
	if err != nil {
		return nil, err
	}

	return &desc.GetRefreshTokenResponse{RefreshToken: refreshToken}, nil
}
