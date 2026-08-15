package model

import (
	"time"

	utils_string "github.com/markgredasov/viewcode-sandbox/internal/utils/string"
)

type TestsResult struct {
	TotalTestsNumber    int64
	PassedTestsNumber   int64
	FailedTestLog       string
	PassingDurationInMs float64
}

func NewTestsResult(totalTestsNumber, passedTestsNumber int64, failedTestLog string, duration float64) TestsResult {
	return TestsResult{
		TotalTestsNumber:    totalTestsNumber,
		PassedTestsNumber:   passedTestsNumber,
		FailedTestLog:       failedTestLog,
		PassingDurationInMs: duration,
	}
}

type ExecutionResult struct {
	StdOut   string
	StdErr   string
	Duration time.Duration
}

func NewExecutionResult(stdOut string, stdErr string, duration time.Duration) ExecutionResult {
	cleanedStdOut := utils_string.CleanString(stdOut)
	cleanedStdErr := utils_string.CleanString(stdErr)
	return ExecutionResult{
		StdOut:   cleanedStdOut,
		StdErr:   cleanedStdErr,
		Duration: duration,
	}
}
