package app

import (
	"context"
	userApi "github.com/gomscourse/auth/internal/api/user"
	"github.com/gomscourse/auth/internal/client/db"
	"github.com/gomscourse/auth/internal/client/db/pg"
	"github.com/gomscourse/auth/internal/closer"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/config/env"
	"github.com/gomscourse/auth/internal/repository"
	userRepo "github.com/gomscourse/auth/internal/repository/user"
	"github.com/gomscourse/auth/internal/service"
	userService "github.com/gomscourse/auth/internal/service/user"
	"log"
)

type serviceProvider struct {
	pgConfig   config.PGConfig
	grpcConfig config.GRPCConfig

	dbClient           db.Client
	userRepository     repository.UserRepository
	userService        service.UserService
	userImplementation *userApi.Implementation
}

func newServiceProvider() *serviceProvider {
	return &serviceProvider{}
}

func (sp *serviceProvider) PgConfig() config.PGConfig {
	if sp.pgConfig == nil {
		cfg, err := env.NewPGConfig()
		if err != nil {
			log.Fatalf("failed to load PG config: %s", err.Error())
		}
		sp.pgConfig = cfg
	}

	return sp.pgConfig

}

func (sp *serviceProvider) GRPCConfig() config.GRPCConfig {
	if sp.grpcConfig == nil {
		cfg, err := env.NewGRPCConfig()
		if err != nil {
			log.Fatalf("failed to load GRPC config: %s", err.Error())
		}

		sp.grpcConfig = cfg
	}

	return sp.grpcConfig
}

func (sp *serviceProvider) DbClient(ctx context.Context) db.Client {
	if sp.dbClient == nil {
		client, err := pg.New(ctx, sp.PgConfig().DSN())
		if err != nil {
			log.Fatalf("failed to load DB client: %s", err.Error())
		}

		err = client.DB().Ping(ctx)
		if err != nil {
			log.Fatalf("failed to ping DB: %s", err.Error())
		}

		closer.Add(client.Close)

		sp.dbClient = client
	}

	return sp.dbClient
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepository == nil {
		sp.userRepository = userRepo.NewRepository(sp.DbClient(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewService(sp.UserRepository(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) UserImplementation(ctx context.Context) *userApi.Implementation {
	if sp.userImplementation == nil {
		sp.userImplementation = userApi.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImplementation
}
