package app

import (
	"context"
	"github.com/gomscourse/auth/internal/config"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/gomscourse/common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"net/http"
	"sync"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
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

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := app.runGRPCServer()
		if err != nil {
			log.Fatalf("failed to start GRPC server: %v", err)
		}
	}()

	go func() {
		defer wg.Done()

		err := app.runHTTPServer()
		if err != nil {
			log.Fatalf("failed to start HTTP server: %v", err)
		}
	}()

	wg.Wait()

	return nil
}

func (app *App) initDeps(ctx context.Context) error {
	funcs := []func(ctx context.Context) error{
		app.initConfig,
		app.initServiceProvider,
		app.initGRPCServer,
		app.initHTTPServer,
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

func (app *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := desc.RegisterUserV1HandlerFromEndpoint(ctx, mux, app.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	app.httpServer = &http.Server{
		Addr:    app.serviceProvider.HTTPConfig().Address(),
		Handler: mux,
	}

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

func (app *App) runHTTPServer() error {
	log.Printf("HTTP server is running on %s", app.serviceProvider.HTTPConfig().Address())

	err := app.httpServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
