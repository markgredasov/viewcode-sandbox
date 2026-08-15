package handler_executions

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
)

type ExecuteDTORequest struct {
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

type ExecuteDTOResponse struct {
	IsPassed        bool                       `json:"is_passed"`
	TestsResult     TestsResultDTOResponse     `json:"tests"`
	ExecutionResult ExecutionResultDTOResponse `json:"execution"`
}

type TestsResultDTOResponse struct {
	TotalTestsNumber  int64   `json:"total"`
	PassedTestsNumber int64   `json:"passed"`
	FailedTestLog     string  `json:"failed_log"`
	PassingDuration   float64 `json:"passing_duration_in_ms"`
}

type ExecutionResultDTOResponse struct {
	StdOut   string        `json:"stdout"`
	StdErr   string        `json:"stderr"`
	Duration time.Duration `json:"duration"`
}

func fromModelToDTO(m model.Execution) ExecuteDTOResponse {
	executionInfo := ExecutionResultDTOResponse{
		StdOut:   m.ExecutionResult.StdOut,
		StdErr:   m.ExecutionResult.StdErr,
		Duration: m.ExecutionResult.Duration,
	}

	testsInfo := TestsResultDTOResponse{
		TotalTestsNumber:  m.TestsResult.TotalTestsNumber,
		PassedTestsNumber: m.TestsResult.PassedTestsNumber,
		FailedTestLog:     m.TestsResult.FailedTestLog,
		PassingDuration:   m.TestsResult.PassingDurationInMs,
	}

	isPassed := testsInfo.PassedTestsNumber == testsInfo.TotalTestsNumber

	return ExecuteDTOResponse{
		IsPassed:        isPassed,
		TestsResult:     testsInfo,
		ExecutionResult: executionInfo,
	}
}
