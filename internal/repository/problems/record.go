package problems

import (
	"time"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ProblemRecord struct {
	ProblemID   string    `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"description"`
	CreatedAt   time.Time `db:"created_at"`
	Languages   []string  `db:"langs"`
	Complexity  string    `db:"complexity_rate"`
	Categories  []string  `db:"categories"`
}

func toModel(pr ProblemRecord) (model.Problem, error) {
	problem, err := model.NewProblem(
		pr.ProblemID,
		pr.Name,
		pr.Description,
		pr.CreatedAt,
		pr.Languages,
		pr.Complexity,
		pr.Categories,
	)
	if err != nil {
		return model.Problem{}, err
	}

	return problem, nil
}

func toModels(r []ProblemRecord) ([]model.Problem, error) {
	problems := make([]model.Problem, len(r))
	for i := range r {
		problem, err := toModel(r[i])
		if err != nil {
			return nil, errs.Internal("scan model", err)
		}
		problems[i] = problem
	}
	return problems, nil
}
