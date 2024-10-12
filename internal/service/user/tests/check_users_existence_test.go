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

func TestCheckUsersExistence(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx       context.Context
		usernames []string
	}

	txManagerMock := commonMocks.NewTxManagerMock(t)

	var (
		ctx       = context.Background()
		mc        = minimock.NewController(t)
		usernames = []string{gofakeit.Username(), gofakeit.Username()}
		queryName = "create query"
		queryRow  = "sql query"

		q = &db.Query{
			Name:     queryName,
			QueryRow: queryRow,
		}

		repoError = fmt.Errorf("repo error")
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
				ctx:       ctx,
				usernames: usernames,
			},
			err: nil,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(q, nil)
				return mock
			},
		},
		{
			name: "repo error case",
			args: args{
				ctx:       ctx,
				usernames: usernames,
			},
			err: repoError,
			userRepositoryMock: func(mc *minimock.Controller) repository.UserRepository {
				mock := repositoryMocks.NewUserRepositoryMock(t)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(nil, repoError)
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

				err := service.CheckUsersExistence(tt.args.ctx, tt.args.usernames)
				require.Equal(t, tt.err, err)
			},
		)
	}
}
