package problems

import database_postgres "github.com/markgredasov/viewcode-sandbox/internal/database/postgres"

type Repository struct {
	pool *database_postgres.Pool
}

func NewProblemsRepository(pool *database_postgres.Pool) *Repository {
	return &Repository{
		pool: pool,
	}
}
