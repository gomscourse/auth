package config

import (
	"github.com/joho/godotenv"
	"time"
)

func Load() error {
	err := godotenv.Load(".env")
	if err != nil {
		return err
	}

	return nil
}

type GRPCConfig interface {
	Address() string
}

type JWTConfig interface {
	RefreshTokenSecret() string
	AccessTokenSecret() string
	RefreshTokenLifetime() time.Duration
	AccessTokenLifetime() time.Duration
}

type PGConfig interface {
	DSN() string
}

type HTTPConfig interface {
	Address() string
}

type SwaggerConfig interface {
	Address() string
}
