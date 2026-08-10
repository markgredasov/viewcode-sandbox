package errs

import "errors"

var (
	ErrNoTestsProvided      = errors.New("no tests found")
	ErrProblemAlreadyExists = errors.New("problem with such name already exists")
	ErrInvalidArgument      = errors.New("invalid argument")
	ErrNotFound             = errors.New("not found")
)
