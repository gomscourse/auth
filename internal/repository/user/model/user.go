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
