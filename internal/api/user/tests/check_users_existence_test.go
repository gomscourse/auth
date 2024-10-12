package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/api/user"
	"github.com/gomscourse/auth/internal/service"
	"github.com/gomscourse/auth/internal/service/mocks"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCheckUsersExistence(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CheckUsersExistenceRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		usernames    = []string{gofakeit.Username(), gofakeit.Username()}
		serviceError = errors.New("service error")

		req = &desc.CheckUsersExistenceRequest{
			Usernames: usernames,
		}
	)

	tests := []struct {
		name            string
		args            args
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			err: serviceError,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CheckUsersExistenceMock.Expect(ctx, usernames).Return(serviceError)
				return mock
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

				userServiceMock := tt.userServiceMock(mc)
				api := user.NewImplementation(userServiceMock)

				_, err := api.CheckUsersExistence(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
			},
		)
	}

}
