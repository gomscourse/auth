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
	CheckUsersExistence(ctx context.Context, usernames []string) error
}

type AuthService interface {
	Login(ctx context.Context, username, password string) (string, error)
	GetRefreshToken(ctx context.Context, refreshToken string) (string, error)
	GetAccessToken(ctx context.Context, refreshToken string) (string, error)
}

type AccessService interface {
	Check(ctx context.Context, endpointAddress string) error
}
