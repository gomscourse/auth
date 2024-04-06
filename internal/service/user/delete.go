package user

import (
	"context"
)

func (s *serv) Delete(ctx context.Context, id int64) error {
	err := s.txManager.ReadCommitted(ctx, func(ctx context.Context) error {
		q, err := s.userRepository.Delete(ctx, id)
		if err != nil {
			return err
		}

		err = s.userRepository.CreateLog(ctx, "user.Delete", "user", id, q)
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
