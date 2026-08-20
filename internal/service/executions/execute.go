package service_executions

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

func (s *ExecutionsService) Execute(ctx context.Context, exec model.Execution) (model.Execution, error) {
	problemID := exec.ProblemID.String()

	tests, err := s.testsRepo.GetByProblemIDAndLanguage(ctx, problemID, exec.Language)
	if err != nil {
		return model.Execution{}, err
	}

	execution, err := s.testsRunningPort.RunTests(ctx, exec, tests)
	if err != nil {
		return model.Execution{}, err
	}

	testsResult, err := s.parserPort.Parse(execution.Language, execution.ExecutionResult.StdOut)
	if err != nil {
		return model.Execution{}, err
	}

	execution.WithTestsResult(testsResult)

	return execution, nil
}
