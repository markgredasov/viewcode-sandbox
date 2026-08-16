package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgconn"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *Repository) Create(ctx context.Context, problem model.Problem) (model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		INSERT INTO public.problems (id, name, description, complexity_rate, categories)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, name, description, created_at, langs, complexity_rate, categories;
	`

	row := r.pool.QueryRow(
		ctx,
		sqlQuery,
		problem.ProblemID,
		problem.Name,
		problem.Description,
		problem.ComplexityRate.String(),
		model.ConvertFromValueObjsArray(problem.Categories),
	)

	var pr ProblemRecord
	if err := row.Scan(
		&pr.ProblemID,
		&pr.Name,
		&pr.Description,
		&pr.CreatedAt,
		&pr.Languages,
		&pr.Complexity,
		&pr.Categories,
	); err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return model.Problem{}, fmt.Errorf(
				"problem with this name='%s' already exists: %w",
				problem.Name,
				errs.ErrAlreadyExists,
			)
		}
		return model.Problem{}, err
	}

	createdProblem, err := fromProblemRecordToModel(pr)
	if err != nil {
		return model.Problem{}, err
	}

	return createdProblem, nil
}
