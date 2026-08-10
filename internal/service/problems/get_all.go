package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ProblemsService) GetAll(ctx context.Context) ([]model.Problem, error) {
	return s.problemsRepository.GetAll(ctx)
}
