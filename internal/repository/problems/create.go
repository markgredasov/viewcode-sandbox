package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
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

	rows, err := r.pool.Query(ctx, sqlQuery, problem.ProblemID, problem.Name, problem.Description,
		problem.ComplexityRate.String(), model.ConvertFromValueObjsArray(problem.Categories))
	if err != nil {
		return model.Problem{}, errs.Internal("query rows", err)
	}

	rec, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[ProblemRecord])
	if err != nil {
		msg := "row scan"
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			msg = fmt.Sprintf("problem with this name='%s' already exists", problem.Name)
			return model.Problem{}, errs.AlreadyExists(msg)
		}
		return model.Problem{}, errs.Internal(msg, err)
	}

	return toModel(rec)
}
