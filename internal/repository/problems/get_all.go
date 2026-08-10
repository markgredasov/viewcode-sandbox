package problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (r *ProblemsRepository) GetAll(ctx context.Context) ([]model.Problem, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OpTimeout())
	defer cancel()

	sqlQuery := `
		SELECT id, name, description, created_at, langs
		FROM problems;
	`

	rows, err := r.pool.Query(ctx, sqlQuery)
	if err != nil {
		return []model.Problem{}, err
	}
	defer rows.Close()

	problemModels := make([]model.Problem, 0)
	for rows.Next() {
		var pr ProblemRecord

		err := rows.Scan(
			&pr.ProblemID,
			&pr.Name,
			&pr.Description,
			&pr.CreatedAt,
			&pr.Languages,
		)
		if err != nil {
			return []model.Problem{}, err
		}

		problemModel, err := fromProblemRecordToModel(pr)
		if err != nil {
			return []model.Problem{}, err
		}

		problemModels = append(problemModels, problemModel)
	}

	if err := rows.Err(); err != nil {
		return []model.Problem{}, err
	}

	return problemModels, nil
}
