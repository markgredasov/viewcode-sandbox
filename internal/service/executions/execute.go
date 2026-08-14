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

	return s.testsRunningPort.RunTests(ctx, exec, tests)
}
