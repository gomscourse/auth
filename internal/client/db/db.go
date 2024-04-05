package db

import (
	"context"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
)

// Client клиент для работы с базой данных
type Client interface {
	DB() DB
	Close() error
}

// Query обертка над запросом, хранящая имя и сам запрос
// Имя запроса используется для логирования и потенциально может использоваться где-то еще, например, для трейсинга
type Query struct {
	Name     string
	QueryRow string
}

// SQLExecutor комбинирует NamedExecutor и QueryExecutor
type SQLExecutor interface {
	NamedExecutor
	QueryExecutor
}

// NamedExecutor интерфейс для работы с именованными запросами с помощью тегов и структур
type NamedExecutor interface {
	ScanOneContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
	ScanAllContext(ctx context.Context, dest interface{}, q Query, args ...interface{}) error
}

// QueryExecutor интерфейс для работы с обычными запрсоами
// TODO: убрать привязку к постгрес
type QueryExecutor interface {
	ExecContext(ctx context.Context, q Query, args ...interface{}) (pgconn.CommandTag, error)
	QueryContext(ctx context.Context, q Query, args ...interface{}) (pgx.Rows, error)
	QueryRowContext(ctx context.Context, q Query, args ...interface{}) pgx.Row
}

// Pinger интерфейс для проверки соединения с БД
type Pinger interface {
	Ping(ctx context.Context) error
}

type DB interface {
	SQLExecutor
	Pinger
	BeginTx(ctx context.Context, txOpts pgx.TxOptions) (pgx.Tx, error)
	Close()
}

// Handler - функция, которая выполняется в транзакции
type Handler func(ctx context.Context) error

// TxManager - менеджер транзакций, который выполняет указанный пользователем обработчик в транзакции
type TxManager interface {
	ReadCommitted(ctx context.Context, handler Handler) error
}

// Transactor - интерфейс для работы с транзакциями
type Transactor interface {
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}
