package service_executions

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/execution"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
	"go.uber.org/zap"
)

func (s *ExecutionsService) Execute(ctx context.Context, userID, problemID, language, code string) (execution.Execution, error) {
	lang, err := model.NewLanguage(language)
	if err != nil {
		return execution.Execution{}, err
	}

	tests, err := s.testsRepo.GetByProblemIDAndLanguage(ctx, problemID, lang)
	if err != nil {
		return execution.Execution{}, err
	}

	log := logger.FromContext(ctx)

	subm, err := s.submissionsPort.CreateSubmission(ctx, userID, problemID, tests.TestsNumber, lang)
	if err != nil {
		log.Error("failed to create submission", zap.String("problemID", problemID), zap.String("userID", userID))
	}

	codeModel, err := model.NewCode(code)
	if err != nil {
		return execution.Execution{}, err
	}

	subm, err = s.submissionsPort.SetRunning(ctx, subm.ID.String())
	if err != nil {
		log.Error("failed to set submission running", zap.String("submissionID", subm.ID.String()))
	}

	execRes, err := s.testsRunningPort.RunTests(ctx, codeModel, tests)
	if err != nil {
		return execution.Execution{}, err
	}

	e, err := s.parseToExecution(ctx, log, execRes, subm, codeModel)
	if err != nil {
		return execution.Execution{}, err
	}

	return e, nil
}

func (s *ExecutionsService) parseToExecution(ctx context.Context, log *logger.Logger, res execution.ExecutionResult,
	subm submission.Submission, code model.Code) (execution.Execution, error) {
	testsResult, err := s.parserPort.Parse(model.LanguageGo, res.StdOut)
	if err != nil {
		return execution.Execution{}, err
	}

	log.Debug("tests result", zap.Any("tests result", testsResult))

	subm, err = s.submissionsPort.SetResult(ctx, subm.ID.String(), submission.StatusAccepted, testsResult.PassedTestsNumber, testsResult.PassingDurationInMs)
	if err != nil {
		log.Error("failed to set submission result", zap.String("submissionID", subm.ID.String()))
	}

	return execution.New(code, subm, res), nil
}
