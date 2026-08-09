package problems

import database_postgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"

type ProblemsRepository struct {
	pool database_postgres.Pool
}

func NewProblemsRepository(pool database_postgres.Pool) *ProblemsRepository {
	return &ProblemsRepository{
		pool: pool,
	}
}
