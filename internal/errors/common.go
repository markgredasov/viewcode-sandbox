package errs

import "errors"

var (
	ErrInvalidLanguage            = errors.New("this language is unavailable")
	ErrEmptyCode                  = errors.New("code is empty")
	ErrTooLargeCode               = errors.New("code is too large")
	ErrNoTestsProvided            = errors.New("no tests found")
	ErrEmptyProblemName           = errors.New("empty problem name")
	ErrTooLargeProblemName        = errors.New("too large problem name")
	ErrEmptyProblemDescription    = errors.New("empty problem description")
	ErrTooLargeProblemDescription = errors.New("too large problem description")
	ErrProblemAlreadyExists       = errors.New("problem with such name already exists")
)
