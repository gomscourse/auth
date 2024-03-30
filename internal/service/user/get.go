package user

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Get(ctx context.Context, id int64) (*model.User, error) {
	userObj, err := s.userRepository.Get(ctx, id)
	if err != nil {
		return &model.User{}, err
	}

	return userObj, nil
}
