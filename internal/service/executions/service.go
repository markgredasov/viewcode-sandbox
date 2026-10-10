package service_executions

import (
	"context"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/execution"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
)

type ExecutionsService struct {
	testsRunningPort TestsRunningPort
	testsRepo        TestsRepository
	parserPort       ParserPort
	submissionsPort  SubmissionsPort
}

type TestsRunningPort interface {
	RunTests(ctx context.Context, code model.Code, tests model.Verification) (execution.ExecutionResult, error)
}

type TestsRepository interface {
	GetByProblemIDAndLanguage(ctx context.Context, problemID string, language model.Language) (model.Verification, error)
}

type ParserPort interface {
	Parse(language model.Language, stdOut string) (model.TestsResult, error)
}

type SubmissionsPort interface {
	CreateSubmission(ctx context.Context, userID, problemID string, testsTotal int, language model.Language) (submission.Submission, error)
	SetRunning(ctx context.Context, id string) (submission.Submission, error)
	SetResult(ctx context.Context, id string, status submission.Status, testsPassed int, executionTimeInMS float64) (submission.Submission, error)
}

func NewExecutionsService(testsRunningPort TestsRunningPort, testsRepo TestsRepository,
	parserPort ParserPort, submissionsPort SubmissionsPort,
) *ExecutionsService {
	return &ExecutionsService{
		testsRunningPort: testsRunningPort,
		testsRepo:        testsRepo,
		parserPort:       parserPort,
		submissionsPort:  submissionsPort,
	}
}
