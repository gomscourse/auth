package user

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	var userObj *model.User
	err := s.txManager.ReadCommitted(
		ctx, func(ctx context.Context) error {
			user, q, err := s.userRepository.Get(ctx, id)
			if err != nil {
				return err
			}

			err = s.userRepository.CreateLog(ctx, "user.Get", "user", user.ID, q)
			if err != nil {
				return err
			}

			userObj = user

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return userObj, nil
}
