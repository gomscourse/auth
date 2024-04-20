package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/repository"
)

type userRepositoryMockFunc func(mc *minimock.Controller) repository.UserRepository
