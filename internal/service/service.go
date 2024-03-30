package service

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
)

type UserService interface {
	Create(ctx context.Context, info *model.UserCreateInfo) (int64, error)
	Get(ctx context.Context, id int64) (*model.User, error)
	Update(ctx context.Context, info *model.UserUpdateInfo) error
	Delete(ctx context.Context, id int64) error
}
