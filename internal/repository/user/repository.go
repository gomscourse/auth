package user

import (
	"context"
	"errors"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/auth/internal/client/db"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/repository/user/converter"
	repoModel "github.com/gomscourse/auth/internal/repository/user/model"
	"github.com/jackc/pgx/v4"
	"time"
)

const (
	tableName = "auth_user"

	idColumn        = "id"
	nameColumn      = "name"
	emailColumn     = "email"
	roleColumn      = "role"
	createdAtColumn = "created_at"
	updatedAtColumn = "updated_at"
)

const (
	logTableName      = "log"
	logActionColumn   = "action"
	logModelColumn    = "model"
	logModelIdColumn  = "model_id"
	logQueryColumn    = "query"
	logQueryRowColumn = "query_row"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.UserRepository {
	return &repo{db: db}
}

func (r *repo) Create(ctx context.Context, info *model.UserCreateInfo) (int64, *db.Query, error) {
	password := info.Password
	if password == "" {
		return 0, nil, errors.New("password can't be empty")
	}

	// TODO: generate password hash

	builderInsert := sq.Insert(tableName).
		PlaceholderFormat(sq.Dollar).
		Columns(nameColumn, emailColumn, roleColumn).
		Values(info.Name, info.Email, info.Role).
		Suffix("RETURNING id")

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return 0, nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "create_user_query",
		QueryRow: query,
	}

	var userID int64
	err = r.db.DB().QueryRowContext(ctx, q, args...).Scan(&userID)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, &q, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, *db.Query, error) {
	builderSelect := sq.Select(idColumn, nameColumn, emailColumn, roleColumn, createdAtColumn, updatedAtColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{idColumn: id}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return &model.User{}, nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "get_user_query",
		QueryRow: query,
	}

	var user repoModel.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return &model.User{}, nil, fmt.Errorf("user with id %d not found", id)
	}

	if err != nil {
		return &model.User{}, nil, fmt.Errorf("failed to get user: %w", err)
	}

	return converter.ToUserFromRepo(&user), &q, nil
}

func (r *repo) Update(ctx context.Context, info *model.UserUpdateInfo) (*db.Query, error) {

	buildUpdate := sq.Update(tableName).
		PlaceholderFormat(sq.Dollar).
		Set(nameColumn, info.Name).
		Set(emailColumn, info.Email).
		Set(updatedAtColumn, time.Now()).
		Where(sq.Eq{idColumn: info.ID})

	query, args, err := buildUpdate.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "update_user_query",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return &q, nil
}

func (r *repo) Delete(ctx context.Context, id int64) (*db.Query, error) {
	deleteBuilder := sq.Delete(tableName).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: id})
	query, args, err := deleteBuilder.ToSql()

	q := db.Query{
		Name:     "delete_user_query",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to delete user: %w", err)
	}
	return &q, nil
}

func (r *repo) CreateLog(ctx context.Context, action, model string, modelId int64, loggedQ *db.Query) error {
	builderInsert := sq.Insert(logTableName).
		PlaceholderFormat(sq.Dollar).
		Columns(logActionColumn, logModelColumn, logModelIdColumn, logQueryColumn, logQueryRowColumn).
		Values(action, model, modelId, loggedQ.Name, loggedQ.QueryRow)

	query, args, err := builderInsert.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "create_log_query",
		QueryRow: query,
	}

	_, err = r.db.DB().ExecContext(ctx, q, args...)
	if err != nil {
		return fmt.Errorf("failed to insert log: %w", err)
	}

	return nil
}
