package tests

import (
	"context"
	"database/sql"
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

func TestGet(t *testing.T) {
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

		id      = gofakeit.Int64()
		name    = gofakeit.Name()
		email   = gofakeit.Email()
		role    = gofakeit.Int32()
		created = gofakeit.Date()
		updated = sql.NullTime{
			Time:  created,
			Valid: true,
		}

		queryName = "get query"
		queryRow  = "sql query"

		action    = "user.Get"
		modelName = "user"

		repoErrorGet = fmt.Errorf("repo error get")
		repoErrorLog = fmt.Errorf("repo error log")

		user = &model.User{
			ID:        id,
			Username:  name,
			Email:     email,
			Role:      role,
			CreatedAt: created,
			UpdatedAt: updated,
		}

		q = &db.Query{
			Name:     queryName,
			QueryRow: queryRow,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               *model.User
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: user,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.GetMock.Expect(ctx, id).Return(user, q, nil)
				mock.CreateLogMock.Expect(ctx, action, modelName, id, q).Return(nil)
				return mock
			},
		},
		{
			name: "repo get error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErrorGet,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.GetMock.Expect(ctx, id).Return(nil, nil, repoErrorGet)
				return mock
			},
		},
		{
			name: "repo create log error case",
			args: args{
				ctx: ctx,
				id:  id,
			},
			want: nil,
			err:  repoErrorLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.GetMock.Expect(ctx, id).Return(user, q, nil)
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

				result, err := service.Get(tt.args.ctx, tt.args.id)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
