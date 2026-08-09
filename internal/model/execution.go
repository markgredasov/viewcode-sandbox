package model

import (
	"time"

	"github.com/google/uuid"
)

type Execution struct {
	ProblemID uuid.UUID

	Language Language
	Code     Code
	Tests    Verification

	StdOut   *string
	StdErr   *string
	Duration *time.Duration
}

func NewExecution(id, language, code string) (Execution, error) {
	uuidObject, err := uuid.Parse(id)
	if err != nil {
		return Execution{}, err
	}

	languageObject, err := NewLanguage(language)
	if err != nil {
		return Execution{}, err
	}

	codeObject, err := NewCode(code)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		ProblemID: uuidObject,
		Language:  languageObject,
		Code:      codeObject,
	}, nil
}

type Result struct {
	StdOut            string
	StdErr            string
	ExecutionTimeInMS int
}
