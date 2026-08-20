package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ProblemsService) Create(
	ctx context.Context,
	name string,
	description string,
	complexity string,
	categories []string,
) (model.Problem, error) {
	problem, err := model.NewProblemUnitialized(name, description)
	if err != nil {
		return model.Problem{}, err
	}

	complexityRate, err := model.NewComplexity(complexity)
	if err != nil {
		return model.Problem{}, err
	}

	categoryValueObjects, err := model.ConvertToValueObjsArray(categories, model.NewCategory)
	if err != nil {
		return model.Problem{}, err
	}

	problem.WithComplexityRate(complexityRate)
	problem.WithCategories(categoryValueObjects)

	return s.problemsRepository.Create(ctx, problem)
}
