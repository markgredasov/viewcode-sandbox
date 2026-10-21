package verifications

import (
	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/platform/database/postgres"
)

type CodeVerificationRepository struct {
	pool *dbpostgres.Pool
}

func NewCodeVerificationRepository(pool *dbpostgres.Pool) *CodeVerificationRepository {
	return &CodeVerificationRepository{
		pool: pool,
	}
}
