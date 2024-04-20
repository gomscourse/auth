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

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx  context.Context
		info *model.UserCreateInfo
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

		name               = gofakeit.Name()
		email              = gofakeit.Email()
		password           = gofakeit.Password(true, true, true, true, true, 8)
		passwordConfirm    = password
		passwordConfirmErr = gofakeit.Password(true, true, true, true, true, 9)
		role               = gofakeit.Int32()
		id                 = gofakeit.Int64()

		queryName = "create query"
		queryRow  = "sql query"

		action    = "user.Create"
		modelName = "user"

		passwordsErr    = fmt.Errorf("passwords are not equal")
		repoErrorCreate = fmt.Errorf("repo error create")
		repoErrorLog    = fmt.Errorf("repo error log")

		info = &model.UserCreateInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}

		infoErr = &model.UserCreateInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirmErr,
			Role:            role,
		}

		q = &db.Query{
			Name:     queryName,
			QueryRow: queryRow,
		}
	)

	tests := []struct {
		name               string
		args               args
		want               int64
		err                error
		userRepositoryMock userRepositoryMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: id,
			err:  nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.CreateMock.Expect(ctx, info).Return(id, q, nil)
				mock.CreateLogMock.Expect(ctx, action, modelName, id, q).Return(nil)
				return mock
			},
		},
		{
			name: "passwords not equal case",
			args: args{
				ctx:  ctx,
				info: infoErr,
			},
			want: 0,
			err:  passwordsErr,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				return repositoryMocks.NewUserRepositoryMock(t)
			},
		},
		{
			name: "repo create error case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: 0,
			err:  repoErrorCreate,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.CreateMock.Expect(ctx, info).Return(0, nil, repoErrorCreate)
				return mock
			},
		},
		{
			name: "repo create log error case",
			args: args{
				ctx:  ctx,
				info: info,
			},
			want: 0,
			err:  repoErrorLog,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.CreateMock.Expect(ctx, info).Return(id, q, nil)
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

				result, err := service.Create(tt.args.ctx, tt.args.info)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}
}
