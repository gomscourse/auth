package tests

import (
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/service"
)

type userServiceMockFunc func(mc *minimock.Controller) service.UserService
