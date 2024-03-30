package user

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

func (s *serv) Create(ctx context.Context, info *model.UserCreateInfo) (int64, error) {
	userID, err := s.userRepository.Create(ctx, info)
	if err != nil {
		//return &desc.CreateResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
		return 0, err
	}

	return userID, nil
}
