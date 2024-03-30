package user

import (
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/service"
)

type serv struct {
	userRepository repository.UserRepository
}

func NewService(userRepository repository.UserRepository) service.UserService {
	return &serv{userRepository: userRepository}
}
