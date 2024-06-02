package model

import (
	"database/sql"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type User struct {
	ID                            int64
	Username, Email, PasswordHash string
	Role                          int32
	CreatedAt                     time.Time
	UpdatedAt                     sql.NullTime
}

type UserClaims struct {
	jwt.StandardClaims
	Username string `json:"username"`
	Role     int32  `json:"role"`
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
