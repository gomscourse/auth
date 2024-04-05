package transaction

import (
	"context"
	"github.com/gomscourse/auth/internal/client/db"
	"github.com/gomscourse/auth/internal/client/db/pg"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager - создает новый менеджер транзакций, который удовлетворяет интерфейсу db.TxManager
func NewTransactionManager(transactor db.Transactor) db.TxManager {
	return &manager{
		db: transactor,
	}
}

func (m *manager) ReadCommitted(ctx context.Context, handler db.Handler) error {
	txOpts := pgx.TxOptions{IsoLevel: pgx.ReadCommitted}
	return m.transaction(ctx, txOpts, handler)
}

// transaction - основная функция, которая выполняет указанный пользователем обработчик в транзакции
func (m *manager) transaction(ctx context.Context, opts pgx.TxOptions, handler db.Handler) (err error) {
	// если это вложенная транзакция, пропускаем инициацию новой транзакции и запускаем обработчик
	tx, ok := ctx.Value(pg.TxKey).(pgx.Tx)
	if ok {
		// в данном случае ошибка возвращенная обработчиком будет записана в именованное возвращаемео занчение (err)
		// и отложенная функция будет иметь к ней доступ и сможет также закомитить или откатить транзакцию
		return handler(ctx)
	}

	// стартуем новую транзакцию
	tx, err = m.db.BeginTx(ctx, opts)
	if err != nil {
		return errors.Wrap(err, "can't begin transaction")
	}

	// кладем транзакцию в контекст
	ctx = pg.MakeContextTx(ctx, tx)

	// настраиваем функцию отсрочки для отката или комита транзакции
	defer func() {
		// восстанавливаемся после паники
		if r := recover(); r != nil {
			err = errors.Errorf("panic recovered: %v", err)
		}

		// откатываем транзацию если произошла ошибка
		if err != nil {
			if errRollback := tx.Rollback(ctx); errRollback != nil {
				err = errors.Wrapf(err, "rollback error: %v", errRollback)
			}

			return
		}

		if err == nil {
			err = tx.Commit(ctx)
			if err != nil {
				err = errors.Wrap(err, "commit failed")
			}
		}
	}()

	// выполняем код внутри транзакции
	// если фукнци потерпела неудачу, записываем ошибку и функция отсрочки выполнит откат
	// или в противном случае закомитит транзакцию
	if err = handler(ctx); err != nil {
		err = errors.Wrap(err, "failed executing handler inside transaction")
	}

	return err
}
