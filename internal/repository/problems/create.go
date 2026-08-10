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
		RETURNING id, name, description, created_at, langs;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, problem.ProblemID, problem.Name, problem.Description)

	var pr ProblemRecord
	if err := row.Scan(
		&pr.ProblemID,
		&pr.Name,
		&pr.Description,
		&pr.CreatedAt,
		&pr.Languages,
	); err != nil {
		if _, ok := errors.AsType[*pgconn.PgError](err); ok {
			return model.Problem{}, errs.ErrProblemAlreadyExists
		}
		return model.Problem{}, err
	}

	createdProblem, err := fromProblemRecordToModel(pr)
	if err != nil {
		return model.Problem{}, err
	}

	return createdProblem, nil
}
