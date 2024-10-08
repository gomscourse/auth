package auth

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/utils"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
	"github.com/gomscourse/common/pkg/sys/messages"
	"github.com/pkg/errors"
)

func (s *serv) GetRefreshToken(ctx context.Context, refreshToken string) (string, error) {
	claims, err := utils.VerifyToken(refreshToken, []byte(s.jwtConfig.RefreshTokenSecret()))
	if err != nil {
		return "", sys.NewCommonError(messages.RefreshTokenInvalid, codes.PermissionDenied)
	}

	user, _, err := s.userRepo.GetByUsername(ctx, claims.Username)
	if err != nil {
		return "", err
	}

	newRefreshToken, err := utils.GenerateToken(
		model.User{
			Username: claims.Username,
			Role:     user.Role,
		},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenLifetime(),
	)

	if err != nil {
		return "", errors.New("failed to generate token")
	}

	return newRefreshToken, nil
}
