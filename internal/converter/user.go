package converter

import (
	"github.com/gomscourse/auth/internal/model"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func ToUserFromService(user *model.User) *desc.User {
	var updatedAt *timestamppb.Timestamp
	if user.UpdatedAt.Valid {
		updatedAt = timestamppb.New(user.UpdatedAt.Time)
	}

	return &desc.User{
		Id:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		Role:      desc.Role(user.Role),
		CreatedAt: timestamppb.New(user.CreatedAt),
		UpdatedAt: updatedAt,
	}
}

func ToUserCreateInfoFromDesc(info *desc.UserCreateInfo) *model.UserCreateInfo {
	return &model.UserCreateInfo{
		Name:            info.GetUsername(),
		Email:           info.GetEmail(),
		Password:        info.GetPassword(),
		PasswordConfirm: info.GetPasswordConfirm(),
		Role:            int32(info.GetRole()),
	}
}

func ToUserUpdateInfoFromDesc(info *desc.UpdateUserInfo, id int64) *model.UserUpdateInfo {
	return &model.UserUpdateInfo{
		ID:    id,
		Name:  info.GetUsername().GetValue(),
		Email: info.GetEmail().GetValue(),
	}
}
