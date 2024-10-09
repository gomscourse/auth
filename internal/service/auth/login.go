package auth

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/utils"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
)

func (s *serv) Login(ctx context.Context, username, password string) (string, error) {
	// Лезем в базу или кэш за данными пользователя
	user, _, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return "", err
	}
	// Сверяем хэши пароля
	if !utils.VerifyPassword(user.PasswordHash, password) {
		return "", sys.NewCommonError("Invalid login or password", codes.Unauthenticated)
	}

	refreshToken, err := utils.GenerateToken(
		model.User{
			Username: username,
			Role:     user.Role,
		},
		[]byte(s.jwtConfig.RefreshTokenSecret()),
		s.jwtConfig.RefreshTokenLifetime(),
	)

	if err != nil {
		return "", sys.NewCommonError("failed to generate token", codes.Unauthenticated)
	}

	return refreshToken, nil
}
