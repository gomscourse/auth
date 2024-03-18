package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/auth/internal/config"
	"github.com/gomscourse/auth/internal/config/env"
	desc "github.com/gomscourse/auth/pkg/user_v1"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"log"
	"net"
	"time"
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

	builderInsert := sq.Insert("auth_user").
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
	builderSelect := sq.Select("id", "name", "email", "role", "created_at", "updated_at").
		From("auth_user").
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{"id": req.GetId()}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return &desc.GetResponse{}, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	row := s.pool.QueryRow(ctx, query, args...)
	if err != nil {
		return &desc.GetResponse{}, status.Errorf(codes.Internal, "failed to select user: %v", err)
	}

	var id int
	var name, email string
	var role int32
	var createdAt time.Time
	var updatedAt sql.NullTime

	err = row.Scan(&id, &name, &email, &role, &createdAt, &updatedAt)
	if errors.Is(err, pgx.ErrNoRows) {
		return &desc.GetResponse{}, status.Errorf(codes.NotFound, "user with id %d not found", req.GetId())
	}

	if err != nil {
		return &desc.GetResponse{}, status.Errorf(codes.Internal, "failed to get user: %v", err)
	}

	user := &desc.User{
		Id:        int64(id),
		Name:      name,
		Email:     email,
		Role:      desc.Role(role),
		CreatedAt: timestamppb.New(createdAt),
	}

	if updatedAt.Valid {
		user.UpdatedAt = timestamppb.New(updatedAt.Time)
	}

	return &desc.GetResponse{User: user}, nil
}

func (s *server) Update(ctx context.Context, req *desc.UpdateRequest) (*emptypb.Empty, error) {
	info := req.GetInfo()
	userID := req.GetId()

	buildUpdate := sq.Update("auth_user").
		PlaceholderFormat(sq.Dollar).
		Set("name", info.GetName().GetValue()).
		Set("email", info.GetEmail().GetValue()).
		Set("updated_at", time.Now()).
		Where(sq.Eq{"id": userID})

	query, args, err := buildUpdate.ToSql()
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to build query: %v", err)
	}

	_, err = s.pool.Exec(ctx, query, args...)
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "failed to update user: %v", err)
	}

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
