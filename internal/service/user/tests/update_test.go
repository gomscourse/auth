package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/repository"
	repositoryMocks "github.com/gomscourse/auth/internal/repository/mocks"
	userService "github.com/gomscourse/auth/internal/service/user"
	"github.com/gomscourse/common/pkg/db"
	commonMocks "github.com/gomscourse/common/pkg/db/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestUpdate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		info *model.UserUpdateInfo
	}

	txManagerMock := commonMocks.NewTxManagerMock(t)
	txManagerMock.ReadCommittedMock.Set(
		func(ctx context.Context, handler db.Handler) (err error) {
			return handler(ctx)
		},
	)

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name  = gofakeit.Name()
		email = gofakeit.Email()
		id    = gofakeit.Int64()

		queryName = "update query"
		queryRow  = "sql query"

		action    = "user.Update"
		modelName = "user"

		repoErrorUpdate = fmt.Errorf("repo error update")
		repoErrorLog    = fmt.Errorf("repo error log")

		info = &model.UserUpdateInfo{
			ID:    id,
			Name:  name,
			Email: email,
		}

		q = &db.Query{
			Name:     queryName,
			QueryRow: queryRow,
		}
	)

	tests := []struct {
		name               string
		args               args
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.UpdateMock.Expect(ctx, info).Return(q, nil)
				mock.CreateLogMock.Expect(ctx, action, modelName, id, q).Return(nil)
				return mock
			},
		},
		{
			name: "repo update error case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			err: repoErrorUpdate,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.UpdateMock.Expect(ctx, info).Return(nil, repoErrorUpdate)
				return mock
			},
		},
		{
			name: "repo create log error case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			err: repoErrorLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.UpdateMock.Expect(ctx, info).Return(q, nil)
				mock.CreateLogMock.Expect(ctx, action, modelName, id, q).Return(repoErrorLog)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				chatRepoMock := tt.userRepositoryMock(mc)
				service := userService.NewTestService(chatRepoMock, txManagerMock)

				err := service.Update(tt.args.ctx, tt.args.info)
				require.Equal(t, tt.err, err)
			},
		)
	}
}
