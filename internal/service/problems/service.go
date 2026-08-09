package service_problems

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ProblemsService struct {
	problemsRepository ProblemsRepository
}

type ProblemsRepository interface {
	Create(ctx context.Context, problem model.Problem) (model.Problem, error)
}

func NewProblemsService(problemsRepo ProblemsRepository) *ProblemsService {
	return &ProblemsService{
		problemsRepository: problemsRepo,
	}
}
