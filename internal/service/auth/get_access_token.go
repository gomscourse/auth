package auth

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/utils"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
	"github.com/gomscourse/common/pkg/sys/messages"
)

func (s *serv) GetAccessToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return "", sys.NewCommonError(messages.RefreshTokenInvalid, codes.PermissionDenied)
	}

	user, _, err := s.userRepo.GetByUsername(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	accessToken, err := utils.GenerateToken(
		model.User{
			Username: claims.Username,
			Role:     user.Role,
		},
		[]byte(s.jwtConfig.AccessTokenSecret()),
		s.jwtConfig.AccessTokenLifetime(),
	)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}
