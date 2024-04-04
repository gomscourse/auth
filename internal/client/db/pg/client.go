package pg

import (
	"context"
	"github.com/gomscourse/auth/internal/client/db"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/pkg/errors"
)

type pgClient struct {
	masterDBC db.DB
}

func (pc *pgClient) DB() db.DB {
	return pc.masterDBC
}

func (pc *pgClient) Close() error {
	if pc.masterDBC != nil {
		pc.masterDBC.Close()
	}

	return nil
}

func New(ctx context.Context, dsn string) (db.Client, error) {
	dbc, err := pgxpool.Connect(ctx, dsn)
	if err != nil {
		return nil, errors.Errorf("failed to connect to DB: %v", err)
	}

	return &pgClient{
		masterDBC: &pg{dbc: dbc},
	}, nil
}
