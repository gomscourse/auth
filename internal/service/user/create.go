package user

import (
	"context"
	"errors"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserCreateInfo) (int64, error) {
	if info.Password != info.PasswordConfirm {
		return 0, errors.New("passwords are not equal")
	}

	var userID int64
	err := s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			id, q, err := s.userRepository.Create(ctx, info)
			if err != nil {
				return err
			}

			err = s.userRepository.CreateLog(ctx, "user.Create", "user", id, q)
			if err != nil {
				return err
			}

			userID = id

			return nil
		},
	)

	if err != nil {
		return 0, err
	}

	return userID, nil
}
