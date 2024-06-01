package auth

import "github.com/gomscourse/auth/internal/repository"

type serv struct {
	userRepo repository.UserRepository
}
