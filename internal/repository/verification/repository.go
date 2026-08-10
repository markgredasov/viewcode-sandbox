package verification

import database_postgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"

type CodeVerificationRepository struct {
	pool database_postgres.Pool
}

func NewCodeVerificationRepository(pool database_postgres.Pool) *CodeVerificationRepository {
	return &CodeVerificationRepository{
		pool: pool,
	}
}
