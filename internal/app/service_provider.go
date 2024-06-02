package app

import (
	"context"
	accessApi "github.com/gomscourse/auth/internal/api/access"
	authApi "github.com/gomscourse/auth/internal/api/auth"
	userApi "github.com/gomscourse/auth/internal/api/user"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/config/env"
	"github.com/gomscourse/auth/internal/repository"
	accessRepo "github.com/gomscourse/auth/internal/repository/access"
	userRepo "github.com/gomscourse/auth/internal/repository/user"
	"github.com/gomscourse/auth/internal/service"
	accessService "github.com/gomscourse/auth/internal/service/access"
	authService "github.com/gomscourse/auth/internal/service/auth"
	userService "github.com/gomscourse/auth/internal/service/user"
	"github.com/gomscourse/common/pkg/closer"
	"github.com/gomscourse/common/pkg/db"
	"github.com/gomscourse/common/pkg/db/pg"
	"github.com/gomscourse/common/pkg/db/transaction"
	"log"
)

type serviceProvider struct {
	pgConfig      config.PGConfig
	grpcConfig    config.GRPCConfig
	httpConfig    config.HTTPConfig
	swaggerConfig config.HTTPConfig
	jwtConfig     config.JWTConfig

	dbClient             db.Client
	txManager            db.TxManager
	userRepository       repository.UserRepository
	accessRepository     repository.AccessRepository
	userService          service.UserService
	authService          service.AuthService
	accessService        service.AccessService
	userImplementation   *userApi.Implementation
	authImplementation   *authApi.Implementation
	accessImplementation *accessApi.Implementation
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

func (sp *serviceProvider) HTTPConfig() config.HTTPConfig {
	if sp.httpConfig == nil {
		cfg, err := env.NewHTTPConfig()
		if err != nil {
			log.Fatalf("failed to load HTTP config: %s", err.Error())
		}

		sp.httpConfig = cfg
	}

	return sp.httpConfig
}

func (sp *serviceProvider) SwaggerConfig() config.SwaggerConfig {
	if sp.swaggerConfig == nil {
		cfg, err := env.NewSwaggerConfig()
		if err != nil {
			log.Fatalf("failed to load Swagger config: %s", err.Error())
		}

		sp.swaggerConfig = cfg
	}

	return sp.swaggerConfig
}

func (sp *serviceProvider) JWTConfig() config.JWTConfig {
	if sp.jwtConfig == nil {
		cfg, err := env.NewJWTConfig()
		if err != nil {
			log.Fatalf("failed to load JWT config: %s", err.Error())
		}

		sp.jwtConfig = cfg
	}

	return sp.jwtConfig
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

func (sp *serviceProvider) TxManager(ctx context.Context) db.TxManager {
	if sp.txManager == nil {
		sp.txManager = transaction.NewTransactionManager(sp.DbClient(ctx).DB())
	}

	return sp.txManager
}

func (sp *serviceProvider) UserRepository(ctx context.Context) repository.UserRepository {
	if sp.userRepository == nil {
		sp.userRepository = userRepo.NewRepository(sp.DbClient(ctx))
	}

	return sp.userRepository
}

func (sp *serviceProvider) AccessRepository(ctx context.Context) repository.AccessRepository {
	if sp.accessRepository == nil {
		sp.accessRepository = accessRepo.NewRepository(sp.DbClient(ctx))
	}

	return sp.accessRepository
}

func (sp *serviceProvider) UserService(ctx context.Context) service.UserService {
	if sp.userService == nil {
		sp.userService = userService.NewService(sp.UserRepository(ctx), sp.TxManager(ctx))
	}

	return sp.userService
}

func (sp *serviceProvider) AuthService(ctx context.Context) service.AuthService {
	if sp.authService == nil {
		sp.authService = authService.NewService(sp.UserRepository(ctx), sp.JWTConfig())
	}

	return sp.authService
}

func (sp *serviceProvider) AccessService(ctx context.Context) service.AccessService {
	if sp.accessService == nil {
		sp.accessService = accessService.NewService(sp.UserRepository(ctx), sp.AccessRepository(ctx), sp.JWTConfig())
	}

	return sp.accessService
}

func (sp *serviceProvider) UserImplementation(ctx context.Context) *userApi.Implementation {
	if sp.userImplementation == nil {
		sp.userImplementation = userApi.NewImplementation(sp.UserService(ctx))
	}

	return sp.userImplementation
}

func (sp *serviceProvider) AuthImplementation(ctx context.Context) *authApi.Implementation {
	if sp.authImplementation == nil {
		sp.authImplementation = authApi.NewImplementation(sp.AuthService(ctx))
	}

	return sp.authImplementation
}

func (sp *serviceProvider) AccessImplementation(ctx context.Context) *accessApi.Implementation {
	if sp.accessImplementation == nil {
		sp.accessImplementation = accessApi.NewImplementation(sp.AccessService(ctx))
	}

	return sp.accessImplementation
}
