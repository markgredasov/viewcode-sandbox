package problems

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *Repository) GetByID(ctx context.Context, id string) (model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, description, created_at, langs
		FROM problems
		WHERE id = $1;
	`

	row := r.pool.QueryRow(ctx, sqlQuery, id)

	var pr ProblemRecord
	if err := row.Scan(&pr.ProblemID, &pr.Name, &pr.Description, &pr.CreatedAt, &pr.Languages); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Problem{}, fmt.Errorf("problem not found: %w", errs.ErrNotFound)
		}
		return model.Problem{}, err
	}

	problem, err := fromProblemRecordToModel(pr)
	if err != nil {
		return model.Problem{}, err
	}

	return problem, nil
}
