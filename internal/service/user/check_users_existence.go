package user

import (
	"context"
)

func (s *serv) CheckUsersExistence(ctx context.Context, usernames []string) error {
	_, err := s.userRepository.CheckUsersExistence(ctx, usernames)
	return err
}
