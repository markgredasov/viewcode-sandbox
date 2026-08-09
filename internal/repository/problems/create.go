package problems

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *ProblemsRepository) Create(ctx context.Context, problem model.Problem) (model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		INSERT INTO public.problems (id, name, description)
		VALUES ($1, $2, $3)
		RETURNING id, name, description, created_at;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, problem.ProblemID, problem.Name, problem.Description)

	var createdProblem model.Problem
	if err := row.Scan(&createdProblem.ProblemID, &createdProblem.Name, &createdProblem.Description, &createdProblem.CreatedAt); err != nil {
		if _, ok := errors.AsType[*pgconn.PgError](err); ok {
			return model.Problem{}, errs.ErrProblemAlreadyExists
		}
		return model.Problem{}, err
	}

	return createdProblem, nil
}
