package auth

import (
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/service"
)

type serv struct {
	userRepo  repository.UserRepository
	jwtConfig config.JWTConfig
}

func NewService(userRepository repository.UserRepository, jwtConfig config.JWTConfig) service.AuthService {
	return &serv{userRepo: userRepository, jwtConfig: jwtConfig}
}
