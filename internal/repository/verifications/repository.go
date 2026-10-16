package verifications

import (
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"

	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"
)

type CodeVerificationRepository struct {
	pool   *dbpostgres.ConnectionPool
	getter *trmpgx.CtxGetter
}

func NewCodeVerificationRepository(pool *dbpostgres.ConnectionPool, getter *trmpgx.CtxGetter) *CodeVerificationRepository {
	return &CodeVerificationRepository{
		pool:   pool,
		getter: getter,
	}
}
