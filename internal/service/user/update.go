package user

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Update(ctx context.Context, info *model.UserUpdateInfo) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		q, err := s.userRepository.Update(ctx, info)
		if err != nil {
			return err
		}

		err = s.userRepository.CreateLog(ctx, "user.Update", "user", info.ID, q)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}
