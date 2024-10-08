package user

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/repository/user/converter"
	repoModel "github.com/gomscourse/auth/internal/repository/user/model"
	"github.com/gomscourse/common/pkg/db"
	"github.com/gomscourse/common/pkg/sys"
	"github.com/gomscourse/common/pkg/sys/codes"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"slices"
	"strings"
	"time"
)

const (
	usersTable = "auth_user"

	idColumn           = "id"
	usernameColumn     = "username"
	passwordHashColumn = "password_hash"
	emailColumn        = "email"
	roleColumn         = "role"
	createdAtColumn    = "created_at"
	updatedAtColumn    = "updated_at"
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
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to generate password hash: %w", err)
	}

	builderInsert := sq.Insert(usersTable).
		PlaceholderFormat(sq.Dollar).
		Columns(usernameColumn, passwordHashColumn, emailColumn, roleColumn).
		Values(info.Name, string(passwordHash), info.Email, info.Role).
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
	err = r.db.DB().QueryRowContextScan(ctx, &userID, q, args...)
	if err != nil {
		return 0, nil, fmt.Errorf("failed to insert user: %w", err)
	}

	return userID, &q, nil
}

func (r *repo) Get(ctx context.Context, id int64) (*model.User, *db.Query, error) {
	return r.GetOneByColumn(ctx, idColumn, id)
}

func (r *repo) GetByUsername(ctx context.Context, username string) (*model.User, *db.Query, error) {
	return r.GetOneByColumn(ctx, usernameColumn, username)
}

func (r *repo) GetOneByColumn(ctx context.Context, column string, value any) (*model.User, *db.Query, error) {
	builderSelect := sq.Select(
		idColumn,
		usernameColumn,
		passwordHashColumn,
		emailColumn,
		roleColumn,
		createdAtColumn,
		updatedAtColumn,
	).
		From(usersTable).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{column: value}).
		Limit(1)

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return &model.User{}, nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     fmt.Sprintf("get_user_by_%s_query", column),
		QueryRow: query,
	}

	var user repoModel.User
	err = r.db.DB().ScanOneContext(ctx, &user, q, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return &model.User{}, nil, fmt.Errorf("user with %s %v not found", column, value)
	}

	if err != nil {
		return &model.User{}, nil, fmt.Errorf("failed to get user: %w", err)
	}

	return converter.ToUserFromRepo(&user), &q, nil
}

func (r *repo) Update(ctx context.Context, info *model.UserUpdateInfo) (*db.Query, error) {

	buildUpdate := sq.Update(usersTable).
		PlaceholderFormat(sq.Dollar).
		Set(usernameColumn, info.Name).
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
	deleteBuilder := sq.Delete(usersTable).PlaceholderFormat(sq.Dollar).Where(sq.Eq{idColumn: id})
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

func (r *repo) CheckUsersExistence(ctx context.Context, usernames []string) (*db.Query, error) {
	var users []string

	querySelect := sq.Select(usernameColumn).
		PlaceholderFormat(sq.Dollar).
		From(usersTable).
		Where(sq.Eq{usernameColumn: usernames})

	query, args, err := querySelect.ToSql()
	if err != nil {
		return nil, errors.Errorf("failed to build query: %s", err)
	}

	q := db.Query{
		Name:     "select users by names",
		QueryRow: query,
	}

	err = r.db.DB().ScanAllContext(ctx, &users, q, args...)
	if err != nil {
		return &q, errors.Errorf("failed to get users: %s", err)
	}

	requestedCount := len(usernames)
	actualCount := len(users)
	if actualCount < requestedCount {
		absent := make([]string, 0, requestedCount-actualCount)

		for _, u := range usernames {
			if !slices.Contains(users, u) {
				absent = append(absent, u)
			}
		}

		return &q, sys.NewCommonError(
			fmt.Sprintf("The following users do not exist: %s", strings.Join(absent, ",")),
			codes.InvalidArgument,
		)
	}

	return &q, nil
}
