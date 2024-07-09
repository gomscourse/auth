package env

import (
	"github.com/gomscourse/auth/internal/config"
	"net"
	"os"
	"time"

	"github.com/pkg/errors"
)

const (
	grpcHostEnvName = "GRPC_HOST"
	grpcPortEnvName = "GRPC_PORT"
)

type grpcConfig struct {
	host              string
	port              string
	requestLimitCount int
	requestLimitTime  time.Duration
}

func NewGRPCConfig() (config.GRPCConfig, error) {
	host := os.Getenv(grpcHostEnvName)
	if len(host) == 0 {
		return nil, errors.New("grpc host not found")
	}

	port := os.Getenv(grpcPortEnvName)
	if len(port) == 0 {
		return nil, errors.New("grpc port not found")
	}

	return &grpcConfig{
		host:              host,
		port:              port,
		requestLimitCount: 100,
		requestLimitTime:  time.Second,
	}, nil
}

func (cfg *grpcConfig) Address() string {
	return net.JoinHostPort(cfg.host, cfg.port)
}

func (cfg *grpcConfig) RateLimit() (int, time.Duration) {
	return cfg.requestLimitCount, cfg.requestLimitTime
}
