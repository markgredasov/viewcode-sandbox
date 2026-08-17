package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ProblemsService) GetAll(ctx context.Context, complexity *string, categories *[]string) ([]model.Problem, error) {
	var complexityValueObj model.Complexity
	if complexity != nil {
		val, err := model.NewComplexity(*complexity)
		if err != nil {
			return []model.Problem{}, err
		}
		complexityValueObj = val
	}

	var categoriesValueObj []model.Category
	if categories != nil {
		val, err := model.ConvertToValueObjsArray(*categories, model.NewCategory)
		if err != nil {
			return []model.Problem{}, err
		}
		categoriesValueObj = val
	}

	return s.problemsRepository.GetAll(ctx, complexityValueObj.String(), model.ConvertFromValueObjsArray(categoriesValueObj))
}
