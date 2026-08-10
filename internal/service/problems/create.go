package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ProblemsService) Create(ctx context.Context, name, description string) (model.Problem, error) {
	problem, err := model.NewProblemUnitialized(name, description)
	if err != nil {
		return model.Problem{}, err
	}

	return s.problemsRepository.Create(ctx, problem)
}
