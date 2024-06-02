package model

type AccessRule struct {
	Endpoint string `db:"endpoint"`
	Role     int32  `db:"role"`
}
