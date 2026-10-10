package service_executions

import (
	"context"
	"errors"

	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
	"github.com/markgredasov/viewcode-sandbox/internal/logger"
	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/execution"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
	"go.uber.org/zap"
)

func (s *ExecutionsService) Execute(ctx context.Context, userID, problemID, language, code string) (execution.Execution, error) {
	log := logger.FromContext(ctx)

	lang, err := model.NewLanguage(language)
	if err != nil {
		return execution.Execution{}, err
	}

	tests, err := s.testsRepo.GetByProblemIDAndLanguage(ctx, problemID, lang)
	if err != nil {
		return execution.Execution{}, err
	}

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

	execRes, err := s.runTests(ctx, subm, codeModel, tests)
	if err != nil {
		return execution.Execution{}, err
	}

	e, err := s.parseToExecution(ctx, log, execRes, subm, codeModel)
	if err != nil {
		return execution.Execution{}, err
	}

	return e, nil
}

func (s *ExecutionsService) runTests(ctx context.Context, subm submission.Submission,
	codeModel model.Code, tests model.Verification,
) (execution.ExecutionResult, error) {
	res, err := s.testsRunningPort.RunTests(ctx, codeModel, tests)
	if err != nil {
		status := submission.StatusInternal
		if errors.Is(err, errs.ErrTimeoutExceeded) {
			status = submission.StatusTimeout
		}

		_, pErr := s.submissionsPort.SetResult(ctx, subm.ID.String(), status, subm.TestsPassed, 0)
		if pErr != nil {
			return execution.ExecutionResult{}, errors.Join(err, pErr)
		}
		return execution.ExecutionResult{}, err
	}

	if res.StdErr != "" {
		err = errs.Internal("stderr is not empty", errs.ErrInternal)
		_, pErr := s.submissionsPort.SetResult(ctx, subm.ID.String(), submission.StatusInternal, subm.TestsPassed, 0)
		if pErr != nil {
			return execution.ExecutionResult{}, errors.Join(err, pErr)
		}
		return execution.ExecutionResult{}, err
	}

	return res, nil
}

func (s *ExecutionsService) parseToExecution(ctx context.Context, log *logger.Logger, res execution.ExecutionResult,
	subm submission.Submission, code model.Code,
) (execution.Execution, error) {
	testsResult, err := s.parserPort.Parse(model.LanguageGo, res.StdOut)
	if err != nil {
		return execution.Execution{}, err
	}

	var status submission.Status
	if testsResult.TotalTestsNumber == testsResult.PassedTestsNumber {
		status = submission.StatusAccepted
	} else {
		status = submission.StatusWrongAnswer
	}

	subm, err = s.submissionsPort.SetResult(ctx, subm.ID.String(), status, testsResult.PassedTestsNumber, testsResult.PassingDurationInMs)
	if err != nil {
		log.Error("failed to set submission result", zap.String("submissionID", subm.ID.String()))
		subm.TestsTotal = testsResult.TotalTestsNumber
		subm.TestsPassed = testsResult.PassedTestsNumber
		subm.ExecutionTimeInMS = int(testsResult.PassingDurationInMs)
		subm.Status = submission.StatusAccepted
	}

	return execution.New(code, subm, res), nil
}
