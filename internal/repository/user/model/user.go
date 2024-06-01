package model

import (
	"database/sql"
	"time"
)

type User struct {
	ID           int64        `db:"id"`
	Username     string       `db:"username"`
	PasswordHash string       `db:"password_hash"`
	Email        string       `db:"email"`
	Role         int32        `db:"role"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
}
