package problems

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	dbpostgres "github.com/markgredasov/viewcode-sandbox/internal/platform/database/postgres"
)

func (r *Repository) GetByID(ctx context.Context, id string) (model.Problem, error) {
	sqlQuery := `
		SELECT id, name, description, created_at, langs, complexity_rate, categories
		FROM problems
		WHERE id = $1;
	`

	conn := dbpostgres.ExecutorFromContext(ctx, r.pool)
	rows, err := conn.Query(ctx, sqlQuery, id)
	if err != nil {
		return model.Problem{}, errs.Internal("query row", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ProblemRecord])
	if err != nil {
		msg := "scan row"
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Problem{}, errs.NotFound("problem not found")
		}
		return model.Problem{}, errs.Internal(msg, err)
	}

	return toModel(rec)
}
