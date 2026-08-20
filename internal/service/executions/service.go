package service_executions

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ExecutionsService struct {
	testsRunningPort TestsRunningPort
	testsRepo        TestsRepository
	parserPort       ParserPort
}

type TestsRunningPort interface {
	RunTests(ctx context.Context, exec model.Execution, tests model.Verification) (model.Execution, error)
}

type TestsRepository interface {
	GetByProblemIDAndLanguage(ctx context.Context, problemID string, language model.Language) (model.Verification, error)
}

type ParserPort interface {
	Parse(language model.Language, stdOut string) (model.TestsResult, error)
}

func NewExecutionsService(testsRunningPort TestsRunningPort, testsRepo TestsRepository, parserPort ParserPort) *ExecutionsService {
	return &ExecutionsService{
		testsRunningPort: testsRunningPort,
		testsRepo:        testsRepo,
		parserPort:       parserPort,
	}
}
