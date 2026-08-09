package model

import (
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

type Problem struct {
	ProblemID   string
	Name        string
	Description string
	CreatedAt   time.Time
}

func NewProblem(name, description string) (Problem, error) {
	if utf8.RuneCountInString(name) <= 0 {
		return Problem{}, errs.ErrEmptyProblemName
	}

	if utf8.RuneCountInString(name) >= 70 {
		return Problem{}, errs.ErrTooLargeProblemName
	}

	if utf8.RuneCountInString(description) <= 0 {
		return Problem{}, errs.ErrEmptyProblemDescription
	}

	if utf8.RuneCountInString(description) >= 10000 {
		return Problem{}, errs.ErrTooLargeProblemDescription
	}

	return Problem{
		ProblemID:   uuid.NewString(),
		Name:        name,
		Description: description,
	}, nil
}
