package problems

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ProblemRecord struct {
	ProblemID   string
	Name        string
	Description string
	CreatedAt   time.Time
	Languages   []string
	Complexity  string
	Categories  []string
}

func fromProblemRecordToModel(pr ProblemRecord) (model.Problem, error) {
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
