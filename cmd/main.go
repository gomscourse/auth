package main

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/config/env"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type server struct {
	desc.UnimplementedUserV1Server
	pool *pgxpool.Pool
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info := req.GetInfo()
	password := info.GetPassword()
	if password == "" {
		return &desc.CreateResponse{}, status.Errorf(codes.InvalidArgument, "password can't be empty")
	}

	// TODO: generate password hash

	builderInsert := sq.Insert(`"user"`).
		PlaceholderFormat(sq.Dollar).
		Columns("name", "email", "role").
		Values(info.GetName(), info.GetEmail(), info.GetRole()).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return &desc.CreateResponse{}, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	var userID int64
	err = s.pool.QueryRow(ctx, query, args...).Scan(&userID)
	if err != nil {
		return &desc.CreateResponse{}, status.Errorf(codes.Internal, "failed to insert user: %v", err)
	}

	log.Printf("inserted user with id: %d", userID)
	return &desc.CreateResponse{Id: userID}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	fmt.Printf("get id: %d\n", req.GetId())
	return &desc.GetResponse{User: &desc.User{Role: desc.Role_ADMIN}}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	fmt.Printf("%+v\n", req.GetInfo())
	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Printf("get id: %d\n", req.GetId())
	return &emptypb.Empty{}, nil
}

func main() {
	ctx := context.Background()
	// Считываем переменные окружения
	err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	grpcConfig, err := env.NewGRPCConfig()
	if err != nil {
		log.Fatalf("failed to get grpc config: %v", err)
	}

	pgConfig, err := env.NewPGConfig()
	if err != nil {
		log.Fatalf("failed to get pg config: %v", err)
	}

	lis, err := net.Listen("tcp", grpcConfig.Address())
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Создаем пул соединений с базой данных
	pool, err := pgxpool.Connect(ctx, pgConfig.DSN())
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer pool.Close()

	s := grpc.NewServer()
	reflection.Register(s)
	desc.RegisterUserV1Server(s, &server{pool: pool})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
