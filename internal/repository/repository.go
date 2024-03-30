package repository

import (
	"context"
	desc "github.com/gomscourse/auth/pkg/user_v1"
)

type UserRepository interface {
	Create(ctx context.Context, info *desc.UserCreateInfo) (int64, error)
	Get(ctx context.Context, id int64) (*desc.User, error)
	Update(ctx context.Context, id int64, info *desc.UpdateUserInfo) error
	Delete(ctx context.Context, id int64) error
}
