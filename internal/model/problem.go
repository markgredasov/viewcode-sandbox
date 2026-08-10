package model

import (
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	errs "github.com/markgredasov/viewcode-sandbox/internal/errors"
)

const (
	maxProblemNameLength = 70
	maxDescriptionLength = 10_000
)

type Problem struct {
	ProblemID   string
	Name        string
	Description string
	CreatedAt   time.Time
	Languages   []Language
}

func NewProblem(problemID, name, description string, createdAt time.Time, languages []string) (Problem, error) {
	if utf8.RuneCountInString(name) <= 0 {
		return Problem{}, fmt.Errorf("empty problem name: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(name) >= maxProblemNameLength {
		return Problem{}, fmt.Errorf("problem name is too large: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) <= 0 {
		return Problem{}, fmt.Errorf("description is empty: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) >= maxDescriptionLength {
		return Problem{}, fmt.Errorf("description is too large: %w", errs.ErrInvalidArgument)
	}

	languageValueObjects, err := ConvertTo(languages)
	if err != nil {
		return Problem{}, fmt.Errorf("convert languages: %w", err)
	}

	return Problem{
		ProblemID:   problemID,
		Name:        name,
		Description: description,
		CreatedAt:   createdAt,
		Languages:   languageValueObjects,
	}, nil
}

func NewProblemUnitialized(name, description string) (Problem, error) {
	if utf8.RuneCountInString(name) <= 0 {
		return Problem{}, fmt.Errorf("empty problem name: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(name) >= maxProblemNameLength {
		return Problem{}, fmt.Errorf("problem name is too large: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) <= 0 {
		return Problem{}, fmt.Errorf("description is empty: %w", errs.ErrInvalidArgument)
	}

	if utf8.RuneCountInString(description) >= maxDescriptionLength {
		return Problem{}, fmt.Errorf("description is too large: %w", errs.ErrInvalidArgument)
	}

	return Problem{
		ProblemID:   uuid.NewString(),
		Name:        name,
		Description: description,
	}, nil
}
