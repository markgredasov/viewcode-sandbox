package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ProblemsService) GetByID(ctx context.Context, id string) (model.Problem, error) {
	return s.problemsRepository.GetByID(ctx, id)
}
