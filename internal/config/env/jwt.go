package env

import (
	"github.com/gomscourse/auth/internal/config"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	refreshTokenSecretEnvName = "REFRESH_TOKEN_SECRET"
	accessTokenSecretEnvName  = "ACCESS_TOKEN_SECRET"
)

type jwtConfig struct {
	refreshTokenSecret string
	accessTokenSecret  string
	refreshLifetime    time.Duration
	accessLifetime     time.Duration
}

func (j jwtConfig) RefreshTokenSecret() string {
	return j.refreshTokenSecret
}

func (j jwtConfig) AccessTokenSecret() string {
	return j.accessTokenSecret
}

func (j jwtConfig) RefreshTokenLifetime() time.Duration {
	return j.refreshLifetime
}

func (j jwtConfig) AccessTokenLifetime() time.Duration {
	return j.accessLifetime
}

func NewJWTConfig() (config.JWTConfig, error) {
	refreshSecret := os.Getenv(refreshTokenSecretEnvName)
	if len(refreshSecret) == 0 {
		return nil, errors.New("refresh token secret not found")
	}

	accessSecret := os.Getenv(accessTokenSecretEnvName)
	if len(accessSecret) == 0 {
		return nil, errors.New("access secret token not found")
	}

	return &jwtConfig{
		refreshTokenSecret: refreshSecret,
		accessTokenSecret:  accessSecret,
		refreshLifetime:    24 * 60 * time.Minute,
		accessLifetime:     10 * time.Minute,
	}, nil
}
