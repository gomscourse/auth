package access

import (
	"context"
	"fmt"
	sq "github.com/Masterminds/squirrel"
	"github.com/gomscourse/auth/internal/model"
	"github.com/gomscourse/auth/internal/repository"
	"github.com/gomscourse/auth/internal/repository/access/converter"
	repoModel "github.com/gomscourse/auth/internal/repository/access/model"
	"github.com/gomscourse/common/pkg/db"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

const (
	tableName = "access_rule"

	roleColumn     = "role"
	endpointColumn = "endpoint"
)

type repo struct {
	db db.Client
}

func NewRepository(db db.Client) repository.AccessRepository {
	return &repo{db: db}
}

func (r *repo) GetRuleByEndpoint(ctx context.Context, endpoint string) ([]*model.AccessRule, *db.Query, error) {
	builderSelect := sq.Select(roleColumn).
		From(tableName).
		PlaceholderFormat(sq.Dollar).
		Where(sq.Eq{endpointColumn: endpointColumn})

	query, args, err := builderSelect.ToSql()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to build query: %w", err)
	}

	q := db.Query{
		Name:     "get_rules_by_endpoint",
		QueryRow: query,
	}

	var rules []*repoModel.AccessRule
	err = r.db.DB().ScanAllContext(ctx, &rules, q, args...)

	if errors.Is(err, pgx.ErrNoRows) {
		return []*model.AccessRule{}, nil, nil
	}

	if err != nil {
		return nil, nil, fmt.Errorf("failed to get rules for endpoint %s: %w", endpoint, err)
	}

	result := make([]*model.AccessRule, len(rules))

	for _, rule := range rules {
		result = append(result, converter.ToAccessRuleFromRepo(rule))
	}

	return result, &q, nil
}
