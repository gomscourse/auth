package tests

import (
	"context"
	"database/sql"
	"errors"
	"github.com/brianvoe/gofakeit"
	"github.com/gojuno/minimock/v3"
	userApi "github.com/gomscourse/auth/internal/api/user"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/service"
	"github.com/gomscourse/auth/internal/service/mocks"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"
	"testing"
	"time"
)

func TestGet(t *testing.T) {
	t.Parallel()

	type args struct {
		ctx context.Context
		req *desc.GetRequest
	}

	var (
		ctx = context.Background()
		mc  = minimock.NewController(t)

		id      = gofakeit.Int64()
		name    = gofakeit.Name()
		email   = gofakeit.Email()
		role    = gofakeit.Int32()
		created = time.Now()
		updated = sql.NullTime{
			Time:  created,
			Valid: true,
		}

		serviceError = errors.New("service error")

		req = &desc.GetRequest{
			Id: id,
		}

		res = &desc.GetResponse{
			User: &desc.User{
				Id:        id,
				Username:  name,
				Email:     email,
				Role:      desc.Role(role),
				CreatedAt: timestamppb.New(created),
				UpdatedAt: timestamppb.New(updated.Time),
			},
		}

		user = &model.User{
			ID:        id,
			Username:  name,
			Email:     email,
			Role:      role,
			CreatedAt: created,
			UpdatedAt: updated,
		}
	)

	tests := []struct {
		name            string
		args            args
		want            *desc.GetResponse
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
				mock.GetMock.Expect(ctx, id).Return(user, nil)
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
				mock.GetMock.Expect(ctx, id).Return(nil, serviceError)
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
				api := userApi.NewImplementation(userServiceMock)

				result, err := api.Get(tt.args.ctx, tt.args.req)
				require.Equal(t, tt.err, err)
				require.Equal(t, tt.want, result)
			},
		)
	}

}
