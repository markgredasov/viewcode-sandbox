package service_executions

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ExecutionsService struct {
	testsRunningPort TestsRunningPort
	testsRepo        TestsRepository
}

type TestsRunningPort interface {
	RunTests(ctx context.Context, exec model.Execution, tests model.Verification) (model.Execution, error)
}

type TestsRepository interface {
	GetByID(ctx context.Context, problemID string) (model.Verification, error)
}

func NewExecutionsService(testsRunningPort TestsRunningPort, testsRepo TestsRepository) *ExecutionsService {
	return &ExecutionsService{
		testsRunningPort: testsRunningPort,
		testsRepo:        testsRepo,
	}
}
