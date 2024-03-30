package main

import (
	"context"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/config/env"
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/repository/user"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
	"log"
	"net"
)

type server struct {
	desc.UnimplementedUserV1Server
	userRepository repository.UserRepository
}

func (s *server) Create(ctx context.Context, req *desc.CreateRequest) (*desc.CreateResponse, error) {
	info := req.GetInfo()
	userID, err := s.userRepository.Create(ctx, info)
	if err != nil {
		//return &desc.CreateResponse{}, status.Errorf(codes.InvalidArgument, err.Error())
		return &desc.CreateResponse{}, err
	}

	log.Printf("inserted user with id: %d", userID)
	return &desc.CreateResponse{Id: userID}, nil
}

func (s *server) Get(ctx context.Context, req *desc.GetRequest) (*desc.GetResponse, error) {
	userObj, err := s.userRepository.Get(ctx, req.GetId())
	if err != nil {
		return &desc.GetResponse{}, err
	}

	return &desc.GetResponse{User: userObj}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	info := req.GetInfo()
	userID := req.GetId()

	err := s.userRepository.Update(ctx, userID, info)
	if err != nil {
		return &emptypb.Empty{}, err
	}

	return &emptypb.Empty{}, nil
}

func (s *server) Delete(ctx context.Context, req *desc.DeleteRequest) (*emptypb.Empty, error) {
	err := s.userRepository.Delete(ctx, req.GetId())
	if err != nil {
		return &emptypb.Empty{}, err
	}

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
	repo := user.NewRepository(pool)
	desc.RegisterUserV1Server(s, &server{userRepository: repo})

	log.Printf("server listening at %v", lis.Addr())

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
