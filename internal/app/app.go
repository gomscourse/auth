package app

import (
	"context"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/interceptor"
	descAccess "github.com/gomscourse/auth/pkg/access_v1"
	descAuth "github.com/gomscourse/auth/pkg/auth_v1"
	descUser "github.com/gomscourse/auth/pkg/user_v1"
	_ "github.com/gomscourse/auth/statik"
	"github.com/gomscourse/common/pkg/closer"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rakyll/statik/fs"
	"github.com/rs/cors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
	httpServer      *http.Server
	swaggerServer   *http.Server
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
	wg.Add(3)

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

	go func() {
		defer wg.Done()

		err := app.runSwaggerServer()
		if err != nil {
			log.Fatalf("failed to start Swagger server: %v", err)
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
		app.initSwaggerServer,
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
	app.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.UnaryInterceptor(interceptor.ValidateInterceptor),
	)

	reflection.Register(app.grpcServer)

	descUser.RegisterUserV1Server(app.grpcServer, app.serviceProvider.UserImplementation(ctx))
	descAuth.RegisterAuthV1Server(app.grpcServer, app.serviceProvider.AuthImplementation(ctx))
	descAccess.RegisterAccessV1Server(app.grpcServer, app.serviceProvider.AccessImplementation(ctx))
	return nil
}

func (app *App) initHTTPServer(ctx context.Context) error {
	mux := runtime.NewServeMux()

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}

	err := descUser.RegisterUserV1HandlerFromEndpoint(ctx, mux, app.serviceProvider.GRPCConfig().Address(), opts)
	if err != nil {
		return err
	}

	corsMiddleware := cors.New(
		cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "Content-Length", "Authorization"},
			AllowCredentials: true,
		},
	)

	app.httpServer = &http.Server{
		Addr:    app.serviceProvider.HTTPConfig().Address(),
		Handler: corsMiddleware.Handler(mux),
	}

	return nil
}

func (app *App) initSwaggerServer(_ context.Context) error {
	statikFs, err := fs.New()
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/", http.StripPrefix("/", http.FileServer(statikFs)))
	mux.HandleFunc("/api.swagger.json", serveSwaggerFile("/api.swagger.json"))

	app.swaggerServer = &http.Server{
		Addr:    app.serviceProvider.SwaggerConfig().Address(),
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

func (app *App) runSwaggerServer() error {
	log.Printf("Swagger server is running on %s", app.serviceProvider.SwaggerConfig().Address())

	err := app.swaggerServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}

func serveSwaggerFile(path string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Serving swagger file: %s", path)

		statikFs, err := fs.New()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Open swagger file: %s", path)

		file, err := statikFs.Open(path)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer file.Close()

		log.Printf("Read swagger file: %s", path)

		content, err := io.ReadAll(file)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Write swagger file: %s", path)

		w.Header().Set("Content-Type", "application/json")
		_, err = w.Write(content)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		log.Printf("Served swagger file: %s", path)
	}
}
