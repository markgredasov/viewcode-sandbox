package dbpostgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type QueryExecutor interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
}

type txKeyType struct{}

var txKey txKeyType

type TxManager struct {
	pool *Pool
}

func NewTxManager(pool *Pool) *TxManager {
	return &TxManager{pool: pool}
}

func (tm *TxManager) WithinTx(ctx context.Context, fn func(ctx context.Context) error) error {
	tx, err := tm.pool.BeginTx(ctx, pgx.TxOptions{})
	if err != nil {
		return errs.Internal("begin tx", err)
	}

	txCtx := context.WithValue(ctx, txKey, tx)
	if err := fn(txCtx); err != nil {
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
			return errs.Internal("rollback tx", rollbackErr)
		}

		return err
	}

	if err := tx.Commit(ctx); err != nil {
		return errs.Internal("commit tx", err)
	}

	return nil
}

func ExecutorFromContext(ctx context.Context, fallback QueryExecutor) QueryExecutor {
	if tx, ok := ctx.Value(txKey).(pgx.Tx); ok {
		return tx
	}
	return fallback
}
