package execution

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model"
	"github.com/markgredasov/viewcode-sandbox/internal/model/submission"
	utils_string "github.com/markgredasov/viewcode-sandbox/internal/utils/string"
)

type ExecutionResult struct {
	StdOut   string
	StdErr   string
	Duration time.Duration
}

func NewResult(stdOut string, stdErr string, duration time.Duration) ExecutionResult {
	cleanedStdOut := utils_string.CleanString(stdOut)
	cleanedStdErr := utils_string.CleanString(stdErr)
	return ExecutionResult{
		StdOut:   cleanedStdOut,
		StdErr:   cleanedStdErr,
		Duration: duration,
	}
}

type Execution struct {
	Code       model.Code
	Submission submission.Submission
	Result     ExecutionResult
}

func New(code model.Code, submission submission.Submission, result ExecutionResult) Execution {
	return Execution{
		Code:       code,
		Submission: submission,
		Result:     result,
	}
}
