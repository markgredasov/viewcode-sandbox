package handler_executions

import (
	"time"

	"github.com/markgredasov/viewcode-sandbox/internal/model/execution"
)

type ExecuteDTORequest struct {
	ProblemID string `json:"problem_id"`
	Language  string `json:"language"`
	Code      string `json:"code"`
}

type ExecuteDTOResponse struct {
	Submission SubmissionDTO `json:"submission"`
	Result     ResultDTO     `json:"result"`
}

type SubmissionDTO struct {
	ID                string
	ProblemID         string
	TestsTotal        int
	TestsPassed       int
	Status            string
	Language          string
	ExecutionTimeInMS int
	Timestamp         time.Time
}

type ResultDTO struct {
	StdOut   string        `json:"stdout"`
	StdErr   string        `json:"stderr"`
	Duration time.Duration `json:"duration"`
}

func fromModelToDTO(m execution.Execution) ExecuteDTOResponse {
	s := SubmissionDTO{
		ID:                m.Submission.ID.String(),
		ProblemID:         m.Submission.ProblemID.String(),
		TestsTotal:        m.Submission.TestsTotal,
		TestsPassed:       m.Submission.TestsPassed,
		Status:            m.Submission.Status.String(),
		Language:          m.Submission.Language.String(),
		ExecutionTimeInMS: m.Submission.ExecutionTimeInMS,
		Timestamp:         m.Submission.UpdatedAt,
	}

	r := ResultDTO{
		StdOut:   m.Result.StdOut,
		StdErr:   m.Result.StdErr,
		Duration: m.Result.Duration,
	}

	return ExecuteDTOResponse{
		Submission: s,
		Result:     r,
	}
}
