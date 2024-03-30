package user

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, info *model.UserUpdateInfo) error {
	err := s.userRepository.Update(ctx, info)
	if err != nil {
		return err
	}

	return nil
}
