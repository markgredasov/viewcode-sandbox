package dbpostgres

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Pool struct {
	pool *pgxpool.Pool

	opTimeout time.Duration
}

type rowsWithCancel struct {
	pgx.Rows
	cancel context.CancelFunc
}

type rowWithCancel struct {
	pgx.Row
	cancel context.CancelFunc
}

func NewConnectionPool(
	ctx context.Context,
	config Config,
) (*Pool, error) {
	hostAndPort := net.JoinHostPort(config.Host, config.Port)
	connectionString := fmt.Sprintf(
		"postgres://%s:%s@%s/%s?sslmode=disable",
		config.Username,
		config.Password,
		hostAndPort,
		config.DB,
	)

	pgxConfig, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		return nil, fmt.Errorf("parse pgxconfig: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("create pgxpool: %w", err)
	}

	if err = pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pgxpool ping: %w", err)
	}

	return &Pool{
		pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *Pool) Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error) {
	ctx, cancel := context.WithTimeout(ctx, p.OpTimeout())
	rows, err := p.pool.Query(ctx, sql, args...)
	if err != nil {
		cancel()
		return nil, err
	}

	return &rowsWithCancel{
		Rows:   rows,
		cancel: cancel,
	}, nil
}

func (p *Pool) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	ctx, cancel := context.WithTimeout(ctx, p.OpTimeout())
	return &rowWithCancel{
		Row:    p.pool.QueryRow(ctx, sql, args...),
		cancel: cancel,
	}
}

func (p *Pool) Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error) {
	ctx, cancel := context.WithTimeout(ctx, p.OpTimeout())
	defer cancel()
	return p.pool.Exec(ctx, sql, args...)
}

func (p *Pool) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	ctx, cancel := context.WithTimeout(ctx, p.OpTimeout())
	defer cancel()
	return p.pool.BeginTx(ctx, txOptions)
}

func (p *Pool) OpTimeout() time.Duration {
	return p.opTimeout
}

func (r *rowWithCancel) Scan(dest ...any) error {
	defer r.cancel()
	return r.Row.Scan(dest...)
}

func (p *Pool) Close() {
	p.pool.Close()
}

func (r *rowsWithCancel) Close() {
	r.Rows.Close()
	r.cancel()
}
