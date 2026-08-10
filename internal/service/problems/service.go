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
	GetAll(ctx context.Context) ([]model.Problem, error)
	GetByID(ctx context.Context, id string) (model.Problem, error)
	Delete(ctx context.Context, id string) error
}

func NewProblemsService(problemsRepo ProblemsRepository) *ProblemsService {
	return &ProblemsService{
		problemsRepository: problemsRepo,
	}
}
