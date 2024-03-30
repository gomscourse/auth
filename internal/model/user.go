package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID          int64
	Name, Email string
	Role        int32
	CreatedAt   time.Time
	UpdatedAt   sql.NullTime
}

type UserCreateInfo struct {
	Name            string
	Email           string
	Password        string
	PasswordConfirm string
	Role            int32
}

type UserUpdateInfo struct {
	ID    int64
	Name  string
	Email string
}
