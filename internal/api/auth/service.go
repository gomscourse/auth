package user

import (
	"github.com/gomscourse/auth/internal/service"
	desc "github.com/gomscourse/auth/pkg/auth_v1"
)

type Implementation struct {
	desc.UnimplementedAuthV1Server
	authService service.AuthService
}

func NewImplementation(authService service.AuthService) *Implementation {
	return &Implementation{authService: authService}
}
