package verification

import database_postgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"

type VerificationRepository struct {
	pool database_postgres.Pool
}

func NewVerificationRepository(pool database_postgres.Pool) *VerificationRepository {
	return &VerificationRepository{
		pool: pool,
	}
}
