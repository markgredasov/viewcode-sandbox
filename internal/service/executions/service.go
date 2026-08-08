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
	RunTests(ctx context.Context, exec model.Execution, tests model.Tests) (model.Execution, error)
}

type TestsRepository interface {
	Get(ctx context.Context, problemID string) (model.Tests, error)
}

func NewExecutionsService(testsRunningPort TestsRunningPort) *ExecutionsService {
	return &ExecutionsService{
		testsRunningPort: testsRunningPort,
	}
}
