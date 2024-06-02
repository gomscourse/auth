package repository

import (
	"context"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/common/pkg/db"
)

type UserRepository interface {
	Create(ctx context.Context, info *model.UserCreateInfo) (int64, *db.Query, error)
	Get(ctx context.Context, id int64) (*model.User, *db.Query, error)
	GetByUsername(ctx context.Context, username string) (*model.User, *db.Query, error)
	Update(ctx context.Context, info *model.UserUpdateInfo) (*db.Query, error)
	Delete(ctx context.Context, id int64) (*db.Query, error)

	CreateLog(ctx context.Context, action, model string, modelId int64, loggedQ *db.Query) error
}

type AccessRepository interface {
	GetRuleByEndpoint(ctx context.Context, endpoint string) ([]*model.AccessRule, *db.Query, error)
}
