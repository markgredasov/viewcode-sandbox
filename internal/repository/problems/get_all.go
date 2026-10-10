package problems

import (
	"context"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *Repository) GetAll(ctx context.Context, complexity string, categories []string) ([]model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	query := sq.
		Select("id, name, description, created_at, langs, complexity_rate, categories").
		From("problems").
		PlaceholderFormat(sq.Dollar)

	if complexity != "" {
		query = query.Where(sq.Eq{"complexity_rate": complexity})
	}

	if len(categories) != 0 {
		query = query.Where(sq.Expr("categories && ?::category[]", categories))
	}

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return []model.Problem{}, errs.Internal("building sql query", err)
	}

	rows, err := r.pool.Query(ctx, sqlQuery, args...)
	if err != nil {
		return []model.Problem{}, errs.Internal("rows query", err)
	}

	records, err := pgx.CollectRows(rows, pgx.RowToStructByName[ProblemRecord])
	if err != nil {
		return []model.Problem{}, errs.Internal("collect rows", err)
	}

	return toModels(records)
}
