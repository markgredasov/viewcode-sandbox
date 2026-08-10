package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *ProblemsRepository) GetByID(ctx context.Context, id string) (model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, description, created_at
		FROM problems
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, id)

	var problem model.Problem
	if err := row.Scan(&problem.ProblemID, &problem.Name, &problem.Description, &problem.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Problem{}, fmt.Errorf("problem not found: %w", errs.ErrNotFound)
		}
		return model.Problem{}, err
	}

	return problem, nil
}
