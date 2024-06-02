package user

import (
	"github.com/gomscourse/auth/internal/service"
	desc "github.com/gomscourse/auth/pkg/access_v1"
)

type Implementation struct {
	desc.UnimplementedAccessV1Server
	accessService service.AccessService
}

func NewImplementation(accessService service.AccessService) *Implementation {
	return &Implementation{accessService: accessService}
}
