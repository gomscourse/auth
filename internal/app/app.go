package app

import (
	"context"
	"github.com/gomscourse/auth/internal/closer"
	"github.com/gomscourse/auth/internal/config"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	app := &App{}

	err := app.initDeps(ctx)
	if err != nil {
		return nil, err
	}
	return app, nil
}

func (app *App) Run() error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	return app.runGRPCServer()
}

func (app *App) initDeps(ctx context.Context) error {
	funcs := []func(ctx context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
	}

	for _, f := range funcs {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (app *App) initConfig(_ context.Context) error {
	err := config.Load()

	if err != nil {
		return err
	}

	return nil
}

func (app *App) initServiceProvider(_ context.Context) error {
	app.serviceProvider = newServiceProvider()
	return nil
}

func (app *App) initGRPCServer(ctx context.Context) error {
	app.grpcServer = grpc.NewServer()

	reflection.Register(app.grpcServer)

	desc.RegisterUserV1Server(app.grpcServer, app.serviceProvider.UserImplementation(ctx))
	return nil
}

func (app *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", app.serviceProvider.GRPCConfig().Address())
	if err != nil {
		return err
	}

	log.Printf("server listening at %v", lis.Addr())

	if err = app.grpcServer.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	return nil
}
