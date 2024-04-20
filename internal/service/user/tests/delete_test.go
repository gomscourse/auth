package tests

import (
	"context"
	"fmt"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/repository"
	repositoryMocks "github.com/gomscourse/auth/internal/repository/mocks"
	userService "github.com/gomscourse/auth/internal/service/user"
	"github.com/gomscourse/common/pkg/db"
	commonMocks "github.com/gomscourse/common/pkg/db/mocks"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestDelete(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		id  int64
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

		id = gofakeit.Int64()

		queryName = "delete query"
		queryRow  = "sql query"

		action    = "user.Delete"
		modelName = "user"

		repoErrorDelete = fmt.Errorf("repo error delete")
		repoErrorLog    = fmt.Errorf("repo error log")

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
				ctx: ctx,
				id:  id,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.DeleteMock.Expect(ctx, id).Return(q, nil)
				mock.CreateLogMock.Expect(ctx, action, modelName, id, q).Return(nil)
				return mock
			},
		},
		{
			name: "repo delete error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErrorDelete,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.DeleteMock.Expect(ctx, id).Return(nil, repoErrorDelete)
				return mock
			},
		},
		{
			name: "repo create log error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			err: repoErrorLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.DeleteMock.Expect(ctx, id).Return(q, nil)
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

				err := service.Delete(tt.args.ctx, tt.args.id)
				require.Equal(t, tt.err, err)
			},
		)
	}
}
