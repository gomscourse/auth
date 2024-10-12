package tests

import (
	"context"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	"github.com/gomscourse/auth/internal/api/user"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/service"
	"github.com/gomscourse/auth/internal/service/mocks"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.CreateRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		name            = gofakeit.Name()
		email           = gofakeit.Email()
		password        = gofakeit.Email()
		passwordConfirm = password
		role            = gofakeit.Int32()

		id           = gofakeit.Int64()
		serviceError = errors.New("service error")

		req = &desc.CreateRequest{
			Info: &desc.UserCreateInfo{
				Username:        name,
				Email:           email,
				Password:        password,
				PasswordConfirm: passwordConfirm,
				Role:            desc.Role(role),
			},
		}

		res = &desc.CreateResponse{
			Id: id,
		}

		info = &model.UserCreateInfo{
			Name:            name,
			Email:           email,
			Password:        password,
			PasswordConfirm: passwordConfirm,
			Role:            role,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.CreateResponse
		err             error
		userServiceMock userServiceMockFunc
	}{
		{
			name: "success case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: res,
			err:  nil,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(id, nil)
				return mock
			},
		},
		{
			name: "error case",
			args: args{
				ctx: ctx,
				req: req,
			},
			want: nil,
			err:  serviceError,
			userServiceMock: func(mc *minimock.Controller) service.UserService {
				mock := mocks.NewUserServiceMock(mc)
				mock.CreateMock.Expect(ctx, info).Return(0, serviceError)
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

				result, err := api.Create(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}

}
