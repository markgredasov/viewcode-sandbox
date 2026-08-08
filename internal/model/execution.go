package model

import "time"

type Execution struct {
	Language Language
	Code     Code
	Tests    Tests

	StdOut   *string
	StdErr   *string
	Duration *time.Duration
}

func NewExecution(language, code string) (Execution, error) {
	languageObject, err := NewLanguage(language)
	if err != nil {
		return Execution{}, err
	}

	codeObject, err := NewCode(code)
	if err != nil {
		return Execution{}, err
	}

	return Execution{
		Language: languageObject,
		Code:     codeObject,
	}, nil
}

type Result struct {
	StdOut            string
	StdErr            string
	ExecutionTimeInMS int
}
